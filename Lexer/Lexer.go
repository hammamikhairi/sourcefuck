package lexer

import (
	. "LanguageFuck/Types"
	. "LanguageFuck/Utils"
)

// for imported libs
var CurrentLineTokens []*Token
var ImportedSymbs map[string]uint8 = make(map[string]uint8)

type Lexer struct {
	Content      string
	Content_len  int
	Cursor       int
	Line         int
	LineStart    int
	KeywordsTree *map[string]uint8
}

func LexerInit(content string, tree *map[string]uint8) *Lexer {
	return &Lexer{
		Content:      content,
		Content_len:  len(content),
		Cursor:       0,
		Line:         0,
		LineStart:    0,
		KeywordsTree: tree,
	}

}

func (l *Lexer) ChopChar(len int) {
	for i := 0; i < len; i++ {
		current := string(l.Content[l.Cursor])
		l.Cursor++
		if current == "\n" {
			l.Line++
			l.LineStart = l.Cursor
		}
	}
}

func (l *Lexer) Trim() {
	for l.Cursor < l.Content_len-1 && (IsSpace(string(l.Content[l.Cursor])) || l.getCharAt(l.Cursor) == '\n') {
		l.ChopChar(1)
	}
}

func (l *Lexer) getCharAt(pos int) rune {
	Assert(pos < l.Content_len, "INVALID POS")
	return rune(l.Content[pos])
}

func (l *Lexer) startsWith(prefix string) bool {

	if len(prefix) == 1 {
		return l.getCharAt(l.Cursor) == rune(prefix[0])
	}

	for i := 0; i < len(prefix); i++ {
		if prefix[i] != l.Content[l.Cursor+1] {
			return false
		}
	}
	return true
}

func (l *Lexer) NextToken() *Token {
	l.Trim()
	l.getCharAt(l.Cursor)
	// fmt.Println(CurrentLineTokens)

	token := &Token{}
	token.Addr = Vec2i{X: l.Cursor, Line: l.Line, Origin: l.LineStart}

	st := 0

	if l.Cursor >= l.Content_len {
		token.Kind = TOKEN_END
		token.Len = 1
		return token
	}

	if l.startsWith("\"") {
		token.Kind = TOKEN_STRING
		l.ChopChar(1)
		for l.Cursor < l.Content_len-1 {

			if l.getCharAt(l.Cursor) == '"' && l.getCharAt(l.Cursor-1) != '\\' {
				break
			}

			l.ChopChar(1)
			st++
		}
		st += 2

		l.ChopChar(1)
		token.Len = st
		return token
	}

	if l.startsWith("\t") {
		token.Kind = TOKEN_TAB
		l.ChopChar(1)
		st++
		for l.getCharAt(l.Cursor) == '\t' {
			l.ChopChar(1)
			st++
		}
		token.Len = st
		return token
	}

	if l.startsWith("/") {
		l.ChopChar(1)
		if !l.startsWith("/") {
			l.Cursor--
			l.ChopChar(1)
			token.Kind = TOKEN_SYMBOL
			token.Len = 1
			return token
		}
		token.Kind = TOKEN_COMMENT
		for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != '\n' {
			st++
			l.ChopChar(1)
		}
		if l.Cursor < l.Content_len {
			st++
			l.ChopChar(1)
		}
		token.Len = st
		return token
	}

	if IsAlpha(byte(l.getCharAt(l.Cursor))) {
		token.Kind = TOKEN_SYMBOL
		for l.Cursor < l.Content_len && IsSymbolChar(byte(l.getCharAt(l.Cursor))) {
			l.ChopChar(1)
			st++
		}
		token.Len = st

		// PREPROC
		if l.GetTokenContent(token) == "package" {
			token.Kind = TOKEN_PREPROC
			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != '\n' {
				st++
				l.ChopChar(1)
			}
			token.Len = st
			return token
		}

		// PREPROC
		if l.GetTokenContent(token) == "import" {
			token.Kind = TOKEN_PREPROC
			l.Trim()
			end := '\n'
			if l.startsWith("(") {
				end = ')'
			}

			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != end {
				st++
				l.ChopChar(1)
			}

			if end == ')' {
				st += 1
				l.ChopChar(1)
			}

			st++
			l.ChopChar(1)

			token.Len = st
			return token
		}

		// KEYWORDS
		if val, ok := (*l.KeywordsTree)[l.GetTokenContent(token)]; ok {
			switch val {
			case 0:
				token.Kind = TOKEN_KEYWORD
			case 1:
				token.Kind = TOKEN_TYPE
			case 2:
				token.Kind = TOKEN_LIB
			}
		}

		if token.Kind == TOKEN_LIB {
			for l.Cursor < l.Content_len && (IsSymbolChar(byte(l.getCharAt(l.Cursor))) || l.getCharAt(l.Cursor) == '.') {
				l.ChopChar(1)
				st++
			}

			for _, prev := range CurrentLineTokens {
				if prev.Kind == TOKEN_SYMBOL {
					// fmt.Println(l.GetTokenContent(token))
					ImportedSymbs[l.GetTokenContent(prev)] = 0
					prev.Kind = TOKEN_IMPORTED
				}
			}
		}

		if _, ok := ImportedSymbs[l.GetTokenContent(token)]; ok {
			token.Kind = TOKEN_IMPORTED
		}
		token.Len = st
		return token
	}

	l.ChopChar(1)
	token.Kind = TOKEN_INVALID
	token.Len = 1
	return token
}

func (l *Lexer) GetTokens() *[]*Token {
	tokens := []*Token{}
	oldLine := -1
	for l.Cursor < l.Content_len {
		next := l.NextToken()
		if next.Addr.Origin != oldLine {
			CurrentLineTokens = []*Token{}
			oldLine = next.Addr.Origin
		}
		CurrentLineTokens = append(CurrentLineTokens, next)
		tokens = append(tokens, next)
	}
	return &tokens
}

func (l *Lexer) GetTokenContent(token *Token) string {
	Assert(token.Addr.X+token.Len <= l.Content_len, "TOKEN OUT OF RANGE")
	return l.Content[token.Addr.X : token.Addr.X+token.Len]
}

// to lex multiple files
func (l *Lexer) ResetContent(content string) {
	l.Content = content
	l.Content_len = len(content)
	l.Cursor = 0
	l.Line = 0
	l.LineStart = 0
}
