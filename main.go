package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "first argument should be 'map' or 'list' or 'struct'\n")
		fmt.Fprint(os.Stderr, `
example usage:

	quickimmut map Key Valu
	quickimmut map '*Key' '*Valu'
	quickimmut list Valu
	quickimmut list '*Valu'

Golang syntax is emitted on stdout.
This can be composed into a complete file with shell:

	quickimmut map Key Valu | cat "package foo" - > gen.go

	(quickimmut map Key Valu ; quickimmut list Wow) | cat "package foo" - > gen.go

`)
		os.Exit(1)
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case "map":
		mapInfo{args[0], args[1], false, "Compiler"}.Eval(os.Stdout)
	case "list":
		listInfo{args[0], false, "Compiler"}.Eval(os.Stdout)
	case "struct":
		// It turns out in almost every case I have today, the structs have at least one detail that's "special".
		// In terms of costs, grinding out methods on one type is also a lot less irritating than the multi-type dance.
		// As a result, implementing read-only struct generation it hasn't been worth it (to me; yet).
		fmt.Fprint(os.Stderr, "sorry, struct feature not actually implemented.  PRs maybe welcome.\n")
		os.Exit(1)
	default:
		fmt.Fprint(os.Stderr, "first argument should be 'map' or 'list' or 'struct'\n")
		os.Exit(1)
	}
}
