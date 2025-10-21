package encrypt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

//||------------------------------------------------------------------------------------------------||
//|| VerifyWithPublicKey: Verifies a digital signature using a PKIX PEM public key
//||------------------------------------------------------------------------------------------------||

func VerifyWithPublicKey(message []byte, signatureB64 string, publicKeyPEM string) error {
	// Decode PEM
	pemBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pemBlock == nil || pemBlock.Type != "PUBLIC KEY" {
		return errors.New("invalid public key PEM format (expected PKIX)")
	}

	// Parse key
	keyAny, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return errors.New("invalid public key")
	}
	publicKey, ok := keyAny.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	// Hash message
	hashed := sha256.Sum256(message)

	// Decode signature
	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return errors.New("invalid base64 signature")
	}

	// Verify
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return errors.New("invalid signature")
	}

	return nil
}
