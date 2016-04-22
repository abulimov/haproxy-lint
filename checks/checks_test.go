package checks

import (
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestRun(t *testing.T) {
	lines, err := lib.GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSections(lines)
	problems := Run(sections, false)

	if len(problems) != 9 {
		t.Errorf("Expected %d problems, got %d", 9, len(problems))
	}
}

func BenchmarkRun(b *testing.B) {
	lines, err := lib.GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		b.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSections(lines)
	// run the Run function b.N times
	for n := 0; n < b.N; n++ {
		Run(sections, false)
	}
}
