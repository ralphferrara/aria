package encrypt

import (
	"bytes"
	"testing"
)

//||------------------------------------------------------------------------------------------------||
//|| Test Key Generation, Hashing, Encryption & Decryption
//||------------------------------------------------------------------------------------------------||

func TestEncryptDecryptCycle(t *testing.T) {
	// Generate key pair
	privPEM, pubPEM, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	if !bytes.Contains([]byte(privPEM), []byte("BEGIN PRIVATE KEY")) {
		t.Errorf("expected PKCS#8 private key PEM, got:\n%s", privPEM)
	}
	if !bytes.Contains([]byte(pubPEM), []byte("BEGIN PUBLIC KEY")) {
		t.Errorf("expected PKIX public key PEM, got:\n%s", pubPEM)
	}

	// Generate check key
	checkKey, err := GenerateCheckKey(privPEM)
	if err != nil {
		t.Fatalf("failed to generate check key: %v", err)
	}

	// Validate private key matches check key
	if err := CheckPrivateKey(privPEM, checkKey); err != nil {
		t.Fatalf("check private key failed: %v", err)
	}

	// Test encryption
	original := []byte("complyage-encryption-test")
	ciphertext, err := EncryptWithPublicKey(original, pubPEM)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// Test decryption
	plaintext, err := DecryptWithPrivateKey(ciphertext, privPEM)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	// Verify roundtrip
	if !bytes.Equal(original, plaintext) {
		t.Errorf("decrypted text mismatch.\nGot:  %s\nWant: %s", plaintext, original)
	}
}
