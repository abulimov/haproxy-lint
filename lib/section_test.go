package lib

import "testing"

func TestStripComments(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		s string
		// Expected results.
		want string
	}{
		{
			name: "String without comment",
			s:    "acl h_some hdr(Host) -i some.example.com",
			want: "acl h_some hdr(Host) -i some.example.com",
		},
		{
			name: "String with comment",
			s:    "acl h_some hdr(Host) -i some.example.com # comment",
			want: "acl h_some hdr(Host) -i some.example.com",
		},
		{
			name: "Just comment",
			s:    "#only comment",
			want: "",
		},
	}
	for _, tt := range tests {
		if got := StripComments(tt.s); got != tt.want {
			t.Errorf("%q. StripComments() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		needle string
		slice  []string
		// Expected results.
		want bool
	}{
		{
			name:   "string in slice",
			needle: "needle",
			slice:  []string{"some", "long", "needle", "ok"},
			want:   true,
		},
		{
			name:   "no string in slice",
			needle: "needle",
			slice:  []string{"some", "long", "sfneedle", "ok"},
			want:   false,
		},
		{
			name:   "empty slice",
			needle: "needle",
			slice:  []string{},
			want:   false,
		},
	}
	for _, tt := range tests {
		if got := stringInSlice(tt.needle, tt.slice); got != tt.want {
			t.Errorf("%q. stringInSlice() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

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
	lines, err := ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := GetSection("frontend", lines)

	if len(sections) != 2 {
		t.Errorf("Expected %d sections, got %d", 2, len(sections))
	}
}

func TestGetSections(t *testing.T) {
	lines, err := ReadConfigFile("../testdata/haproxy.cfg")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	sections := GetSections(lines)

	if len(sections) != 7 {
		t.Errorf("Expected %d sections, got %d", 7, len(sections))
	}
}
