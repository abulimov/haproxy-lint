package checks

import (
	"fmt"

	"github.com/abulimov/haproxy-lint/lib"
	A "github.com/abulimov/haproxy-lint/lib/acl"
)

// CheckUnusedACL checks if we have ACLs that are defined but not used
func CheckUnusedACL(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	ACLs := make(map[lib.Entity]bool)
	var otherLines []int
	for i, line := range s.Content {
		acl := A.GetNameFromDeclaration(line)
		if acl != "" {
			ACLs[lib.Entity{Line: i, Name: acl}] = false
		} else {
			otherLines = append(otherLines, i)
		}
	}
	for _, i := range otherLines {
		for acl := range ACLs {
			if A.LineUsesACL(acl.Name, s.Content[i]) {
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

// CheckUnknownACLs checks if we have ACLs that are used in config but not defined
func CheckUnknownACLs(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	ACLs := make(map[string]bool)
	var otherLines []int
	for i, line := range s.Content {
		acl := A.GetNameFromDeclaration(line)
		if acl != "" {
			ACLs[acl] = false
		} else {
			otherLines = append(otherLines, i)
		}
	}
	for _, i := range otherLines {
		usedACLs := A.GetUsedACLNames(s.Content[i])
		for _, acl := range usedACLs {
			if !A.IsInline(acl) {
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
