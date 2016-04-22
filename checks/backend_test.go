package checks

import (
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckUnusedBackends(t *testing.T) {
	lines, err := lib.ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSections(lines)
	problems := CheckUnusedBackends(sections)

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}

func TestCheckUnknownBackends(t *testing.T) {
	lines, err := lib.ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSections(lines)
	problems := CheckUnknownBackends(sections)

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}
