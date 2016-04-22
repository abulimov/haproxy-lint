package checks

import (
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckUnusedACL(t *testing.T) {
	lines, err := lib.GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSection("frontend", lines)
	problems := CheckUnusedACL(sections[1])

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}

func TestCheckUnknownACLs(t *testing.T) {
	lines, err := lib.GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSection("frontend", lines)
	problems := CheckUnknownACLs(sections[1])

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}
