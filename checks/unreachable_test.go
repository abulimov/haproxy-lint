package checks

import (
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckUnreachableRules(t *testing.T) {
	lines, err := lib.GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSection("frontend", lines)
	problems := CheckUnreachableRules(sections[1])

	if len(problems) != 2 {
		t.Errorf("Expected %d problems, got %d", 2, len(problems))
	}
}
