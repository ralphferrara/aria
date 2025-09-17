package encrypt

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateCheckKey: Generates a hash from a PKCS#8 private key PEM
//||------------------------------------------------------------------------------------------------||

func GenerateCheckKey(privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("invalid private key PEM format (expecting PKCS#8)")
	}

	// Parse to confirm it’s valid PKCS#8
	_, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	hash := sha256.Sum256([]byte(privateKeyPEM))
	checkKey := hex.EncodeToString(hash[:])

	return checkKey, nil
}

//||------------------------------------------------------------------------------------------------||
//|| CheckPrivateKey
//||------------------------------------------------------------------------------------------------||

func CheckPrivateKey(privateKeyPEM string, checkKey string) error {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return errors.New("invalid private key PEM format")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	if privateKey.N.String() != checkKey {
		return errors.New("private key does not match the provided key")
	}

	hash := sha256.Sum256([]byte(privateKeyPEM))
	if hex.EncodeToString(hash[:]) != checkKey {
		return errors.New("private key does not match the provided key")
	}

	return nil
}
