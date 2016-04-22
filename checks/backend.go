package checks

import (
	"fmt"

	"github.com/abulimov/haproxy-lint/lib"
	B "github.com/abulimov/haproxy-lint/lib/backend"
)

// CheckUnusedBackends checks if we have declared but not used backends
func CheckUnusedBackends(sections []*lib.Section) []lib.Problem {
	var problems []lib.Problem
	backends := make(map[lib.Entity]bool)
	var other []*lib.Section
	for _, s := range sections {
		if s.Type == "backend" {
			backend := lib.Entity{Line: s.Line, Name: s.Name}
			backends[backend] = false
		} else {
			other = append(other, s)
		}
	}
	for b := range backends {
		for _, s := range other {
			for _, line := range s.Content {
				if B.LineUsesBackend(b.Name, line) {
					backends[b] = true
				}
			}
		}
	}
	for b, used := range backends {
		if !used {
			problems = append(
				problems,
				lib.Problem{
					Line:     b.Line,
					Col:      0,
					Severity: "warning",
					Message:  fmt.Sprintf("backend %s declared but not used", b.Name),
				},
			)
		}
	}
	return problems
}

func CheckUnknownBackends(sections []*lib.Section) []lib.Problem {
	var problems []lib.Problem
	backends := make(map[string]bool)
	var other []*lib.Section
	for _, s := range sections {
		if s.Type == "backend" {
			backends[s.Name] = false
		} else {
			other = append(other, s)
		}
	}

	for _, s := range other {
		for i, line := range s.Content {
			backend := B.GetNameFromDeclaration(line)
			if backend != "" {
				if _, found := backends[backend]; !found {
					problems = append(
						problems,
						lib.Problem{
							Line:     i,
							Col:      0,
							Severity: "critical",
							Message:  fmt.Sprintf("backend %s used but not declared", backend),
						},
					)
				}
			}
		}
	}

	return problems
}
