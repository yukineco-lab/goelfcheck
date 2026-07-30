package main

import (
	"flag"
	"fmt"
	"os"

	"goelfcheck/internal/inspect"
	"goelfcheck/internal/report"
)

func main() {
	verbose := flag.Bool("v", false, "show detailed ELF, Go, DWARF, symbol, and string scan information")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-v] <ELF file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	result, err := inspect.File(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "goelfcheck: %v\n", err)
		os.Exit(1)
	}

	report.Print(os.Stdout, result, *verbose)
	if result.HasFindings() {
		os.Exit(1)
	}
}
