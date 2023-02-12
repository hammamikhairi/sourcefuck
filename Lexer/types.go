package lexer

// TEMPORARILY FOR GO ONLY
var KEYWORDS []string = []string{
	"break ", "default ", "func ", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

type TokenKind uint8

const (
	TOKEN_INVALID TokenKind = iota
	TOKEN_PREPROC
	TOKEN_SYMBOL
	TOKEN_KEYWORD
	TOKEN_COMMENT
	TOKEN_STRING
	TOKEN_TAB
	TOKEN_END
)

// For Token position [inMain, col]
type Vec2i struct {
	X, Line int
}

type Token struct {
	Kind TokenKind
	Addr Vec2i
	Len  int
}
