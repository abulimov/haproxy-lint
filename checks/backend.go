package checks

import (
	"fmt"

	"github.com/abulimov/haproxy-linter/lib"
)

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
				if usesBackend(b.Name, line) {
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

func usesBackend(backend, line string) bool {
	var keywords = []string{"use_backend", "default_backend"}
	name := ""
	for _, kw := range keywords {
		name = lib.GetUsage(kw, line)
		if name == backend {
			return true
		}
	}
	return false
}

func getBackend(line string) string {
	var keywords = []string{"use_backend", "default_backend"}
	name := ""
	for _, kw := range keywords {
		name = lib.GetUsage(kw, line)
		if name != "" {
			return name
		}
	}
	return ""
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
			backend := getBackend(line)
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
