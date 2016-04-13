// Package lib contains core logic for performing checks on HAProxy config
package lib

// Problem is a struct desctibing typical problem our linter can find
type Problem struct {
	Col      int    `json:"col"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// Entity is used to contain line with its number
type Entity struct {
	Line int
	Name string
}

// SectionCheck is a function to run check against single Section
type SectionCheck func(*Section) []Problem

// GlobalCheck is a function to run check against all Sections
type GlobalCheck func([]*Section) []Problem
