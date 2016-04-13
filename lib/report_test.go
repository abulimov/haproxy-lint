package lib

import "testing"

var testProblems = []Problem{
	{
		Col:      0,
		Line:     24,
		Severity: "warning",
		Message:  "ACL h_some declared but not used",
	},
	{
		Col:      0,
		Line:     25,
		Severity: "critical",
		Message:  "backend undefined-servers used but not declared",
	},
}

func TestReportProblems(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		problems []Problem
		// Expected results.
		want string
	}{
		{
			name:     "Empty problems",
			problems: []Problem{},
			want:     "",
		},
		{
			name:     "Normal problems",
			problems: testProblems,
			want: `24:0:warning: ACL h_some declared but not used
25:0:critical: backend undefined-servers used but not declared
`,
		},
	}
	for _, tt := range tests {
		if got := ReportProblems(tt.problems); got != tt.want {
			t.Errorf("%q. ReportProblems() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestReportProblemsJSON(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		problems []Problem
		// Expected results.
		want string
	}{
		{
			name:     "Empty problems",
			problems: []Problem{},
			want:     "[]",
		},
		{
			name:     "Normal problems",
			problems: testProblems,
			want: `[
  {
    "col": 0,
    "line": 24,
    "severity": "warning",
    "message": "ACL h_some declared but not used"
  },
  {
    "col": 0,
    "line": 25,
    "severity": "critical",
    "message": "backend undefined-servers used but not declared"
  }
]`,
		},
	}
	for _, tt := range tests {
		if got := ReportProblemsJSON(tt.problems); got != tt.want {
			t.Errorf("%q. ReportProblemsJSON() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
