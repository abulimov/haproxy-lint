package checks

import (
	"fmt"
	"strings"

	"github.com/abulimov/haproxy-lint/lib"
)

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
	return strings.Contains(line, acl)
}
