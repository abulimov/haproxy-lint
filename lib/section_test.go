package lib

import "testing"

func TestIsSectionDelimiter(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		s string
		// Expected results.
		want bool
	}{
		{
			name: "Real delimiter",
			s:    "frontend https-in",
			want: true,
		},
		{
			name: "not a delimiter",
			s:    "blabla https-in",
			want: false,
		},
	}
	for _, tt := range tests {
		if got := isSectionDelimiter(tt.s); got != tt.want {
			t.Errorf("%q. isSectionDelimiter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetSectionName(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		s string
		// Expected results.
		want string
	}{
		{
			name: "Normal name",
			s:    "frontend https-in",
			want: "https-in",
		},
		{
			name: "Listen section",
			s:    "listen 127.0.0.1:8000",
			want: "127.0.0.1:8000",
		},
		{
			name: "No section name",
			s:    "global",
			want: "",
		},
	}
	for _, tt := range tests {
		if got := getSectionName(tt.s); got != tt.want {
			t.Errorf("%q. getSectionName() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetSection(t *testing.T) {
	lines, err := GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := GetSection("frontend", lines)

	if len(sections) != 2 {
		t.Errorf("Expected %d sections, got %d", 2, len(sections))
	}
}

func TestGetSections(t *testing.T) {
	lines, err := GetConfig("../testdata/haproxy.cfg", "")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := GetSections(lines)

	if len(sections) != 7 {
		t.Errorf("Expected %d sections, got %d", 7, len(sections))
	}
}
