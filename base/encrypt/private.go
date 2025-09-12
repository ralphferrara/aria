package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| Decrypts data using RSA public key in PEM format
//||------------------------------------------------------------------------------------------------||

func DecryptWithPrivateKey(ciphertext []byte, privateKeyPEM string) ([]byte, error) {
	//||------------------------------------------------------------------------------------------------||
	//|| Split base64 parts
	//||------------------------------------------------------------------------------------------------||

	parts := strings.SplitN(string(ciphertext), ".", 3)
	if len(parts) != 3 {
		return nil, errors.New("invalid ciphertext format for hybrid decrypt")
	}

	encryptedAESKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	cipherData, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Decrypt AES Key with RSA Private Key
	//||------------------------------------------------------------------------------------------------||

	pemBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if pemBlock == nil || pemBlock.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM format")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedAESKey)
	if err != nil {
		return nil, err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Decrypt Data with AES-GCM
	//||------------------------------------------------------------------------------------------------||

	blockAES, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(blockAES)
	if err != nil {
		return nil, err
	}
	plainData, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plainData, nil
}
