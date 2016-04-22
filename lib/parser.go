package lib

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
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

func trimNo(line string) string {
	trimmed := strings.TrimSpace(line)
	return strings.TrimSpace(strings.TrimPrefix(trimmed, "no "))
}

// StripComments returns string without comments and extra spaces
func StripComments(s string) string {
	// strip comments
	re := regexp.MustCompile(`\s*#.*`)
	return re.ReplaceAllString(s, "")
}

// GetKeyword returns used keyword
func GetKeyword(line string) string {
	trimmed := trimNo(line)
	splitted := strings.Fields(trimmed)
	if len(splitted) > 0 {
		return splitted[0]
	}
	return ""
}

// GetUsage returns name of used rule, or "" if line does not contain rule
func GetUsage(keyword, line string) string {
	trimmed := trimNo(line)
	if strings.HasPrefix(trimmed, keyword) {
		rest := strings.TrimPrefix(trimmed, keyword)
		splitted := strings.Fields(rest)
		if len(splitted) > 0 {
			return splitted[0]
		}
	}
	return ""
}

// CleanupConfig helps us replace lines matching pattern with empty strings.
// Helpfull when we are trying to strip some template engine conditionals.
// This funtion also removes all comments.
func CleanupConfig(lines []string, pattern string) []string {
	var result []string

	re := regexp.MustCompile(pattern)
	for _, line := range lines {
		if pattern != "" && re.MatchString(line) {
			result = append(result, "")
		} else {
			result = append(result, StripComments(line))
		}
	}
	return result
}

// GetConfig is a small wrapper for reading and cleaning up config
func GetConfig(filePath, ignorePattern string) ([]string, error) {
	config, err := ReadConfigFile(filePath)
	if err != nil {
		return nil, err
	}
	return CleanupConfig(config, ignorePattern), nil
}
