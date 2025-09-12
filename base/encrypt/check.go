package encrypt

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateCheckKEy
//||------------------------------------------------------------------------------------------------||

func GenerateCheckKey(privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("invalid private key PEM format")
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
