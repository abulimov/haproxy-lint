package checks

import (
	"fmt"
	"sort"

	"github.com/abulimov/haproxy-lint/lib"
	A "github.com/abulimov/haproxy-lint/lib/acl"
)

// CheckUnreachableRules checks if we have rules that will never be
// triggered because their acls are less generic than similar acls
// declared before.
func CheckUnreachableRules(s *lib.Section) []lib.Problem {
	checkableKWs := map[string]bool{
		"use_backend": true,
		"redirect":    true,
		"use-server":  true,
	}

	var problems []lib.Problem
	var lineNums []int
	aclsMap := make(map[int][]*A.ACL)
	kwMap := make(map[int]string)
	for i, line := range s.Content {
		acls := A.ParseACLs(line)
		if len(acls) > 0 {
			lineNums = append(lineNums, i)
			aclsMap[i] = acls
			kwMap[i] = lib.GetKeyword(line)
		}
	}
	sort.Ints(lineNums)
	usedACLs := make(map[A.ACL]int)
	// for each line from start of section
	for _, i := range lineNums {
		curACLs := aclsMap[i]
		curKW := kwMap[i]
		// for each acl in current acls
		for _, acl := range curACLs {
			// if we already seen some of current acls
			if prevLine, found := usedACLs[*acl]; found {
				// if current keyword is same as in line when this ACL was used
				if _, ok := checkableKWs[curKW]; ok && curKW == kwMap[prevLine] {
					// if previous acl line is more generic than current
					if A.In(aclsMap[prevLine], curACLs) && !A.HasOrs(curACLs) {
						problems = append(
							problems,
							lib.Problem{
								Line:     prevLine,
								Col:      0,
								Severity: "warning",
								Message:  fmt.Sprintf("This line shadows rule '%s' on line %d", s.Content[i], i),
							},
						)
						problems = append(
							problems,
							lib.Problem{
								Line:     i,
								Col:      0,
								Severity: "warning",
								Message:  fmt.Sprintf("Unreachable rule - '%s' on line %d will match first", s.Content[prevLine], prevLine),
							},
						)
					}
				}
			} else {
				// add current acl to used acls map
				usedACLs[*acl] = i
			}
		}
	}
	return problems
}
