package lib

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ReadConfig returns slice of lines in io.Reader
func ReadConfig(f io.Reader) ([]string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

// ReadConfigFile returns slice of lines in config file
func ReadConfigFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadConfig(f)
}

// GetUsage returns name of used rule, or "" if line does not contain rule
func GetUsage(keyword, line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "no ") {
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "no "))
	}
	if strings.HasPrefix(trimmed, keyword) {
		rest := strings.TrimPrefix(trimmed, keyword)
		splitted := strings.Fields(rest)
		if len(splitted) > 0 {
			return splitted[0]
		}
	}
	return ""
}
