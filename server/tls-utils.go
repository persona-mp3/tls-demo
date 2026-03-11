package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

// Loads a server.crt and server.key file generated
// using the `gen-certs.sh`. It looks in the root directory
// for server.crt and server.key file
func LoadTlsConfig() (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return nil, fmt.Errorf("couldn't open server.crt file because: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	return tlsConfig, nil
}

func InitHandshake(conn net.Conn) bool {
	log.SetFlags(log.Lshortfile)
	tlsConn, valid := conn.(*tls.Conn)
	if !valid {
		log.Printf("unexpectedly got a non-tls connection\n")
		return false
	}

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("tls-handshake wasn't succesfull because: %s\n", err)
		return false
	}

	return true
}
