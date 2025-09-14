//||------------------------------------------------------------------------------------------------||
//|| KeyPair: Generate RSA Private/Public PEM-Encoded Key Pair
//||------------------------------------------------------------------------------------------------||

package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateKeyPair: Returns PEM-encoded private and public keys (RSA 4096-bit)
//||------------------------------------------------------------------------------------------------||

func GenerateKeyPair() (privateKeyPEM string, publicKeyPEM string, err error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Generate RSA Key (4096-bit)
	//||------------------------------------------------------------------------------------------------||

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key: %w", err)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Private Key (PEM)
	//||------------------------------------------------------------------------------------------------||

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Public Key (PEM)
	//||------------------------------------------------------------------------------------------------||

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(privPEM), string(pubPEM), nil
}
