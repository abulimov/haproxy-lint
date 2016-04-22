package acl

import (
	"reflect"
	"testing"
)

func TestLineUsesACL(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		acl  string
		line string
		// Expected results.
		want bool
	}{
		{
			name: "line with h_test acl",
			acl:  "h_test",
			line: "redirect scheme https code 301 if !h_test or h_some",
			want: true,
		},
		{
			name: "line without h_test acl",
			acl:  "h_test",
			line: "redirect scheme https code 301 if !h_missing or h_some",
			want: false,
		},
	}
	for _, tt := range tests {
		if got := LineUsesACL(tt.acl, tt.line); got != tt.want {
			t.Errorf("%q. usesACL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsInline(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		acl string
		// Expected results.
		want bool
	}{
		{
			name: "normal inline ACL",
			acl:  "{ ssl_fc }",
			want: true,
		},
		{
			name: "bad inline ACL",
			acl:  "ssl_fc }",
			want: false,
		},
	}
	for _, tt := range tests {
		if got := IsInline(tt.acl); got != tt.want {
			t.Errorf("%q. isInlineACL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestParseACLs(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		line string
		// Expected results.
		want []*ACL
	}{
		{
			name: "acls with negation",
			line: "redirect scheme https code 301 if !{ ssl_fc } h_missing",
			want: []*ACL{
				{
					Name:       "{ ssl_fc }",
					Negated:    true,
					Inline:     true,
					Predefined: false,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
		},
		{
			name: "pre-defined acl and 'or'",
			line: "redirect scheme https code 301 if METH_GET or h_missing",
			want: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
		},
		{
			name: "short usage",
			line: "block if some_acl",
			want: []*ACL{
				{
					Name:       "some_acl",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
		},
		{
			name: "Jinja-template if",
			line: "{% if something == 'other' %}",
			want: []*ACL(nil),
		},
	}
	for _, tt := range tests {
		if got := ParseACLs(tt.line); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ParseACLs() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetUsedACLNames(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		line string
		// Expected results.
		want []string
	}{
		{
			name: "acls with negation",
			line: "redirect scheme https code 301 if !{ ssl_fc } h_missing",
			want: []string{"{ ssl_fc }", "h_missing"},
		},
		{
			name: "pre-defined acl and 'or'",
			line: "redirect scheme https code 301 if METH_GET or h_missing",
			want: []string{"h_missing"},
		},
		{
			name: "short usage",
			line: "block if some_acl",
			want: []string{"some_acl"},
		},
		{
			name: "Jinja-template if",
			line: "{% if something == 'other' %}",
			want: []string(nil),
		},
	}
	for _, tt := range tests {
		if got := GetUsedACLNames(tt.line); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getUsedACLs() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		first  []*ACL
		second []*ACL
		// Expected results.
		want bool
	}{
		{
			name: "Equal must be In",
			first: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			second: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			want: true,
		},
		{
			name: "Bigger cannot be In",
			first: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
				{
					Name:       "h_ok",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			second: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			want: false,
		},
		{
			name: "Smaller must be In",
			first: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_ok",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			second: []*ACL{
				{
					Name:       "METH_GET",
					Negated:    false,
					Inline:     false,
					Predefined: true,
				},
				{
					Name:       "h_missing",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
				{
					Name:       "h_ok",
					Negated:    false,
					Inline:     false,
					Predefined: false,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := In(tt.first, tt.second); got != tt.want {
			t.Errorf("%q. In() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
