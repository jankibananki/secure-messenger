package protocol

import (
	"bytes"
	"crypto/sha256"
	"io"
	"testing"
	"time"
)

// crypto/sha256 koristimo samo u testovima.
// U aplikaciji cemo koristiti sopstvenu implementaciju.

const testSender = "550e8400-e29b-41d4-a716-446655440000"

func TestMessageRoundTrip(t *testing.T) {
	// Testni ciphertext, bez stvarnog sifrovanja.
	payload := []byte("encrypted test message")

	header := Header{
		Version:   Version,
		Type:      TypeMessage,
		Sender:    testSender,
		Timestamp: time.Now().UnixMilli(),
		Algorithm: AlgorithmTEA,
		Size:      uint32(len(payload)),
		SHA256:    sha256.Sum256(payload),
	}

	var network bytes.Buffer

	// Posiljalac
	err := WriteFrame(
		&network,
		header,
		bytes.NewReader(payload),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Primalac
	receivedHeader, err := ReadHeader(&network)
	if err != nil {
		t.Fatal(err)
	}

	var receivedPayload bytes.Buffer

	err = ReadPayloadVerified(
		&network,
		receivedHeader,
		&receivedPayload,
		sha256.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if receivedHeader != header {
		t.Fatal("header mismatch")
	}

	if !bytes.Equal(receivedPayload.Bytes(), payload) {
		t.Fatal("payload mismatch")
	}

	t.Log("Binary message transfer successful")
}

func TestFSWCommand(t *testing.T) {
	header := Header{
		Version:   Version,
		Type:      TypeFSWStart,
		Sender:    testSender,
		Timestamp: time.Now().UnixMilli(),
		Algorithm: AlgorithmNone,
		Size:      0,
	}

	var network bytes.Buffer

	if err := WriteFrame(&network, header, nil); err != nil {
		t.Fatal(err)
	}

	received, err := ReadHeader(&network)
	if err != nil {
		t.Fatal(err)
	}

	if received.Type != TypeFSWStart {
		t.Fatal("incorrect FSW command")
	}

	if received.Size != 0 {
		t.Fatal("FSW must not have a payload")
	}

	t.Log("FSW command parsed successfully")
}

func TestCorruptedPayload(t *testing.T) {
	original := []byte("original")
	modified := []byte("modified")

	header := Header{
		Version:   Version,
		Type:      TypeMessage,
		Sender:    testSender,
		Timestamp: time.Now().UnixMilli(),
		Algorithm: AlgorithmLEACTR,
		Size:      uint32(len(modified)),
		SHA256:    sha256.Sum256(original),
	}

	var network bytes.Buffer

	if err := WriteFrame(
		&network,
		header,
		bytes.NewReader(modified),
	); err != nil {
		t.Fatal(err)
	}

	receivedHeader, err := ReadHeader(&network)
	if err != nil {
		t.Fatal(err)
	}

	err = ReadPayloadVerified(
		&network,
		receivedHeader,
		io.Discard,
		sha256.New(),
	)

	if err == nil {
		t.Fatal("corruption was not detected")
	}

	t.Log("Corrupted payload detected successfully")
}
