package network

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/jankibananki/secure-messenger/internal/protocol"
)

func TestTCPBinaryFrame(t *testing.T) {
	// port 0 operativni sistem bira
	listener, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	//nesifrovani test bajtovi
	payload := []byte("test transport bytes")

	results := make(chan error, 1)

	//simulacija primaoca
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			results <- err
			return
		}
		defer conn.Close()

		_ = conn.SetDeadline(
			time.Now().Add(5 * time.Second),
		)

		header, err := protocol.ReadHeader(conn)
		if err != nil {
			results <- err
			return
		}

		if header.Type != protocol.TypeMessage {
			results <- fmt.Errorf("unexpected packet type")
			return
		}

		var received bytes.Buffer

		err = protocol.ReadPayloadVerified(
			conn,
			header,
			&received,
			sha256.New(),
		)
		if err != nil {
			results <- err
			return
		}

		if !bytes.Equal(received.Bytes(), payload) {
			results <- fmt.Errorf("payload mismatch")
			return
		}

		results <- nil
	}()

	// Simulacija posiljaoca.
	conn, err := Connect(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(
		time.Now().Add(5 * time.Second),
	)

	header := protocol.Header{
		Version:   protocol.Version,
		Type:      protocol.TypeMessage,
		Sender:    "550e8400-e29b-41d4-a716-446655440000",
		Timestamp: time.Now().UnixMilli(),
		Algorithm: protocol.AlgorithmTEA,
		Size:      uint32(len(payload)),
		SHA256:    sha256.Sum256(payload),
	}

	err = protocol.WriteFrame(
		conn,
		header,
		bytes.NewReader(payload),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Proveravamo rezultat prijema.
	select {
	case err := <-results:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("TCP receive timeout")
	}

	t.Log("TCP binary frame transferred successfully")
}
