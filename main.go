package langfuck

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/hammamikhairi/langfuck/Lexer"
	. "github.com/hammamikhairi/langfuck/Parser"
	. "github.com/hammamikhairi/langfuck/Types"
	. "github.com/hammamikhairi/langfuck/Utils"
)

var (
	tree map[string]uint8
	swap map[string]string

	ext       string
	dec, peek bool
	OutDir    string
	key       int

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
	flag.IntVar(&key, "key", 8, "The encryption/decryption key. (Default is 8)")
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

		pr := ParserInit(tokens, key)
		pr.Parse(l, dec)
		ParsedCode := l.FormatCode(tokens, &pr.Swap)

		if peek {
			fmt.Println("\n<<<<<", currentFilePath, ">>>>>\n ")
			fmt.Print(l.FormatCode(tokens, &pr.Swap))
			continue
		}

		WriteStringToFile(filepath.Join(OutDirPath, file), ParsedCode)
	}
}
