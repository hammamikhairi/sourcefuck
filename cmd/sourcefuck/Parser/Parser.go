package parser

import (
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Encrypter"
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Lexer"
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Types"
)

type Parser struct {
	Tokens *[]*Token
	Swap   map[string]string
	Enc    *Encrypter
}

func ParserInit(tokens *[]*Token, key []byte) *Parser {
	return &Parser{
		tokens,
		make(map[string]string),
		EncrypterInit(key),
	}
}

func (pr *Parser) Parse(l *Lexer, decrypt bool) {
	var (
		enc        string
		token      *Token
		tokens     []*Token = (*pr.Tokens)
		tokens_len int      = len(*pr.Tokens)
	)
	for i := 0; i < tokens_len; i++ {
		token = tokens[i]
		token_content := l.GetTokenContent(token)

		if _, ok := pr.Swap[token_content]; token.Kind == TOKEN_IMPORTED && !ok {
			if decrypt {
				enc = pr.Enc.Decrypt(token_content[1:])
			} else {
				enc = string(token_content[0]) + pr.Enc.Encrypt(token_content)
			}
			for i+1 < tokens_len && l.GetTokenContent(tokens[i+1]) == "." {
				i += 2
			}
			pr.Swap[token_content] = enc
		}

		if token.Kind == TOKEN_SYMBOL {
			if decrypt {
				enc = pr.Enc.Decrypt(token_content[1:])
			} else {
				enc = string(token_content[0]) + pr.Enc.Encrypt(token_content)
			}
			pr.Swap[token_content] = enc
		}
	}
}
