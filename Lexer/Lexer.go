package lexer

type Lexer struct {
	Content      string
	Content_len  int
	Cursor       int
	Line         int
	KeywordsTree *map[string]uint8
}

func LexerInit(content string, tree *map[string]uint8) *Lexer {
	return &Lexer{content, len(content), 0, 0, tree}
}

func (l *Lexer) ChopChar(len int) {
	for i := 0; i < len; i++ {
		current := string(l.Content[l.Cursor])
		l.Cursor++
		if current == "\n" {
			l.Line++
		}
	}
}

func (l *Lexer) Trim() {
	for l.Cursor < l.Content_len && isSpace(string(l.Content[l.Cursor])) || l.getCharAt(l.Cursor) == "\n" {
		l.ChopChar(1)
	}
}

func (l *Lexer) getCharAt(pos int) string {
	return string(l.Content[l.Cursor])
}

func (l *Lexer) startsWith(prefix string) bool {

	if len(prefix) == 1 {
		return l.getCharAt(l.Cursor) == prefix
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

	token := &Token{}
	token.Addr = Vec2i{l.Cursor, l.Line}

	st := 0
	if l.startsWith("\"") {
		token.Kind = TOKEN_STRING
		l.ChopChar(1)
		for l.getCharAt(l.Cursor) != "\"" {
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
		for l.getCharAt(l.Cursor) == "\t" {
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
		for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != "\n" {
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

	if isAlpha(l.getCharAt(l.Cursor)) {
		token.Kind = TOKEN_SYMBOL
		for l.Cursor < l.Content_len && isAlpha(l.getCharAt(l.Cursor)) {
			l.ChopChar(1)
			st++
		}

		lastToken := l.Content[token.Addr.X : token.Addr.X+st]
		// PREPROC
		if lastToken == "package" {
			token.Kind = TOKEN_PREPROC
			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != "\n" {
				st++
				l.ChopChar(1)
			}
			token.Len = st
			return token
		}

		// PREPROC
		if lastToken == "import" {
			token.Kind = TOKEN_PREPROC
			l.Trim()
			end := "\n"
			if l.startsWith("(") {
				end = ")"
			}

			c := 0
			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != end {
				st++
				l.ChopChar(1)
				if l.startsWith("\n") {
					c++
				}
			}

			if end == ")" {
				st += 1
				l.ChopChar(1)
				l.Line -= c
			}

			st++
			l.ChopChar(1)

			token.Len = st
			return token
		}

		// KEYWORDS
		if _, ok := (*l.KeywordsTree)[lastToken]; ok {
			token.Kind = TOKEN_KEYWORD
		}

		token.Len = st
		return token
	}

	l.ChopChar(1)
	token.Kind = TOKEN_INVALID
	token.Len = 1
	return token
}

func (l *Lexer) GetTokens() []*Token {
	tokens := []*Token{}
	for l.Cursor < l.Content_len {
		tokens = append(tokens, l.NextToken())
	}
	return tokens
}
