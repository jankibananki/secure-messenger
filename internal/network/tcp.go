package network

import (
	"fmt"
	"net"
	"time"
)

const connectTimeout = 5 * time.Second

// pokretanja tcp servera koji osluskuje zadatu adresu
func Listen(address string) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf(
			"start TCP listener: %w", err,
		)
	}
	return listener, nil
}

// povezivanje sa drugom instancom aplikacije
func Connect(address string) (net.Conn, error) {
	conn, err := net.DialTimeout(
		"tcp",
		address,
		connectTimeout,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"connect to peer: %w", err,
		)
	}

	return conn, nil
}
