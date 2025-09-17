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
//|| GenerateKeyPair: Returns PEM-encoded private and public keys (RSA 4096-bit, PKCS#8 format)
//||------------------------------------------------------------------------------------------------||

func GenerateKeyPair() (privateKeyPEM string, publicKeyPEM string, err error) {
	// Generate RSA Key (4096-bit)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Encode Private Key (PKCS#8)
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY", // ✅ PKCS#8
		Bytes: privBytes,
	})

	// Encode Public Key (PKIX/SPKI)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY", // ✅ standard PKIX/SPKI
		Bytes: pubBytes,
	})

	return string(privPEM), string(pubPEM), nil
}
