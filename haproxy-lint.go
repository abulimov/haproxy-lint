// Simple linter for HAProxy configuration.
//
// This tool can perform static analysis of HAProxy config.
//
// Usage:
//
//   $ haproxy-lint [options] FILE
//
// Available options:
//
//   -json          output found problems in JSON format
//
//   -run-haproxy   run haproxy binary in check mode and parse
//                  its output to extract alerts and warnigns
//
//   -ignore regexp ignore lines in config file matching given regexp
//                  (works only for native checks)
//
//   -v             show version and exit
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/abulimov/haproxy-lint/checks"
	"github.com/abulimov/haproxy-lint/lib"
)

var version = "0.4.1"

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] haproxy.cfg\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output in json")
	haproxyFlag := flag.Bool("run-haproxy", true, "Try to run HAProxy binary in check mode")
	versionFlag := flag.Bool("version", false, "print haproxy-lint version and exit")
	ignoreFlag := flag.String("ignore", "", "ignore lines in config file matching this regexp")

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

	var problems []lib.Problem
	useHAProxy := *haproxyFlag
	if useHAProxy {
		haproxyProblems, err := lib.RunHAProxyCheck(filePath)
		if err != nil {
			log.Println(err)
			useHAProxy = false
		} else {
			problems = append(problems, haproxyProblems...)
		}
	}

	config, err := lib.ReadConfigFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// if we need to strip some strings from config (for example, Jinja2 conditionals)
	if *ignoreFlag != "" {
		config = lib.Filter(config, *ignoreFlag)
	}

	sections := lib.GetSections(config)

	// if we have local haproxy executable we shouldn't run
	// checks that are implemented in haproxy itself.
	nativeProblems := checks.Run(sections, useHAProxy)

	problems = append(problems, nativeProblems...)

	if len(problems) != 0 {
		if *jsonFlag {
			fmt.Println(lib.ReportProblemsJSON(problems))
		} else {
			fmt.Print(lib.ReportProblems(problems))
		}
	}
}
