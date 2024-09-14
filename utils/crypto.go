package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"
)

// Generate a 32-byte key using SHA-256 hash of the given passphrase
func generateKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

// Encrypt a plaintext string using AES-GCM with the given key
func Encrypt(plaintext string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}

	passphrase := config.JWTSecret
	key := generateKey(passphrase)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Use Galois Counter Mode (GCM) for encryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce (number used once) for encryption
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return the ciphertext as a base64-encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt a base64-encoded ciphertext string using AES-GCM with the given key
func Decrypt(ciphertextBase64 string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}

	passphrase := config.JWTSecret
	key := generateKey(passphrase)

	// Decode the base64-encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Use Galois Counter Mode (GCM) for decryption
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
