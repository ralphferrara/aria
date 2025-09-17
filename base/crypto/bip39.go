package crypto

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateBIP39Keys: Accepts dynamic number of words
//||------------------------------------------------------------------------------------------------||

//||------------------------------------------------------------------------------------------------||
//|| GenerateBIP39Keys: Accepts dynamic number of words, returns PKCS#8 private key + PKIX public key
//||------------------------------------------------------------------------------------------------||

func GenerateBIP39Keys(words []string) (privateKeyPEM string, publicKeyPEM string, err error) {
	const keySize = 2048

	// Normalize words
	var clean []string
	for _, w := range words {
		w = strings.ToLower(strings.TrimSpace(w))
		if w != "" {
			clean = append(clean, w)
		}
	}
	if len(clean) == 0 {
		return "", "", errors.New("no valid BIP39 words provided")
	}

	// Create seed
	mnemonic := strings.Join(clean, " ")
	seed := sha256.Sum256([]byte(mnemonic))

	// Generate RSA key deterministically
	reader := NewDeterministicReader(seed[:])
	privKey, err := rsa.GenerateKey(reader, keySize)
	if err != nil {
		return "", "", err
	}

	// Encode Private Key (PKCS#8)
	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}
	privBlock := pem.Block{
		Type:  "PRIVATE KEY", // ✅ PKCS#8
		Bytes: privDER,
	}
	privateKeyPEM = string(pem.EncodeToMemory(&privBlock))

	// Encode Public Key (PKIX/SPKI)
	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY", // ✅ PKIX/SPKI
		Bytes: pubDER,
	}
	publicKeyPEM = string(pem.EncodeToMemory(&pubBlock))

	return privateKeyPEM, publicKeyPEM, nil
}
