package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "golex: Usage: %s [flags] file.l...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))
	if f == nil {
		fmt.Fprintf(os.Stderr, "golex: could not open '%s': %v\n", flag.Arg(0), err)
		os.Exit(1)
	}

	out, _ := os.Create(flag.Arg(0) + ".go")

	lf := ParseLexFile(f)
	lf.WriteGo(out)

	out.Close()

	f.Close()
}
