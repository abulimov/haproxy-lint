package lib

import (
	"regexp"
	"strings"
)

// Sections is a list of all possible section names in HAProxy config
var Sections = [...]string{
	"global",
	"defaults",
	"listen",
	"frontend",
	"backend",
}

// Section is a section in HAProxy config file (for example, 'frontend')
type Section struct {
	Type    string
	Name    string
	Line    int
	Content map[int]string
}

// StripComments returns string without comments and extra spaces
func StripComments(s string) string {
	// strip comments
	re := regexp.MustCompile("#.*")
	replaced := re.ReplaceAllString(s, "")
	return strings.TrimSpace(replaced)
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

// GetSection constructs Section from slice of lines
func GetSection(section string, lines []string) []*Section {
	var result []*Section
	start := false
	var current *Section

	for i, line := range lines {
		line := strings.TrimSpace(line)
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
				current.Line = i + 1
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

// GetSections returns all Sections from slice of config file lines
func GetSections(lines []string) []*Section {
	var result []*Section

	for _, s := range Sections {
		ps := GetSection(s, lines)
		result = append(result, ps...)
	}
	return result
}
