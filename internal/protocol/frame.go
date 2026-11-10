package protocol

import (
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	Version          uint8 = 1
	FixedHeaderSize        = 101
	MaxFilenameBytes       = 255
)

// Tipovi paketa
type PacketType uint8

const (
	TypeMessage  PacketType = 1
	TypeFile     PacketType = 2
	TypeACK      PacketType = 3
	TypeFSWStart PacketType = 4
	TypeFSWStop  PacketType = 5
)

// Algoritmi
type Algorithm uint8

const (
	AlgorithmNone   Algorithm = 0
	AlgorithmTEA    Algorithm = 1
	AlgorithmLEACTR Algorithm = 2
)

// Binarni header
type Header struct {
	Size      uint32
	Version   uint8
	Type      PacketType
	Sender    string
	Timestamp int64
	Algorithm Algorithm
	IV        [16]byte
	SHA256    [32]byte
	Filename  string
}

// Provera ispravnosti headera
func (h Header) Validate() error {
	if h.Version != Version {
		return fmt.Errorf("unsupported version: %d", h.Version)
	}

	switch h.Type {
	case TypeMessage, TypeFile, TypeACK, TypeFSWStart, TypeFSWStop:
	default:
		return fmt.Errorf("unknown packet type: %d", h.Type)
	}

	if !validUUID(h.Sender) {
		return errors.New("sender must be a valid 36-byte UUID")
	}

	if len(h.Filename) > MaxFilenameBytes ||
		!utf8.ValidString(h.Filename) {
		return errors.New("invalid filename")
	}

	if h.Filename == "." || h.Filename == ".." ||
		strings.ContainsAny(h.Filename, "/\\\x00") {
		return errors.New("unsafe filename")
	}

	switch h.Type {
	case TypeMessage, TypeFile:
		if h.Algorithm != AlgorithmTEA &&
			h.Algorithm != AlgorithmLEACTR {
			return errors.New("unsupported encryption algorithm")
		}

		if h.Type == TypeFile && h.Filename == "" {
			return errors.New("missing filename")
		}

		if h.Type == TypeMessage {
			if h.Size == 0 || h.Filename != "" {
				return errors.New("invalid message packet")
			}
		}

	default:
		// ACK i FSW komande nemaju payload.
		if h.Size != 0 ||
			h.Algorithm != AlgorithmNone ||
			h.Filename != "" {
			return errors.New("invalid control packet")
		}

		if h.IV != [16]byte{} ||
			h.SHA256 != [32]byte{} {
			return errors.New("invalid control packet metadata")
		}
	}

	return nil
}

// Slanje kompletnog frejma:
// [fixed header][filename][encrypted payload]
func WriteFrame(
	w io.Writer,
	h Header,
	ciphertext io.Reader,
) error {
	if err := h.Validate(); err != nil {
		return err
	}

	if h.Size > 0 && ciphertext == nil {
		return errors.New("missing ciphertext")
	}

	var fixed [FixedHeaderSize]byte

	// SIZE: bajtovi 0-3
	binary.BigEndian.PutUint32(fixed[0:4], h.Size)

	// VERSION: bajt 4
	fixed[4] = h.Version

	// TYPE: bajt 5
	fixed[5] = byte(h.Type)

	// SENDER UUID: bajtovi 6-41
	copy(fixed[6:42], h.Sender)

	// TIMESTAMP: bajtovi 42-49
	binary.BigEndian.PutUint64(
		fixed[42:50],
		uint64(h.Timestamp),
	)

	// ALGORITHM: bajt 50
	fixed[50] = byte(h.Algorithm)

	// IV: bajtovi 51-66
	copy(fixed[51:67], h.IV[:])

	// SHA-256: bajtovi 67-98
	copy(fixed[67:99], h.SHA256[:])

	// FILENAME SIZE: bajtovi 99-100
	binary.BigEndian.PutUint16(
		fixed[99:101],
		uint16(len(h.Filename)),
	)

	// Upis fiksnog dela headera
	if err := writeAll(w, fixed[:]); err != nil {
		return err
	}

	// Upis naziva fajla, ako postoji
	if err := writeAll(w, []byte(h.Filename)); err != nil {
		return err
	}

	// Kontrolni paket: nema payload-a
	if h.Size == 0 {
		return nil
	}

	// Upis tacno SIZE bajtova
	_, err := io.CopyN(
		w,
		ciphertext,
		int64(h.Size),
	)

	return err
}

// Citanje binarnog headera iz TCP stream-a
func ReadHeader(r io.Reader) (Header, error) {
	var h Header
	var fixed [FixedHeaderSize]byte

	// Citamo tacno 101 bajt.
	if _, err := io.ReadFull(r, fixed[:]); err != nil {
		return h, err
	}

	h.Size = binary.BigEndian.Uint32(fixed[0:4])
	h.Version = fixed[4]
	h.Type = PacketType(fixed[5])
	h.Sender = string(fixed[6:42])

	h.Timestamp = int64(
		binary.BigEndian.Uint64(fixed[42:50]),
	)

	h.Algorithm = Algorithm(fixed[50])

	copy(h.IV[:], fixed[51:67])
	copy(h.SHA256[:], fixed[67:99])

	filenameSize := binary.BigEndian.Uint16(
		fixed[99:101],
	)

	if filenameSize > MaxFilenameBytes {
		return h, errors.New("filename too long")
	}

	// Naziv fajla je promenljive duzine.
	filename := make([]byte, filenameSize)

	if _, err := io.ReadFull(r, filename); err != nil {
		return h, err
	}

	h.Filename = string(filename)

	if err := h.Validate(); err != nil {
		return h, err
	}

	return h, nil
}

// Cita ciphertext i proverava SHA-256.
// Prosledjeni hasher mora biti nova SHA-256 instanca.
// Koristicemo nasu implementaciju SHA-256 kada je napravimo.
func ReadPayloadVerified(
	r io.Reader,
	h Header,
	dst io.Writer,
	hasher hash.Hash,
) error {
	if err := h.Validate(); err != nil {
		return err
	}

	if h.Size == 0 && h.Type != TypeFile {
		return nil
	}

	if dst == nil || hasher == nil {
		return errors.New("missing destination or hasher")
	}

	if hasher.Size() != 32 {
		return errors.New("SHA-256 hasher required")
	}

	// Istovremeno upisujemo ciphertext i racunamo hash.
	_, err := io.CopyN(
		io.MultiWriter(dst, hasher),
		r,
		int64(h.Size),
	)

	if err != nil {
		return err
	}

	// Provera integriteta pre desifrovanja
	if subtle.ConstantTimeCompare(
		hasher.Sum(nil),
		h.SHA256[:],
	) != 1 {
		return errors.New("ciphertext SHA-256 mismatch")
	}

	return nil
}

// Upisuje sve bajtove i kada Writer vraca krace upise.
func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)

		if err != nil {
			return err
		}

		if n <= 0 {
			return io.ErrShortWrite
		}

		data = data[n:]
	}

	return nil
}

// Proverava format UUID-a sa crticama.
func validUUID(s string) bool {
	if len(s) != 36 {
		return false
	}

	for i := 0; i < len(s); i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if s[i] != '-' {
				return false
			}
			continue
		}

		c := s[i]

		if !((c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}
