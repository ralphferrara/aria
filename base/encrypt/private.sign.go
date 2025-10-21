package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

//||------------------------------------------------------------------------------------------------||
//|| SignWithPrivateKey: Creates a digital signature using a PKCS#8 PEM private key
//||------------------------------------------------------------------------------------------------||

func SignWithPrivateKey(message []byte, privateKeyPEM string) (string, error) {
	// Decode PEM
	pemBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
		return "", errors.New("invalid private key PEM format (expected PKCS#8)")
	}

	// Parse key
	keyAny, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return "", err
	}
	privateKey, ok := keyAny.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// Hash message
	hashed := sha256.Sum256(message)

	// Sign
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
