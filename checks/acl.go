package checks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abulimov/haproxy-lint/lib"
)

// CheckUnusedACL checks if we have ACLs that are defined but not used
func CheckUnusedACL(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	ACLs := make(map[lib.Entity]bool)
	var otherLines []int
	for i, line := range s.Content {
		acl := getACL(line)
		if acl != "" {
			ACLs[lib.Entity{Line: i, Name: acl}] = false
		} else {
			otherLines = append(otherLines, i)
		}
	}
	for _, i := range otherLines {
		for acl := range ACLs {
			if usesACL(acl.Name, s.Content[i]) {
				ACLs[acl] = true
			}
		}
	}
	for acl, used := range ACLs {
		if !used {
			problems = append(
				problems,
				lib.Problem{
					Line:     acl.Line,
					Col:      0,
					Severity: "warning",
					Message:  fmt.Sprintf("ACL %s declared but not used", acl.Name),
				},
			)
		}
	}
	return problems
}

func getACL(l string) string {
	return lib.GetUsage("acl", l)
}

func usesACL(acl, line string) bool {
	return strings.Contains(lib.StripComments(line), acl)
}

func getUsedACLs(line string) []string {
	preDefinedACLs := map[string]bool{
		"FALSE":          true,
		"HTTP":           true,
		"HTTP_1.0":       true,
		"HTTP_1.1":       true,
		"HTTP_CONTENT":   true,
		"HTTP_URL_ABS":   true,
		"HTTP_URL_SLASH": true,
		"HTTP_URL_STAR":  true,
		"LOCALHOST":      true,
		"METH_CONNECT":   true,
		"METH_GET":       true,
		"METH_HEAD":      true,
		"METH_OPTIONS":   true,
		"METH_POST":      true,
		"METH_TRACE":     true,
		"RDP_COOKIE":     true,
		"REQ_CONTENT":    true,
		"TRUE":           true,
		"WAIT_END":       true,
	}
	var acls []string
	var rawACLs []string
	afterIfRE := regexp.MustCompile(`\w+\s+if\s+(.+)`)
	afterIfMatch := afterIfRE.FindAllStringSubmatch(lib.StripComments(line), -1)
	if len(afterIfMatch) > 0 {
		if len(afterIfMatch[0]) > 1 {
			afterIfString := afterIfMatch[0][1]
			word := regexp.MustCompile(`({[^}]+}|\w+)`)
			rawACLs = word.FindAllString(afterIfString, -1)
		}
	}
	for _, acl := range rawACLs {
		_, preDefined := preDefinedACLs[acl]
		if acl != "or" && !preDefined {
			acls = append(acls, acl)
		}
	}
	return acls
}

func isInlineACL(acl string) bool {
	re := regexp.MustCompile(`!?{[^}]+}`)
	return re.MatchString(acl)
}

// CheckUnknownACLs checks if we have ACLs that are used in config but not defined
func CheckUnknownACLs(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	ACLs := make(map[string]bool)
	var otherLines []int
	for i, line := range s.Content {
		acl := getACL(line)
		if acl != "" {
			ACLs[acl] = false
		} else {
			otherLines = append(otherLines, i)
		}
	}
	for _, i := range otherLines {
		usedACLs := getUsedACLs(s.Content[i])
		for _, acl := range usedACLs {
			if !isInlineACL(acl) {
				if _, found := ACLs[acl]; !found {
					problems = append(
						problems,
						lib.Problem{
							Line:     i,
							Col:      0,
							Severity: "critical",
							Message:  fmt.Sprintf("ACL %s used but not declared", acl),
						},
					)
				}
			}
		}

	}
	return problems
}
