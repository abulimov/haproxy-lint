package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ReportProblems returns all found problems as pretty string
func ReportProblems(problems []Problem) string {
	buffer := new(bytes.Buffer)
	for _, p := range problems {
		fmt.Fprintf(buffer, "%d:%d:%s: %s\n", p.Line, p.Col, p.Severity, p.Message)
	}
	return buffer.String()
}

//ReportProblemsJSON returns all found problems as pretty JSON string
func ReportProblemsJSON(problems []Problem) string {
	s, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return ""
	}
	return string(s)
}
