package bip39

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| GenerateBIP39Keys: Accepts BIP39WordList, returns PKCS#8 private key + PKIX public key
//||------------------------------------------------------------------------------------------------||

func GenerateBIP39Keys(words BIP39WordList) (privateKeyPEM string, publicKeyPEM string, err error) {
	const keySize = 2048

	//||------------------------------------------------------------------------------------------------||
	//|| Normalize Words
	//||------------------------------------------------------------------------------------------------||

	var clean []string
	for i := 0; i < len(words); i++ {
		w := strings.ToLower(strings.TrimSpace(words[i]))
		if w != "" {
			clean = append(clean, w)
		}
	}
	if len(clean) == 0 {
		return "", "", errors.New("no valid BIP39 words provided")
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Create Seed from Mnemonic
	//||------------------------------------------------------------------------------------------------||

	mnemonic := strings.Join(clean, " ")
	seed := sha256.Sum256([]byte(mnemonic))

	//||------------------------------------------------------------------------------------------------||
	//|| Generate RSA Key Deterministically
	//||------------------------------------------------------------------------------------------------||

	reader := NewDeterministicReader(seed[:])
	privKey, err := rsa.GenerateKey(reader, keySize)
	if err != nil {
		return "", "", err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Private Key (PKCS#8)
	//||------------------------------------------------------------------------------------------------||

	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}
	privBlock := pem.Block{
		Type:  "PRIVATE KEY", // ✅ PKCS#8
		Bytes: privDER,
	}
	privateKeyPEM = string(pem.EncodeToMemory(&privBlock))

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Public Key (PKIX/SPKI)
	//||------------------------------------------------------------------------------------------------||

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
