package validate

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

//||------------------------------------------------------------------------------------------------||
//|| ValidateKeyPair
//||------------------------------------------------------------------------------------------------||

func ValidateKeyPair(privateKeyPEM, publicKeyPEM string) error {
	// Private (PKCS#8, "PRIVATE KEY")
	privBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privBlock == nil || privBlock.Type != "PRIVATE KEY" {
		return errors.New("invalid private key PEM format (expecting PKCS#8 'PRIVATE KEY')")
	}
	privAny, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}
	priv, ok := privAny.(*rsa.PrivateKey)
	if !ok {
		return errors.New("private key is not RSA")
	}

	// Public (PKIX/SPKI, "PUBLIC KEY")
	pubBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return errors.New("invalid public key PEM format (expecting PKIX 'PUBLIC KEY')")
	}
	pubAny, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PKIX public key: %w", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return errors.New("public key is not RSA")
	}

	// Match modulus & exponent
	if priv.PublicKey.N.Cmp(pub.N) != 0 || priv.PublicKey.E != pub.E {
		return errors.New("public key does not match the private key")
	}
	return nil
}
