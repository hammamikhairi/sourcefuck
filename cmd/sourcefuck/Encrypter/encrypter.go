package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"strings"
)

type Encrypter struct {
	key []byte
}

func EncrypterInit(key []byte) *Encrypter {
	return &Encrypter{
		key,
	}
}

var (
	buffer strings.Builder
)

func customBase32Encode(data []byte) string {
	encoder := base32.NewEncoder(base32.StdEncoding.WithPadding('_'), &buffer)
	encoder.Write(data)
	encoder.Close()

	defer buffer.Reset()
	return buffer.String()
}

func customBase32Decode(encoded string) ([]byte, error) {
	reader := strings.NewReader(encoded)
	decoder := base32.NewDecoder(base32.StdEncoding.WithPadding('_'), reader)

	_, err := io.Copy(&buffer, decoder)
	if err != nil {
		println(encoded)
		fmt.Println("Decoding error:", err)
		return nil, err
	}
	defer buffer.Reset()
	return []byte(buffer.String()), nil
}

func (c *Encrypter) Encrypt(text string) string {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return ""
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return customBase32Encode(ciphertext)
}

func (c *Encrypter) Decrypt(ciphertext string) string {
	ciphertextBytes, err := customBase32Decode(ciphertext)
	if err != nil {
		return ""
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return ""
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	stream.XORKeyStream(plaintext, ciphertextBytes)

	return string(plaintext)
}
