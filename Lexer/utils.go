package lexer

func isSpace(char string) bool {
	return char == " "
}

func isAlpha(char string) bool {
	return (char[0] >= 'a' && char[0] <= 'z') || (char[0] >= 'A' && char[0] <= 'Z')
}

func GetTokenName(tk TokenKind) string {
	switch tk {
	case TOKEN_INVALID:
		return "invalid token"
	case TOKEN_PREPROC:
		return "preprocessor directive"
	case TOKEN_SYMBOL:
		return "symbol"
	case TOKEN_KEYWORD:
		return "keyword"
	case TOKEN_COMMENT:
		return "comment"
	case TOKEN_STRING:
		return "string"
	case TOKEN_TAB:
		return "tabulation"
	case TOKEN_END:
		return "EOF"
	}
	return "UNREACHABLE"
}