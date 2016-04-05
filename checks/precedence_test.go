package checks

import (
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckPrecedence(t *testing.T) {
	lines, err := lib.ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSection("frontend", lines)
	problems := CheckPrecedence(sections[1])

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}
