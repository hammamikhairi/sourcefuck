package parser

import (
	. "LanguageFuck/Encrypter"
	. "LanguageFuck/Lexer"
	. "LanguageFuck/Types"
	"fmt"
)

type Parser struct {
	Tokens *[]*Token
	Swap   map[string]string
	Enc    *Encrypter
}

func ParserInit(tokens *[]*Token, key int) *Parser {
	return &Parser{
		tokens,
		make(map[string]string),
		EncrypterInit(key),
	}
}

func (pr *Parser) Parse(l *Lexer, decrypt bool) {
	var enc string
	for _, token := range *pr.Tokens {
		token_content := l.GetTokenContent(token)
		if _, ok := pr.Swap[token_content]; ok {
			continue
		}
		if token.Kind == TOKEN_SYMBOL {
			if decrypt {
				enc = pr.Enc.Decrypt(token_content)
			} else {
				enc = pr.Enc.Encrypt(token_content)
			}
			pr.Swap[token_content] = enc
		}
	}
	fmt.Println(pr.Swap)
}