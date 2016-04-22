package lib

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ParseHaproxyLine parses line of HAProxy config check output
func ParseHaproxyLine(line string) *Problem {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`\[(?P<severity>\w+)\]\s+\d+\/\d+\s+\(\d+\)\s+:\s+parsing\s\[.+:(?P<line>\d+)\]\s+:\s+(?P<message>.+)`),
		regexp.MustCompile(`\[(?P<severity>\w+)\]\s+\d+\/\d+\s+\(\d+\)\s+:\s+(?P<message>.+)\s+at\s\[.+:(?P<line>\d+)\]\s+.+`),
	}
	stopWords := regexp.MustCompile(`unable to load SSL private key|file|cannot find user id for|cannot find group id for`)
	if stopWords.MatchString(line) {
		return nil
	}
	for _, re := range regexps {
		matches := re.FindAllStringSubmatch(line, -1)
		pos := map[string]int{}
		for i, name := range re.SubexpNames() {
			pos[name] = i
		}
		if len(matches) == 1 {
			if len(matches[0]) == 4 {
				lineNum, err := strconv.Atoi(matches[0][pos["line"]])
				if err != nil {
					return nil
				}
				severity := "critical"
				if matches[0][pos["severity"]] == "WARNING" {
					severity = "warning"
				}
				return &Problem{
					Line:     lineNum,
					Col:      0,
					Severity: severity,
					Message:  matches[0][pos["message"]],
				}
			}
		}
	}
	return nil
}

// ParseHaproxyOutput parses whole HAProxy config check output
func ParseHaproxyOutput(f io.Reader) ([]Problem, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var problems []Problem
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		p := ParseHaproxyLine(line)
		if p != nil {
			problems = append(problems, *p)
		}
	}
	return problems, nil
}

// RunHAProxyCheck executes HAProxy binary in check mode and parses it's output
func RunHAProxyCheck(filePath string) ([]Problem, error) {
	cmd := exec.Command("haproxy", "-c", "-f", filePath)
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		// haproxy exits with exit code 1 on any problems with config,
		// but if we don't have haproxy executable we won't get exit code.
		// Here we check it.
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, err
		}

	}
	return ParseHaproxyOutput(&out)
}

// CreateTempConfig creates temporary file, containing given lines.
// If no error returned, this file needs to be removed by the caller.
func CreateTempConfig(lines []string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "haproxy-lint")
	if err != nil {
		return "", err
	}
	for _, line := range lines {
		_, err := tmpfile.Write([]byte(line + "\n"))
		if err != nil {
			os.Remove(tmpfile.Name())
			return "", err
		}
	}
	return tmpfile.Name(), nil
}

// CheckWithHAProxy runs haproxy checks with temporary file creation
func CheckWithHAProxy(config []string, filePath string, createTempFile bool) ([]Problem, error) {
	if createTempFile {
		// we need to create temp file with content
		tempFilePath, err := CreateTempConfig(config)
		if err != nil {
			log.Printf("Failed to create temp file to filter config: %v\n", err)
		} else {
			defer os.Remove(tempFilePath) // clean up
			filePath = tempFilePath
		}
	}
	return RunHAProxyCheck(filePath)
}
