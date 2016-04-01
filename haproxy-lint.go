package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abulimov/haproxy-lint/checks"
	"github.com/abulimov/haproxy-lint/lib"
)

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] haproxy.cfg\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	argJSON := flag.Bool("json", false, "Output in json")
	flag.Usage = myUsage

	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filePath := flag.Args()[0]

	config, err := lib.ReadConfigFile(filePath)
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
