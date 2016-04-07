package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abulimov/haproxy-lint/checks"
	"github.com/abulimov/haproxy-lint/lib"
)

var version = "0.2.1"

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] haproxy.cfg\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	argJSON := flag.Bool("json", false, "Output in json")
	versionFlag := flag.Bool("version", false, "print haproxy-lint version and exit")

	flag.Usage = myUsage

	flag.Parse()
	if *versionFlag {
		fmt.Printf("haproxy-lint version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filePath := flag.Args()[0]

	config, err := lib.ReadConfigFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sections := lib.GetSections(config)

	problems := checks.Run(sections)

	if len(problems) != 0 {
		if *argJSON {
			fmt.Println(lib.ReportProblemsJSON(problems))
		} else {
			fmt.Print(lib.ReportProblems(problems))
		}
	}
}
