
package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	Version       = 1
	MaxHeaderSize = 16 * 1024
)

const (
	TypeMessage  = "message"
	TypeFile     = "file"
	TypeACK      = "ack"
	TypeFSWStart = "fsw_start"
	TypeFSWStop  = "fsw_stop"
)

//json metadata header.
type Header struct {
	Version   int    `json:"version"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
	Algorithm string `json:"algorithm"`
	Size      int64  `json:"size"`
	MessageID string `json:"message_id"`
	Filename  string `json:"filename,omitempty"`
	IV        string `json:"iv,omitempty"`
	Hash      string `json:"hash,omitempty"`
}

//validate proverava da li header prati protokol
func (h Header) Validate() error {
	if h.Version != Version {
		return fmt.Errorf("unsupported protocol version: %d", h.Version)
	}

	switch h.Type {
	case TypeMessage, TypeFile, TypeACK, TypeFSWStart, TypeFSWStop:
	default:
		return fmt.Errorf("unknown packet type: %s", h.Type)
	}

	if strings.TrimSpace(h.Sender) == "" {
		return fmt.Errorf("missing sender")
	}

	if _, err := time.Parse(time.RFC3339Nano, h.Timestamp); err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}

	if strings.TrimSpace(h.MessageID) == "" {
		return fmt.Errorf("missing message_id")
	}

	if h.Size < 0 {
		return fmt.Errorf("negative payload size")
	}

	if h.Type == TypeFile && h.Filename == "" {
		return fmt.Errorf("missing filename")
	}

	if h.Type == TypeMessage || h.Type == TypeFile {
		if h.Algorithm != "TEA" && h.Algorithm != "LEA-CTR" {
			return fmt.Errorf("unsupported algorithm: %s", h.Algorithm)
		}

		if len(h.Hash) != 64 {
			return fmt.Errorf("hash must contain 64 hex characters")
		}

		if _, err := hex.DecodeString(h.Hash); err != nil {
			return fmt.Errorf("invalid hash: %w", err)
		}

		if h.Algorithm == "LEA-CTR" && h.IV == "" {
			return fmt.Errorf("missing IV for LEA-CTR")
		}
	} else {
		//kontrolni paketi zasad nemaju payload.
		if h.Size != 0 || h.Algorithm != "none" {
			return fmt.Errorf("invalid control packet")
		}
	}

	return nil
}

// WriteFrame zapisuje:
// [4B duzina headera][JSON header][encrypted payload]
func WriteFrame(w io.Writer, h Header, payload io.Reader) error {
	if err := h.Validate(); err != nil {
		return err
	}

	if h.Size > 0 && payload == nil {
		return fmt.Errorf("missing payload")
	}

	headerBytes, err := json.Marshal(h)
	if err != nil {
		return err
	}

	if len(headerBytes) > MaxHeaderSize {
		return fmt.Errorf("header too large")
	}

	var prefix [4]byte
	binary.BigEndian.PutUint32(
		prefix[:],
		uint32(len(headerBytes)),
	)

	if _, err := io.Copy(w, bytes.NewReader(prefix[:])); err != nil {
		return err
	}

	if _, err := io.Copy(w, bytes.NewReader(headerBytes)); err != nil {
		return err
	}

	if h.Size > 0 {
		if _, err := io.CopyN(w, payload, h.Size); err != nil {
			return err
		}
	}

	return nil
}

// readHeader cita duzinu headera i JSON metadata header.
func ReadHeader(r io.Reader) (Header, error) {
	var h Header
	var prefix [4]byte

	if _, err := io.ReadFull(r, prefix[:]); err != nil {
		return h, err
	}

	headerSize := binary.BigEndian.Uint32(prefix[:])

	if headerSize == 0 || headerSize > MaxHeaderSize {
		return h, fmt.Errorf("invalid header size: %d", headerSize)
	}

	headerBytes := make([]byte, headerSize)

	if _, err := io.ReadFull(r, headerBytes); err != nil {
		return h, err
	}

	if err := json.Unmarshal(headerBytes, &h); err != nil {
		return h, err
	}

	if err := h.Validate(); err != nil {
		return h, err
	}

	return h, nil
}

// ReadPayload cita tacno h.Size bajtova, podaci se upisuju u dst bez ucitavanja celog fajla u ram
func ReadPayload(r io.Reader, h Header, dst io.Writer) error {
	if h.Size < 0 {
		return fmt.Errorf("negative payload size")
	}

	_, err := io.CopyN(dst, r, h.Size)
	return err
}
