package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/crypto/hkdf"
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

func IsSymbolChar(char byte) bool {
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

	if path == "" {
		log.Fatal("Please specify a path using the -path flag.")
	}

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
		path = filepath.Base(path)
		files = append(files, path)

		base = filepath.Dir(base)
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

func ProcessFlags(path, ext, OutDir string, dec bool) ([]string, string, string) {

	files, base := GetFiles(path, ext)
	var OutDirPath string = base

	if OutDir == "" {
		if dec {
			OutDirPath = filepath.Join(OutDirPath, "/Dec")
		} else {
			OutDirPath = filepath.Join(OutDirPath, "/Enc")
		}
	} else {
		OutDirPath, _ = filepath.Abs(OutDir)
	}

	return files, base, OutDirPath
}

func Get32BytesKey(key string, useOsSig bool) []byte {
	var base []byte
	if useOsSig {
		base = getOsSignature()
	} else {
		if key != "" {
			base, _ = hex.DecodeString(key)
			Assert(len(base) == 32, "Key must be 32 bytes.")
			return base
		} else {
			base = generateRandomKey()
		}
	}

	hkdf := hkdf.New(sha256.New, base, nil, nil)
	keyBytes := make([]byte, 32)
	n, err := hkdf.Read(keyBytes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return nil
	}

	if useOsSig {
		fmt.Printf("[INFO] Using OS Signature as KEY [%x].\n", base)
	} else {
		fmt.Printf("[INFO] Using Generated KEY [%x].\n", base)
	}
	return keyBytes[:n]
}

func getOsSignature() []byte {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname : %s", err)
		os.Exit(1)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error getting network interfaces : %s", err)
		os.Exit(1)
	}

	var macAddr string
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 {
			macAddr = iface.HardwareAddr.String()
			break
		}
	}

	return []byte(hostname + macAddr)
}

func generateRandomKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Error generating random AES Key")
		os.Exit(1)
	}
	return key
}
