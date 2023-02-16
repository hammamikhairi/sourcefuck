package encrypter

import (
	// "bytes"
	// "fmt"
	// "unicode"
	// "crypto/aes"
	// "crypto/cipher"
	// "unicode"
	. "LanguageFuck/Utils"
)

// TODO : for now i'm testing with caesar cipher, will change later

type Encrypter struct {
	key int
}

func (c *Encrypter) Encrypt(plainText string) string {
	cipherText := ""
	for i := 0; i < len(plainText); i++ {
		char := plainText[i]
		if i == 0 && IsAlpha(char) {
			// Maintain case of first letter
			if IsUpper(char) {
				char = byte(int(char-'A'+byte(c.key))%26 + 'A')
			} else { // Lower
				char = byte(int(char-'a'+byte(c.key))%26 + 'a')
			}
		} else {
			// Apply Caesar cipher
			if IsUpper(char) {
				char = byte(int(char-'A'+byte(c.key))%26 + 'A')
			} else if IsLower(char) {
				char = byte(int(char-'a'+byte(c.key))%26 + 'a')
			}
		}
		cipherText += string(char)
	}
	return cipherText
}

func (c *Encrypter) Decrypt(cipherText string) string {
	// The count of unique letters in the English alphabet
	const ALPHA_LEN = 26

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
