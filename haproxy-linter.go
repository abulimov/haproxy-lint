package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var Sections = [...]string{
	"global",
	"defaults",
	"listen",
	"frontend",
	"backend",
}

type Section struct {
	Type    string
	Name    string
	Content map[int]string
}

func StripComments(s string) string {
	// strip comments
	re := regexp.MustCompile("#.*")
	replaced := re.ReplaceAllString(s, "")
	return replaced
}

func stringInSlice(needle string, slice []string) bool {
	for _, s := range slice {
		if needle == s {
			return true
		}
	}
	return false
}

func isSectionDelimiter(s string) bool {
	for _, prefix := range Sections {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func isSection(s, section string) bool {
	return strings.HasPrefix(s, section)
}

func getSectionName(s string) string {
	splitted := strings.Fields(StripComments(s))
	if len(splitted) > 1 {
		return splitted[1]
	}
	return ""
}

func getSection(section string, lines []string) []*Section {
	var result []*Section
	start := false
	var current *Section

	for i, line := range lines {
		if isSectionDelimiter(line) {
			// if we have already started grabbing content,
			// we need to stop
			if start {
				start = false
				result = append(result, current)
			}
			if isSection(line, section) {
				current = new(Section)
				current.Name = getSectionName(line)
				current.Type = section
				current.Content = make(map[int]string)
				start = true
				continue
			}
		}

		if start {
			current.Content[i+1] = line
		}
	}
	if start {
		result = append(result, current)
	}
	return result
}

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

type Problem struct {
	Col      int    `json:"col"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Check func(*Section) []Problem

func LintSection(s *Section) []Problem {
	rules := []Check{
		CheckUnusedACL,
	}
	var problems []Problem
	for _, r := range rules {
		problems = append(problems, r(s)...)
	}
	return problems
}

func CheckUnusedACL(s *Section) []Problem {
	var problems []Problem
	type acl struct {
		line int
		name string
	}
	ACLs := make(map[acl]bool)
	var otherLines []int
	for i, line := range s.Content {
		a := GetACL(line)
		if a != "" {
			ACLs[acl{line: i, name: a}] = false
		} else {
			otherLines = append(otherLines, i)
		}
	}
	for _, i := range otherLines {
		for a := range ACLs {
			if UsesACL(a.name, s.Content[i]) {
				ACLs[a] = true
			}
		}
	}
	for a, used := range ACLs {
		if !used {
			problems = append(
				problems,
				Problem{
					Line:     a.line,
					Col:      0,
					Severity: "warn",
					Message:  fmt.Sprintf("ACL %s declared but not used", a.name),
				},
			)
		}
	}
	return problems
}

func GetACL(l string) string {
	trimmed := strings.TrimSpace(l)
	if strings.HasPrefix(trimmed, "acl") {
		splitted := strings.Fields(trimmed)
		if len(splitted) > 1 {
			return splitted[1]
		}
	}
	return ""
}

func UsesACL(acl, line string) bool {
	return strings.Contains(line, acl)
}

func ReportProblems(problems []Problem) string {
	buffer := new(bytes.Buffer)
	for _, p := range problems {
		fmt.Fprintf(buffer, "%d:%d:%s: %s\n", p.Line, p.Col, p.Severity, p.Message)
	}
	return buffer.String()
}

func ReportProblemsJSON(problems []Problem) string {
	s, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return ""
	}
	return string(s)
}

func main() {
	argFile := flag.String("file", "haproxy.cfg", "file to check")
	argJSON := flag.Bool("json", false, "Output in json")

	flag.Parse()

	config, err := ReadConfigFile(*argFile)
	if err != nil {
		fmt.Println(err)
	}
	parsed := make(map[string][]*Section)

	for _, s := range Sections {
		parsed[s] = getSection(s, config)
	}

	for _, value := range parsed {
		for _, b := range value {
			//fmt.Println(b)
			problems := LintSection(b)
			if len(problems) != 0 {
				if *argJSON {
					fmt.Println(ReportProblemsJSON(problems))
				} else {
					fmt.Print(ReportProblems(problems))
				}
			}
		}
	}
}
