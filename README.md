quickimmut, a golang code generator
===================================

quickimmut spits out short little blots of golang source code which provide immutable maps and lists.

The output is on standard out, and should be piped to a file (or, use the `-output` flag).
A package name is not prefixed to the output; do that yourself first (this makes it easier to combine several outputs into one file).

The way I'd typically use this is by writing a file with its own generation instructions at the top:

```go
package zowwie

// Code generated by go:generate comments in this file.  DO NOT EDIT below the dashed line.
//
//go:generate sed -i /---/q thisfile.go
//go:generate quickimmut -output=thisfile.go list Foo
//go:generate quickimmut -output=thisfile.go map Foo Bar
//
// ---

// [stuff down here gets nuked and regenerated]
```


Immutability in golang is tricky
--------------------------------

Immutability in a language that doesn't have first class features in it is always tricky.

In Golang it's double-tricky because we don't have generics, so we can't really match the UX of the built-in maps nor lists,
and have to make a lot of interesting twists and turns to try to get reasonable ways to handle those.
(Maybe the in-progress Golang generics will solve this.  I don't know.)

If you'd like to poke about this yourself, here are some Golang Playground snippets I used as experiments:

- https://play.golang.org/p/XoGwVd8zOI0 -- an approach with unexported types... that doesn't work.
- https://play.golang.org/p/N51oP-xdaHZ -- an approach with varargs... that doens't work.
- https://play.golang.org/p/WO4cZW0e5i4 -- another approach with unexported types... that doesn't work.
- https://play.golang.org/p/a-DqVuGOfmn -- an approach for maps that **does** work -- and is the basis for what this tool generates.
- https://play.golang.org/p/YMDuCUCrf-4 -- an approach for lists that **does** work -- and is the basis for what this tool generates.


Should I use this?
------------------

This is definitely one of those "provided without warranty, nor guarantee of merchantability or fitness for a particular purpose" projects.

One might argue that golang isn't really meant to be good at immutability.  One might even carry that so far as arguing it's pointless to try.  One wouldn't be wrong, necessarily.

Use it if it feels like it's going to help you.  Don't if it doesn't.


Generics?
---------

This is pre-generics.

I have no idea if it will remain reasonable-looking after golang's generics are fully shipped.


License
-------

SPDX-License-Identifier: Apache-2.0 OR MIT