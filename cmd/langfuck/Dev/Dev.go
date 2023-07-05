package dev

import (
	"fmt"

	. "github.com/hammamikhairi/langfuck/cmd/langfuck/Lexer"
	. "github.com/hammamikhairi/langfuck/cmd/langfuck/Types"
)

func DebugTokens(l *Lexer, tokens *[]*Token) {
	for _, token := range *tokens {

		if token.Kind == TOKEN_END {
			fmt.Println(GetTokenName(token.Kind))
			continue
		}

		fmt.Println(
			l.GetTokenContent(token),
			"\t<<", token.Addr.Line, token.Addr.X, ">>",
			"\t<<", GetTokenName(token.Kind), ">>",
		)
	}
}
