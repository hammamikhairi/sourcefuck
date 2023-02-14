package utils

import (
	"log"
	"unicode"
)

func IsSpace(char string) bool {
	return char == " "
}

func IsAlpha(char string) bool {
	return (char[0] >= 'a' && char[0] <= 'z') || (char[0] >= 'A' && char[0] <= 'Z')
}

func IsSymbolChar(char string) bool {
	return unicode.IsNumber(rune(char[0])) || IsAlpha(char) || char == "_"
}

func Assert(cond bool, errorM string) {
	if cond == false {
		log.Fatal(errorM)
	}
}