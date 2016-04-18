package lib

import (
	"reflect"
	"testing"
)

func TestGetUsage(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		keyword string
		line    string
		// Expected results.
		want string
	}{
		{
			name:    "Real usage",
			keyword: "acl",
			line:    "acl h_some hdr(Host) -i some.example.com",
			want:    "h_some",
		},
		{
			name:    "Other keyword usage",
			keyword: "backend",
			line:    "acl h_something hdr(Host) -i some.example.com",
			want:    "",
		},
		{
			name:    "Double keyword usage",
			keyword: "http-response set-header",
			line:    "http-response set-header Cache-Control max-age=300,must-revalidate",
			want:    "Cache-Control",
		},
		{
			name:    "no option usage",
			keyword: "option",
			line:    "no option abortonclose",
			want:    "abortonclose",
		},
		{
			name:    "Bad acl line",
			keyword: "acl",
			line:    "acl",
			want:    "",
		},
	}
	for _, tt := range tests {
		if got := GetUsage(tt.keyword, tt.line); got != tt.want {
			t.Errorf("%q. GetUsage() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		lines   []string
		pattern string
		// Expected results.
		want []string
	}{
		{
			name: "File without lines to filter",
			lines: []string{
				"global",
				"    daemon",
				"    maxconn 256",
			},
			pattern: "{%",
			want: []string{
				"global",
				"    daemon",
				"    maxconn 256",
			},
		},
		{
			name: "File with lines to filter",
			lines: []string{
				"global",
				"    daemon",
				"    {% if haproxy_domain == 'ru' %}",
				"    maxconn 1024",
				"    {% else %}",
				"    maxconn 256",
				"    {% endif %}",
			},
			pattern: "{%",
			want: []string{
				"global",
				"    daemon",
				"",
				"    maxconn 1024",
				"",
				"    maxconn 256",
				"",
			},
		},
	}
	for _, tt := range tests {
		if got := Filter(tt.lines, tt.pattern); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
