package backend

import "testing"

func TestLineUsesBackend(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		backend string
		line    string
		// Expected results.
		want bool
	}{
		{
			name:    "backend in 'use_backend'",
			backend: "other-servers",
			line:    "use_backend undefined-servers if h_test",
			want:    false,
		},
		{
			name:    "backend in 'use_backend'",
			backend: "undefined-servers",
			line:    "use_backend undefined-servers if h_test",
			want:    true,
		},
		{
			name:    "backend in 'default_backend'",
			backend: "other-servers",
			line:    "default_backend other-servers",
			want:    true,
		},
		{
			name:    "backend in 'default_backend'",
			backend: "undefined-servers",
			line:    "default_backend other-servers",
			want:    false,
		},
	}
	for _, tt := range tests {
		if got := LineUsesBackend(tt.backend, tt.line); got != tt.want {
			t.Errorf("%q. usesBackend() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetNameFromDeclaration(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		line string
		// Expected results.
		want string
	}{
		{
			name: "backend in 'use_backend'",
			line: "use_backend undefined-servers if h_test",
			want: "undefined-servers",
		},
		{
			name: "backend in 'default_backend'",
			line: "default_backend other-servers",
			want: "other-servers",
		},
	}
	for _, tt := range tests {
		if got := GetNameFromDeclaration(tt.line); got != tt.want {
			t.Errorf("%q. getBackend() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
