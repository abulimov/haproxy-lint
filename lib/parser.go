package lib

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ReadConfig(f io.Reader) ([]string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func ReadConfigFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadConfig(f)
}

func GetUsage(keyword, line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, keyword) {
		splitted := strings.Fields(trimmed)
		if len(splitted) > 1 {
			return splitted[1]
		}
	}
	return ""
}
