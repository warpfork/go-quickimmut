package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	output     = flag.String("output", "", "set to a file name you'd like the output be appended to instead of stdout")
	visibility = flag.Bool("exp", false, "set to true if you want carrier types to be exported symbols")
	attachTo   = flag.String("attach", "", "set to a type name if you want constructor funcs to be attached to it as methods")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, "first argument should be 'map' or 'list' or 'struct'\n")
		fmt.Fprint(os.Stderr, `
example usage:

	quickimmut map Key Valu
	quickimmut map '*Key' '*Valu'
	quickimmut list Valu
	quickimmut list '*Valu'
	quickimmut -exp=true -attach=GroupingType map Key Valu

Type names for keys and values can be exported or non-exported,
and can be pointers or non-pointers,
but can't be anonymous types (e.g. not golang maps nor slices nor the clause 'interface{}', etc).
There is no validation of this; you'll just get nonsensical code out if providing nonsensical input.

Golang syntax is emitted on stdout.
This can be composed into a complete file with shell:

	quickimmut map Key Valu | cat "package foo" - > gen.go

	(quickimmut map Key Valu ; quickimmut list Wow) | cat "package foo" - > gen.go

Alternatively, the output flag can be used to append to a file
(which can be useful if using go:generate annotations, since they don't support shell-like redirection):

	//go:generate -output=file.go

Note that the output flag _appends_; it will not create the file if it does not exist.
(To result in valid golang, the file should already contain the "package" declaration anyway.)

`)
		os.Exit(2)
	}
	out := os.Stdout
	if *output != "" {
		out2, err := os.OpenFile(*output, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		out = out2
	}
	switch args[0] {
	case "map":
		if len(args) < 3 {
			fmt.Fprint(os.Stderr, "generating a map requires two further arguments (the key type and the value type)\n")
			os.Exit(2)
		}
		mapInfo{args[1], args[2], *visibility, *attachTo}.Eval(out)
	case "list":
		if len(args) < 2 {
			fmt.Fprint(os.Stderr, "generating a list requires one further argument (the value type)\n")
			os.Exit(2)
		}
		listInfo{args[1], *visibility, *attachTo}.Eval(out)
	case "struct":
		// It turns out in almost every case I have today, the structs have at least one detail that's "special".
		// In terms of costs, grinding out methods on one type is also a lot less irritating than the multi-type dance.
		// As a result, implementing read-only struct generation it hasn't been worth it (to me; yet).
		fmt.Fprint(os.Stderr, "sorry, struct feature not actually implemented.  PRs maybe welcome.\n")
		os.Exit(2)
	default:
		fmt.Fprint(os.Stderr, "first argument should be 'map' or 'list' or 'struct'\n")
		os.Exit(2)
	}
}
