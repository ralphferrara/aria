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

func GenerateBIP39Keys(words []string) (privateKeyPEM string, publicKeyPEM string, err error) {

	const keySize = 2048

	//||------------------------------------------------------------------------------------------------||
	//|| Normalize Words
	//||------------------------------------------------------------------------------------------------||

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

	//||------------------------------------------------------------------------------------------------||
	//|| Create Seed
	//||------------------------------------------------------------------------------------------------||

	mnemonic := strings.Join(clean, " ")
	seed := sha256.Sum256([]byte(mnemonic))

	//||------------------------------------------------------------------------------------------------||
	//|| Generate RSA Keys
	//||------------------------------------------------------------------------------------------------||

	reader := NewDeterministicReader(seed[:])
	privKey, err := rsa.GenerateKey(reader, keySize)
	if err != nil {
		return "", "", err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Private Key
	//||------------------------------------------------------------------------------------------------||

	privDER := x509.MarshalPKCS1PrivateKey(privKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privateKeyPEM = string(pem.EncodeToMemory(&privBlock))

	//||------------------------------------------------------------------------------------------------||
	//|| Encode Public Key
	//||------------------------------------------------------------------------------------------------||

	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}
	publicKeyPEM = string(pem.EncodeToMemory(&pubBlock))

	return privateKeyPEM, publicKeyPEM, nil
}
