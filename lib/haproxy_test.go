package lib

import (
	"os"
	"reflect"
	"testing"
)

func TestParseHaproxyLine(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		line string
		// Expected results.
		want *Problem
	}{
		{
			name: "Warning output",
			line: "[WARNING] 098/114746 (5447) : parsing [/tmp/tmp.cfg:42] : a 'redirect' rule placed after a 'use_backend' rule will still be processed before.",
			want: &Problem{Line: 42, Col: 0, Severity: "warning", Message: "a 'redirect' rule placed after a 'use_backend' rule will still be processed before."},
		},
		{
			name: "Alert output",
			line: "[ALERT] 098/114746 (5447) : parsing [/tmp/tmp.cfg:36] : unknown keyword 'deffault_backend' in 'frontend' section",
			want: &Problem{Line: 36, Col: 0, Severity: "critical", Message: "unknown keyword 'deffault_backend' in 'frontend' section"},
		},
		{
			name: "Summary output",
			line: "[ALERT] 098/114746 (5447) : Fatal errors found in configuration.",
			want: nil,
		},
		{
			name: "Bad output",
			line: "    [ ALL] accept-proxy",
			want: nil,
		},
		{
			name: "SSL cert output",
			line: "[ALERT] 098/131824 (13707) : parsing [/tmp/tmp.cfg:34] : 'bind :443' : unable to load SSL private key from PEM file '/cert/cert.pem'.",
			want: nil,
		},
		{
			name: "SSL cert problem",
			line: "[ALERT] 098/131824 (13707) : Proxy 'secured': no SSL certificate specified for bind ':443' at [/tmp/tmp.cfg:34] (use 'crt').",
			want: &Problem{Line: 34, Col: 0, Severity: "critical", Message: "Proxy 'secured': no SSL certificate specified for bind ':443'"},
		},
		{
			name: "User id problem",
			line: "[ALERT] 109/150357 (97550) : parsing [/tmp/tmp.cfg:8] : cannot find user id for 'haproxy' (0:Undefined error: 0)",
			want: nil,
		},
	}
	for _, tt := range tests {
		if got := ParseHaproxyLine(tt.line); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ParseHaproxyLine() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestParseHaproxyOutput(t *testing.T) {
	f, err := os.Open("../testdata/example_output.txt")
	if err != nil {
		t.Fatalf("Unexpected error while opening test output file: %v", err)
	}
	defer f.Close()
	problems, err := ParseHaproxyOutput(f)
	if err != nil {
		t.Fatalf("Unexpected error while parsing test output: %v", err)
	}

	if len(problems) != 6 {
		t.Errorf("Expected %d problems, got %d", 6, len(problems))
	}

}
