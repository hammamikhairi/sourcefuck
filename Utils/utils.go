package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// checks if the entire string is a single ASCII whitespace
func IsSpace(s string) bool {
	return s == " "
}

// checks if ASCII uppercase
func IsUpper(char byte) bool {
	return char >= 'A' && char <= 'Z'
}

// checks if ASCII lowercase
func IsLower(char byte) bool {
	return char >= 'a' && char <= 'z'
}

// checks if the 1st byte-character is an ASCII alphabet letter
func IsAlpha(char byte) bool {
	return IsLower(char) || IsUpper(char)
}

func IsSymbolChar(s string) bool {
	char := s[0]
	return unicode.IsNumber(rune(char)) || IsAlpha(char) || char == '_'
}

func Assert(cond bool, errorM string) {
	if !cond {
		log.Fatal(errorM)
	}
}

func WriteStringToFile(path, content string) error {
	// Open the file for writing, creating it if it doesn't exist
	CreateFileWithPath(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func getExtension(fileName string) string {
	dotIndex := strings.LastIndex(fileName, ".")
	// special case for no dots (-1), or dotfile (0)
	if dotIndex < 1 {
		return ""
	}
	beforeDot := fileName[dotIndex-1]
	// Unix and Windows dotfiles
	if beforeDot == '/' || beforeDot == '\\' {
		return ""
	}
	return fileName[dotIndex+1:]
}

func GetFiles(path, ext string) ([]string, string) {

	var files []string
	base, _ := filepath.Abs(path)
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
				files = append(files, path[len(base):])
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking the directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		path = filepath.Clean(path)
		if ext == "" {
			ext = getExtension(path)
		}
		files = append(files, path)
		base = base[:len(base)-len(files[0])]
	}

	return files, base
}

func CreateFileWithPath(path string) (*os.File, error) {
	// Ensure that the directory exists by creating any missing parent directories
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}

	// Create the file at the given path
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
