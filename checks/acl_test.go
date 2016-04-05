package checks

import (
	"reflect"
	"testing"

	"github.com/abulimov/haproxy-lint/lib"
)

func TestCheckUnusedACL(t *testing.T) {
	lines, err := lib.ReadConfigFile("../testdata/haproxy.cfg")
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
	lines, err := lib.ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := lib.GetSection("frontend", lines)
	problems := CheckUnknownACLs(sections[1])

	if len(problems) != 1 {
		t.Errorf("Expected %d problems, got %d", 1, len(problems))
	}
}

func TestUsesACL(t *testing.T) {
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
		if got := usesACL(tt.acl, tt.line); got != tt.want {
			t.Errorf("%q. usesACL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetUsedACLs(t *testing.T) {
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
		if got := getUsedACLs(tt.line); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getUsedACLs() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsInlineACL(t *testing.T) {
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
		if got := isInlineACL(tt.acl); got != tt.want {
			t.Errorf("%q. isInlineACL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
