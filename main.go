package main

import (
	. "LanguageFuck/Lexer"
	"fmt"
	"os"
)

// TODOs :
//   X Make a lexer <dynamic -- Each language has  a toml file>
//   - Make sure that the thing is scalable to support multiple files in a project

var tree map[string]uint8

func init() {
	tree = make(map[string]uint8, len(KEYWORDS))
	for _, keyword := range KEYWORDS {
		tree[keyword] = 0
	}
}

func main() {
	// str := "package main \n // hello ducj\nimport \"fmt\" \n\nfunc main() {\n  fmt.Println(1/2) \n  return 0 \n}"
	b, _ := os.ReadFile("main.go")
	str := string(b)
	l := LexerInit(str, &tree)
	tokens := l.GetTokens()
	// old := -1
	// for _, token := range tokens {
	// 	// fmt.Println(
	// 	// 	l.Content[token.Addr.X:token.Addr.X+token.Len],
	// 	// 	"\t<<", token.Addr.Line, ">>",
	// 	// )
	// 	fmt.Print(
	// 		l.Content[token.Addr.X : token.Addr.X+token.Len],
	// 	)
	// 	if token.Addr.Line != old {
	// 		fmt.Print("\n")
	// 		old = token.Addr.Line
	// 	}
	// }
	old := -1
	cursor := 0
	for _, token := range tokens {
		for cursor != token.Addr.X {
			cursor++
			fmt.Print(" ")
		}
		if token.Addr.Line != old {
			for old != token.Addr.Line {
				fmt.Print("\n")
				old++
			}
			// old = token.Addr.Line
		}
		fmt.Print(l.Content[token.Addr.X : token.Addr.X+token.Len])
		cursor += token.Len
	}
}