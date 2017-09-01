# golex
#### <span style="color: #333">flex-compatible lexical analyser generator</span>

## introduction

_golex_ is a [flex](http://flex.sourceforge.net)-compatible lexical analyser generator, written for Go 1.

The below description has been pilfered from flex's description in Debian, adapted to describe _golex_:

_golex_ is a tool for generating scanners: programs which recognize lexical patterns in text. It reads the given input files for a description of a scanner to generate. The description is in the form of pairs of regular expressions and Go code, called rules. _golex_ generates as output a Go source file, which defines a routine `yylex()`. When the routine is run, it analyzes its input for occurrences of the regular expressions. Whenever it finds one, it executes the corresponding Go code.

## notes

_golex_ supports all features for regular expression matching as described in [flex's manual](http://flex.sourceforge.net/manual/Patterns.html#Patterns), _except_:

 * character class set operations `[a-z]{-}[aeiou]`, and
 * matching EOF `<<EOF>>`.

EOF-matching is intended to be added to a future release of _golex_. Character class operations, however, will not, unless Go's own regular expression library (based on [RE2](http://code.google.com/p/re2/)) comes to.

A number of utility functions required for full flex emulation (mostly concerning manipulating the buffer (stack)) are also not yet available.

The full set of omissions (in regular expressions and otherwise) is detailed in the GitHub Issues for this repository.

_golex_ and the scanners it generates are _not_ fast (unlike those of flex).  Rather than implementing its own regular expression engine and crafting a state machine based on that, _golex_ simply defers to Go's built-in regular expressions, and matches character-by-character.  Pull requests to right this wrong gratefully accepted! :)

## examples

Self-contained examples, taken from throughout the flex manual, have been converted to Go and are included as `*.l` in this distribution.  I invite you to compare them to the original _flex_ examples to note how similar they are. [Here are a few examples](http://flex.sourceforge.net/manual/Simple-Examples.html#Simple-Examples), found as `username.l`, `counter.l`, and `toypascal.l` in the _golex_ distribution.

A `test` script for building and running an example is included.  For example:

`./test toypascal.l`

will build _golex_, run _golex_ on `toypascal.l`, build the resulting Go code, and then run the resulting lexer.

## other golexen

This is not the first attempt at writing a _golex_ utility, though it might be the first with the aim of behaving as similarly to the original _flex_ as possible.

Other golexen include (but are not limited to):

 * [Ben Lynn](http://cs.stanford.edu/~blynn/)'s [Nex](http://cs.stanford.edu/~blynn/nex/) tool.
 * [CZ.NIC](http://www.nic.cz)'s package at `git://git.nic.cz/go/lex`.
 * [CZ.NIC](http://www.nic.cz)'s tool at `git://git.nic.cz/go/golex` (it's not like the name is terribly original!).

## license

Copyright 2011-2017 Ashe Connor. Licensed under the [2-Clause BSD License](https://opensource.org/licenses/BSD-2-Clause).
