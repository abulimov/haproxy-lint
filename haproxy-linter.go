package main

import (
	"flag"
	"fmt"

	"github.com/abulimov/haproxy-linter/checks"
	"github.com/abulimov/haproxy-linter/lib"
)

func main() {
	argFile := flag.String("file", "haproxy.cfg", "file to check")
	argJSON := flag.Bool("json", false, "Output in json")

	flag.Parse()

	config, err := lib.ReadConfigFile(*argFile)
	if err != nil {
		fmt.Println(err)
	}
	var sections []*lib.Section

	for _, s := range lib.Sections {
		ps := lib.GetSection(s, config)
		sections = append(sections, ps...)
	}

	problems := checks.Run(sections)

	if len(problems) != 0 {
		if *argJSON {
			fmt.Println(lib.ReportProblemsJSON(problems))
		} else {
			fmt.Print(lib.ReportProblems(problems))
		}
	}
}
