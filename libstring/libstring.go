// Package libstring provides string related library functions.
package libstring

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ExpandTilde is a convenience function that expands ~ to full path.
func ExpandTilde(path string) string {
	newPath, err := homedir.Expand(path)
	if err != nil {
		return path
	}

	return newPath
}

// ExpandTilde is a convenience function that expands both ~ and $ENV.
func ExpandTildeAndEnv(path string) string {
	path = ExpandTilde(path)
	return os.ExpandEnv(path)
}

// GeneratePassword returns password.
// size determines length of initial seed bytes.
func GeneratePassword(size int) (string, error) {
	// Force minimum size to 32
	if size < 32 {
		size = 32
	}

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

// StripChars removes multiple characters in a string.
func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

// PrettyPrintJSON indents JSON nicely.
func PrettyPrintJSON(jsonBytes []byte) ([]byte, error) {
	var out bytes.Buffer

	err := json.Indent(&out, jsonBytes, "", "    ")
	if err != nil {
		return nil, err
	}

	return out.Bytes(), err
}

// FlattenPhone removes non-number characters from phone number.
func FlattenPhone(phone string) string {
	result := strings.Replace(phone, "-", "", -1)
	result = strings.Replace(result, "(", "", -1)
	result = strings.Replace(result, ")", "", -1)
	result = strings.Replace(result, " ", "", -1)

	return result
}

func JSONToText(input string) string {
	result := strings.Replace(input, "{", "", 1)
	result = strings.Replace(input, "}", "", 1)
	result = strings.Replace(input, `"`, "", -1)

	return result
}
