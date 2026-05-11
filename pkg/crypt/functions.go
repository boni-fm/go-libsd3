package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// generateKey menghasilkan 32-byte key menggunakan SHA256 (sama seperti C#)
func generateKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// generateIV menghasilkan 16-byte IV menggunakan MD5 (sama seperti C#)
func generateIV(key string) []byte {
	hash := md5.Sum([]byte(key))
	return hash[:]
}

// EncryptJSON mengenkripsi struct/map apapun menjadi encrypted base64 string
func EncryptJSON(v any, key string) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return Encrypt(string(jsonBytes), key)
}

// DecryptJSON mendekripsi encrypted base64 string menjadi struct yang diinginkan
func DecryptJSON(cipherText, key string, v any) error {
	plainText, err := Decrypt(cipherText, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	if err := json.Unmarshal([]byte(plainText), v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Encrypt mengenkripsi plainText menggunakan AES-CBC
func Encrypt(plainText, key string) (string, error) {
	keyBytes := generateKey(key)
	iv := generateIV(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// PKCS7 padding
	plainBytes := []byte(plainText)
	blockSize := block.BlockSize()
	padding := blockSize - len(plainBytes)%blockSize
	padded := make([]byte, len(plainBytes)+padding)
	copy(padded, plainBytes)
	for i := len(plainBytes); i < len(padded); i++ {
		padded[i] = byte(padding)
	}

	encrypted := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, padded)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt mendekripsi cipherText menggunakan AES-CBC
func Decrypt(cipherText, key string) (string, error) {
	keyBytes := generateKey(key)
	iv := generateIV(key)

	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	if len(cipherBytes)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherBytes, cipherBytes)

	// Remove PKCS7 padding
	padding := int(cipherBytes[len(cipherBytes)-1])
	if padding > block.BlockSize() || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}
	decrypted := cipherBytes[:len(cipherBytes)-padding]

	return string(decrypted), nil
}
