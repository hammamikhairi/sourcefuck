package types

var KEYWORDS []string = []string{
	"break ", "default ", "func ", "interface", "select",
	"case", "defer", "go", "map", "struct", "func", "main",
	"chan", "else", "goto", "package", "switch", "init",
	"const", "fallthrough", "if", "range", "type", "len",
	"continue", "for", "import", "return", "var", "make",
	"panic", "false", "true",
}

var TYPES []string = []string{
	"bool", "string", "int", "int8", "int16", "int32",
	"int64", "uint", "uint8", "uint16", "uint32", "uint64",
	"uintptr", "float32", "float64", "complex64", "complex128",
	"byte", "rune", "nil",
}

var LIBRARIES []string = []string{
	"fmt", "os", "log", "net", "filepath", "ftp",
}

// For Token position [inMain, col]
type Vec2i struct {
	X, Line, Origin int
}

type Token struct {
	Kind TokenKind
	Addr Vec2i
	Len  int
}
