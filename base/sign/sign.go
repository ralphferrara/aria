package sign

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

//||------------------------------------------------------------------------------------------------||
//|| Sign: Sign a message using a PEM-encoded private key
//||------------------------------------------------------------------------------------------------||

func Sign(message []byte, privateKeyPEM string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an ECDSA private key")
	}

	// Hash and sign
	fmt.Printf("PLAIN: %x\n", message)

	return ecdsa.SignASN1(rand.Reader, priv, message)
}
