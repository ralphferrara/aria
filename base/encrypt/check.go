package encrypt

import (
	"crypto/rsa"
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

	// Parse to confirm it’s valid PKCS#8 RSA
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}
	if _, ok := key.(*rsa.PrivateKey); !ok {
		return "", errors.New("not an RSA private key")
	}

	// Hash the full PEM text
	hash := sha256.Sum256([]byte(privateKeyPEM))
	checkKey := hex.EncodeToString(hash[:])

	return checkKey, nil
}

//||------------------------------------------------------------------------------------------------||
//|| CheckPrivateKey: Validates that the PEM matches the expected hash
//||------------------------------------------------------------------------------------------------||

func CheckPrivateKey(privateKeyPEM string, checkKey string) error {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return errors.New("invalid private key PEM format (expecting PKCS#8)")
	}

	// Parse to confirm it’s valid PKCS#8 RSA
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}
	if _, ok := key.(*rsa.PrivateKey); !ok {
		return errors.New("not an RSA private key")
	}

	// Verify hash
	hash := sha256.Sum256([]byte(privateKeyPEM))
	if hex.EncodeToString(hash[:]) != checkKey {
		return errors.New("private key does not match the provided check key")
	}

	return nil
}

//||------------------------------------------------------------------------------------------------||
//|| CheckKeyPair: Verifies that the given PEM-encoded private and public keys form a valid pair
//||------------------------------------------------------------------------------------------------||

func CheckKeyPair(privateKeyPEM, publicKeyPEM string) error {
	// Parse private key PEM
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("invalid private key PEM")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not RSA")
	}

	// Parse public key PEM
	pubBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	// Compare modulus and exponent
	if rsaPriv.PublicKey.N.Cmp(rsaPub.N) != 0 || rsaPriv.PublicKey.E != rsaPub.E {
		return fmt.Errorf("keys do not match")
	}

	return nil
}
