package acl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abulimov/haproxy-lint/lib"
)

// ACL contains acl name and all possible usage attributes
type ACL struct {
	Name       string
	Predefined bool
	Negated    bool
	Inline     bool
}

func (a ACL) String() string {
	return fmt.Sprintf(
		"Name: %s, Predefined: %t, Negated: %t, Inline: %t",
		a.Name,
		a.Predefined,
		a.Negated,
		a.Inline,
	)
}

// IsInline checks if given acl name is inline acl
func IsInline(acl string) bool {
	re := regexp.MustCompile(`!?{[^}]+}`)
	return re.MatchString(acl)
}

// IsPredefined checks if acl name is predefined
func IsPredefined(acl string) bool {
	preDefinedACLs := map[string]bool{
		"FALSE":          true,
		"HTTP":           true,
		"HTTP_1.0":       true,
		"HTTP_1.1":       true,
		"HTTP_CONTENT":   true,
		"HTTP_URL_ABS":   true,
		"HTTP_URL_SLASH": true,
		"HTTP_URL_STAR":  true,
		"LOCALHOST":      true,
		"METH_CONNECT":   true,
		"METH_GET":       true,
		"METH_HEAD":      true,
		"METH_OPTIONS":   true,
		"METH_POST":      true,
		"METH_TRACE":     true,
		"RDP_COOKIE":     true,
		"REQ_CONTENT":    true,
		"TRUE":           true,
		"WAIT_END":       true,
	}
	_, predefined := preDefinedACLs[acl]
	return predefined
}

// ParseACLs returns slice of ACLs parsed from given string
func ParseACLs(line string) []*ACL {
	var acls []*ACL
	afterIfRE := regexp.MustCompile(`\w+\s+if\s+(.+)`)
	afterIfMatch := afterIfRE.FindAllStringSubmatch(line, -1)
	if len(afterIfMatch) > 0 {
		if len(afterIfMatch[0]) > 1 {
			afterIfString := afterIfMatch[0][1]
			word := regexp.MustCompile(`({[^}]+}|\w+|!)`)
			words := word.FindAllString(afterIfString, -1)
			negated := false
			for _, w := range words {
				switch w {
				case "!":
					negated = true
				case "or":
					continue
				default:
					acls = append(acls, &ACL{
						Name:       w,
						Predefined: IsPredefined(w),
						Negated:    negated,
						Inline:     IsInline(w),
					})
					negated = false
				}
			}
		}
	}

	return acls
}

// GetUsedACLNames returns list of acl names from string using acls
func GetUsedACLNames(line string) []string {
	var acls []string
	rawACLs := ParseACLs(line)
	for _, acl := range rawACLs {
		if !acl.Predefined {
			acls = append(acls, acl.Name)
		}
	}
	return acls
}

// GetNameFromDeclaration returns acl name from it's declaration
func GetNameFromDeclaration(l string) string {
	return lib.GetUsage("acl", l)
}

// LineUsesACL checks if given line uses given acl
func LineUsesACL(acl, line string) bool {
	return strings.Contains(line, acl)
}

// In checks if second list of acls contains first
func In(first, second []*ACL) bool {
	m := make(map[ACL]bool)
	for _, f := range first {
		m[*f] = false
	}
	for _, s := range second {
		if _, found := m[*s]; found {
			m[*s] = true
		}
	}
	for _, v := range m {
		if !v {
			return false
		}
	}

	return true
}
