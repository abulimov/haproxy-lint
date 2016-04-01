package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ReportProblems(problems []Problem) string {
	buffer := new(bytes.Buffer)
	for _, p := range problems {
		fmt.Fprintf(buffer, "%d:%d:%s: %s\n", p.Line, p.Col, p.Severity, p.Message)
	}
	return buffer.String()
}

func ReportProblemsJSON(problems []Problem) string {
	s, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return ""
	}
	return string(s)
}
