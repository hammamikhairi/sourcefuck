package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// checks if the entire string is a single whitespace
func IsSpace(char string) bool {
	return char == " "
}

// checks if the 1st byte-character in the string is an uppercase alphabet letter
func IsUpper(char string) bool {
	return char[0] >= 'A' && char[0] <= 'Z'
}

// checks if the 1st byte-character in the string is a lowercase alphabet letter
func IsLower(char string) bool {
	return char[0] >= 'a' && char[0] <= 'z'
}

// checks if the 1st byte-character in the string is an alphabet letter
func IsAlpha(char string) bool {
	return IsLower(char) || IsUpper(char)
}

func IsSymbolChar(char string) bool {
	return unicode.IsNumber(rune(char[0])) || IsAlpha(char) || char == "_"
}

func Assert(cond bool, errorM string) {
	if cond == false {
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
	if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -2 {
		return fileName[dotIndex+1:]
	}
	return ""
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
