package lib

type Problem struct {
	Col      int    `json:"col"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Entity struct {
	Line int
	Name string
}

type SectionCheck func(*Section) []Problem
type GlobalCheck func([]*Section) []Problem
