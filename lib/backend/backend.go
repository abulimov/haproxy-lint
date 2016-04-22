package backend

import "github.com/abulimov/haproxy-lint/lib"

// LineUsesBackend checks if given line uses given backend
func LineUsesBackend(backend, line string) bool {
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

// GetNameFromDeclaration returns backend name from it's declaration
func GetNameFromDeclaration(line string) string {
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
