package authenticator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// generateKey menghasilkan 32-byte key menggunakan SHA256 (sama seperti C#)
func generateKey() []byte {
	hash := sha256.Sum256([]byte("SD3Indomaret"))
	return hash[:]
}

// generateIV menghasilkan 16-byte IV menggunakan MD5 (sama seperti C#)
func generateIV() []byte {
	hash := md5.Sum([]byte("SD3Indomaret"))
	return hash[:]
}

// Encrypt mengenkripsi plainText menggunakan AES-CBC
func Encrypt(plainText string) (string, error) {
	keyBytes := generateKey()
	iv := generateIV()

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
func Decrypt(cipherText string) (string, error) {
	keyBytes := generateKey()
	iv := generateIV()

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
