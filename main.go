package main

import (
	. "LanguageFuck/Lexer"
	"log"
	"path/filepath"
	"strings"

	. "LanguageFuck/Parser"
	. "LanguageFuck/Types"
	"flag"
	"fmt"
	"os"
)

// TODOs :
//   X Make a lexer <dynamic -- Each language has  a toml file>
//   - Make sure that the thing is scalable to support multiple files in a project

var tree map[string]uint8
var swap map[string]string

var SypportedLangs []string = []string{
	"go",
}

func getExtension(fileName string) string {
	if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 {
		return fileName[dotIndex+1:]
	}
	return ""
}

var files []string
var ext string
var dec bool

func init() {
	// flags
	var path string
	flag.StringVar(&path, "path", "", "the path to a file or directory")
	flag.StringVar(&ext, "ext", "", "file extension")
	flag.BoolVar(&dec, "decrypt", false, "flag to decrypt")
	flag.Parse()

	if path == "" {
		log.Fatal("Please specify a path using the -path flag.")
	}

	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Error getting information about the path: %s\n %v\n", path, err)
	}

	if info.IsDir() {

		if ext == "" {
			log.Fatal("Please specify the extension using the -ext flag when dealing with directories.")
		}

		err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && getExtension(p) == ext {
				path, _ := filepath.Abs(p)
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		path = filepath.Clean(path)
		path, _ = filepath.Abs(path)
		if ext == "" {
			ext = getExtension(path)
		}
		files = append(files, path)
	}
	fmt.Println(files, ext)

	tree = make(map[string]uint8, len(KEYWORDS))
	for _, keyword := range KEYWORDS {
		tree[keyword] = 0
	}

	for _, lang_type := range TYPES {
		tree[lang_type] = 1
	}

	for _, lib := range LIBRARIES {
		tree[lib] = 2
	}
}

func main() {

	l := LexerInit("", &tree)
	for _, file := range files {
		fmt.Println("\n<<<<<", file, ">>>>>\n ")
		b, _ := os.ReadFile(file)
		l.ResetContent(string(b))
		tokens := l.GetTokens()

		// TODO : parse and change
		pr := ParserInit(tokens, 9)
		pr.Parse(l, dec)

		// PrintCode(l, tokens)
		PrintNewCode(l, tokens, &pr.Swap)
	}
	// b, _ := os.ReadFile("main.go")
	// l := LexerInit(string(b), &tree)
	// tokens := l.GetTokens()

}

func PrintNewCode(l *Lexer, tokens *[]*Token, swap *map[string]string) {
	old := 0
	cursor := 0
	for _, token := range *tokens {

		if token.Kind == TOKEN_END {
			continue
		}

		for cursor != token.Addr.X {
			cursor++
			fmt.Print(" ")
		}

		if token.Addr.Line != old {
			for old != token.Addr.Line {
				fmt.Print("\n")
				old++
			}
		}
		if token.Kind == TOKEN_SYMBOL {
			fmt.Print((*swap)[l.GetTokenContent(token)])
		} else {
			fmt.Print(l.GetTokenContent(token))
		}
		cursor += token.Len
	}
}

func PrintTokens(l *Lexer, tokens *[]*Token) {
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

func PrintCode(l *Lexer, tokens *[]*Token) {
	old := 0
	cursor := 0
	for _, token := range *tokens {

		if token.Kind == TOKEN_END {
			continue
		}

		for cursor != token.Addr.X {
			cursor++
			fmt.Print(" ")
		}

		if token.Addr.Line != old {
			for old != token.Addr.Line {
				fmt.Print("\n")
				old++
			}
		}
		fmt.Print(l.GetTokenContent(token))
		cursor += token.Len
	}
}