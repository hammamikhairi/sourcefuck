package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Lexer"
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Parser"
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Types"
	. "github.com/hammamikhairi/sourcefuck/cmd/sourcefuck/Utils"
)

var (
	tree map[string]uint8
	swap map[string]string

	ext       string
	dec, peek bool
	OutDir    string
	userKey   string
	useOsSig  bool

	files            []string
	base, OutDirPath string
)

func init() {
	var path string
	flag.StringVar(&path, "path", "", "The path to a file or directory.")
	flag.StringVar(&ext, "ext", "", "The file extension for targeted files within a directory. (Required for directories)")
	flag.BoolVar(&dec, "dec", false, "Specify this flag to decrypt an encrypted code using the provided key.")
	flag.BoolVar(&peek, "peek", false, "Specify this flag to only print the obfuscated/decrypted code to the console.")
	flag.StringVar(&OutDir, "out", "", "Specify the output directory for the obfuscated/decrypted code.")
	flag.BoolVar(&useOsSig, "osSig", false, "Use the OS signature as Enc/Dec Key")
	flag.StringVar(&userKey, "key", "", "The encryption/decryption key. (A random key will be generated if unspecified)")
	flag.Parse()

	files, base, OutDirPath = ProcessFlags(path, ext, OutDir, dec)

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
		currentFilePath := filepath.Join(base, file)

		b, err := os.ReadFile(currentFilePath)
		if err != nil {
			panic(err)
		}

		l.SetContent(string(b))
		tokens := l.GetTokens()

		key := Get32BytesKey(userKey, useOsSig)
		pr := ParserInit(tokens, key)
		pr.Parse(l, dec)
		ParsedCode := l.FormatCode(tokens, &pr.Swap)

		if peek {
			fmt.Println("\n<<<<<", currentFilePath, ">>>>>\n ")
			fmt.Print(l.FormatCode(tokens, &pr.Swap))
			continue
		}

		fmt.Printf("[INFO] Writing Code To %s", filepath.Join(OutDirPath, file))
		WriteStringToFile(filepath.Join(OutDirPath, file), ParsedCode)
	}
}
