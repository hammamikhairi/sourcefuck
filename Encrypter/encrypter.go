package encrypter

import (
	// "bytes"
	// "fmt"
	// "unicode"
	// "crypto/aes"
	// "crypto/cipher"
	// "unicode"
	. "github.com/hammamikhairi/langfuck/Utils"
)

// TODO : for now i'm testing with caesar cipher, will change later

type Encrypter struct {
	key int
}

// The count of unique letters in the English alphabet
const ALPHA_LEN byte = 26

func (c *Encrypter) Encrypt(plainText string) string {
	cipherText := ""
	for i := 0; i < len(plainText); i++ {
		char := plainText[i]

		// Apply Caesar cipher, and maintain case of letters
		if IsAlpha(char) {
			A := byte('A')
			if IsLower(char) {
				A = 'a'
			}
			char = (char-A+byte(c.key))%ALPHA_LEN + A
		}
		cipherText += string(char)
	}
	return cipherText
}

func (c *Encrypter) Decrypt(cipherText string) string {
	plainText := ""
	for i := 0; i < len(cipherText); i++ {
		char := cipherText[i]

		// Apply Caesar cipher, and maintain case of letters
		if IsAlpha(char) {
			A := byte('A')
			if IsLower(char) {
				A = 'a'
			}
			char = (char-A+ALPHA_LEN-byte(c.key))%ALPHA_LEN + A
		}
		plainText += string(char)
	}
	return plainText
}

func EncrypterInit(pwd int) *Encrypter {
	return &Encrypter{
		key: pwd,
	}
}
