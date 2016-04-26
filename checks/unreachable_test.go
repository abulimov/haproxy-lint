package checks

import (
	"reflect"
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckUnreachableRules(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		s *lib.Section
		// Expected results.
		want []lib.Problem
	}{
		{
			name: "ACL with 'or' should't match",
			s: &lib.Section{
				Type: "frontend",
				Name: "name",
				Line: 0,
				Content: map[int]string{
					1: "use_backend undefined-servers if h_test",
					2: "use_backend servers if h_test or { ssl_fc }",
					3: "use_backend servers if h_test h_some",
				},
			},
			want: []lib.Problem{
				{
					Col:      0,
					Line:     1,
					Severity: "warning",
					Message:  "This line shadows rule 'use_backend servers if h_test h_some' on line 3",
				},
				{
					Col:      0,
					Line:     3,
					Severity: "warning",
					Message:  "Unreachable rule - 'use_backend undefined-servers if h_test' on line 1 will match first",
				},
			},
		},
		{
			name: "multiple ACLs should match",
			s: &lib.Section{
				Type: "frontend",
				Name: "name",
				Line: 0,
				Content: map[int]string{
					1: "use_backend undefined-servers if h_test",
					2: "use_backend servers if h_test { ssl_fc }",
					3: "use_backend servers if h_test h_some",
				},
			},
			want: []lib.Problem{
				{
					Col:      0,
					Line:     1,
					Severity: "warning",
					Message:  "This line shadows rule 'use_backend servers if h_test { ssl_fc }' on line 2",
				},
				{
					Col:      0,
					Line:     2,
					Severity: "warning",
					Message:  "Unreachable rule - 'use_backend undefined-servers if h_test' on line 1 will match first",
				},
				{
					Col:      0,
					Line:     1,
					Severity: "warning",
					Message:  "This line shadows rule 'use_backend servers if h_test h_some' on line 3",
				},
				{
					Col:      0,
					Line:     3,
					Severity: "warning",
					Message:  "Unreachable rule - 'use_backend undefined-servers if h_test' on line 1 will match first",
				},
			},
		},
		{
			name: "General acl with 'or' should be found",
			s: &lib.Section{
				Type: "frontend",
				Name: "name",
				Line: 0,
				Content: map[int]string{
					1: "use_backend undefined-servers if h_test or { ssl_fc }",
					2: "use_backend servers if h_test h_some",
				},
			},
			want: []lib.Problem{
				{
					Col:      0,
					Line:     1,
					Severity: "warning",
					Message:  "This line shadows rule 'use_backend servers if h_test h_some' on line 2",
				},
				{
					Col:      0,
					Line:     2,
					Severity: "warning",
					Message:  "Unreachable rule - 'use_backend undefined-servers if h_test or { ssl_fc }' on line 1 will match first",
				},
			},
		},
	}
	for _, tt := range tests {
		if got := CheckUnreachableRules(tt.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CheckUnreachableRules() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
