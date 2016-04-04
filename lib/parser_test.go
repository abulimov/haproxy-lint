package lib

import "testing"

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
