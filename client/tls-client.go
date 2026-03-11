package client

import (
	"crypto/x509"
	"fmt"
	"os"
)

func LoadServerCerts() (*x509.CertPool, error) {
	pem, err := os.ReadFile("./server.crt")
	if err != nil {
		return nil, fmt.Errorf("error loading server.crt: %w", err)
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("valid pem file? run ./gen-tls.sh file: %w", err)
	}

	fmt.Println("successfully loaded pem file")
	return cp, nil
}
