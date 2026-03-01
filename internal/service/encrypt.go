package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptAndSplitFile 对文件加密并分片
func EncryptAndSplitFile(file io.Reader, chunkSize int) ([][]byte, []byte, []byte, error) {
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		return nil, nil, nil, err
	}
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	plainData, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, nil, err
	}
	// PKCS7填充
	padLen := aes.BlockSize - len(plainData)%aes.BlockSize
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	plainData = append(plainData, pad...)

	encrypted := make([]byte, len(plainData))
	mode.CryptBlocks(encrypted, plainData)

	// 分片
	var chunks [][]byte
	for i := 0; i < len(encrypted); i += chunkSize {
		end := i + chunkSize
		if end > len(encrypted) {
			end = len(encrypted)
		}
		chunks = append(chunks, encrypted[i:end])
	}
	return chunks, key, iv, nil
}

// MergeAndDecryptFile 合并分片并解密
func MergeAndDecryptFile(chunks [][]byte, key, iv []byte) ([]byte, error) {
	var encrypted []byte
	for _, chunk := range chunks {
		encrypted = append(encrypted, chunk...)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)
	// 去除PKCS7填充
	padLen := int(decrypted[len(decrypted)-1])
	return decrypted[:len(decrypted)-padLen], nil
}
