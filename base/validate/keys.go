package validate

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//||------------------------------------------------------------------------------------------------||
//|| ValidateKeyPair
//||------------------------------------------------------------------------------------------------||

func ValidateKeyPair(privateKeyPEM, publicKeyPEM string) error {
	// Decode and parse the private key
	privBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return errors.New("invalid private key PEM format")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return errors.New("failed to parse RSA private key: " + err.Error())
	}

	// Decode and parse the public key
	pubBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pubBlock == nil || pubBlock.Type != "RSA PUBLIC KEY" {
		return errors.New("invalid public key PEM format")
	}
	pubKeyIface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return errors.New("failed to parse RSA public key: " + err.Error())
	}
	pubKey, ok := pubKeyIface.(*rsa.PublicKey)
	if !ok {
		return errors.New("provided public key is not an RSA key")
	}

	// Ensure public key matches the private key
	if pubKey.N.Cmp(privateKey.PublicKey.N) != 0 || pubKey.E != privateKey.PublicKey.E {
		return errors.New("public key does not match the private key")
	}

	return nil
}
