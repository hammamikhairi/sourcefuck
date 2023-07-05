package types

// TEMPORARILY FOR GO ONLY
type TokenKind uint8

const (
	TOKEN_INVALID TokenKind = iota
	TOKEN_PREPROC
	TOKEN_SYMBOL
	TOKEN_KEYWORD
	TOKEN_TYPE
	TOKEN_LIB
	TOKEN_IMPORTED
	TOKEN_COMMENT
	TOKEN_STRING
	TOKEN_TAB
	TOKEN_END
)

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
	case TOKEN_TYPE:
		return "type"
	case TOKEN_LIB:
		return "lib"
	case TOKEN_IMPORTED:
		return "imported type"
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
