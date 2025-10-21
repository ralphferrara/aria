package sign

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

//||------------------------------------------------------------------------------------------------||
//|| Verify: Verify a message using a PEM-encoded public key
//||------------------------------------------------------------------------------------------------||

func Verify(message, signature []byte, publicKeyPEM string) (bool, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, errors.New("invalid public key PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	pub, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("not an ECDSA public key")
	}

	// Hash and verify
	hash := sha256.Sum256(message)
	okVerify := ecdsa.VerifyASN1(pub, hash[:], signature)
	return okVerify, nil
}
