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
	"io"
)

//||------------------------------------------------------------------------------------------------||
//|| Encrypts data using RSA public key in PEM format
//||------------------------------------------------------------------------------------------------||

func EncryptWithPublicKey(data []byte, publicKeyPEM string) ([]byte, error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Generate AES-256 Key
	//||------------------------------------------------------------------------------------------------||

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Encrypt Data with AES-GCM
	//||------------------------------------------------------------------------------------------------||

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherData := aesGCM.Seal(nil, nonce, data, nil)

	//||------------------------------------------------------------------------------------------------||
	//|| Encrypt AES Key with RSA Public Key
	//||------------------------------------------------------------------------------------------------||

	pemBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pemBlock == nil || pemBlock.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid public key PEM format")
	}
	pub, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	encryptedAESKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return nil, err
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Join as base64 strings: encryptedAESKey.nonce.cipherData
	//||------------------------------------------------------------------------------------------------||

	final := base64.StdEncoding.EncodeToString(encryptedAESKey) + "." +
		base64.StdEncoding.EncodeToString(nonce) + "." +
		base64.StdEncoding.EncodeToString(cipherData)

	return []byte(final), nil
}
