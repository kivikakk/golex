# golex
#### <span style="color: #333">flex-compatible lexical analyser generator</span>

## introduction

_golex_ is a [flex](http://flex.sourceforge.net)-compatible lexical analyser generator.

The below description has been pilfered from flex's description in Debian, adapted to describe _golex_:

_golex_ is a tool for generating scanners: programs which recognize lexical patterns in text. It reads the given input files for a description of a scanner to generate. The description is in the form of pairs of regular expressions and Go code, called rules. _golex_ generates as output a Go source file, which defines a routine `yylex()`. When the routine is run, it analyzes its input for occurrences of the regular expressions. Whenever it finds one, it executes the corresponding Go code.

## notes

_golex_ supports all features for regular expression matching as described in [flex's manual](http://flex.sourceforge.net/manual/Patterns.html#Patterns), _except_:

 * character class set operations `[a-z]{-}[aeiou]`
 * character class expressions `[:alnum:]`
 * regular expression option setting `(?is-x:pattern)`
 * matching EOF `<<EOF>>`

The above restrictions are intended to be removed in future releases of _golex_ (except, probably the regular expression option setting, as we're actually just using Go's built-in regular expressions for matching. the level of parsing required to allow such options is possibly above me).

A number of utility functions required for full flex emulation (mostly concerning manipulating the buffer [stack]) are also not yet available.

The full set of omissions (in regular expressions and otherwise) is detailed in the file `BUGS`.

_golex_, and the scanners it generates, are _not_ fast (unlike flex).

## examples

Examples taken from throughout the flex manual have been converted to Go and are included as `*.l` in this distribution.

To try one, if you're using `6g` et al., something like the following should work:

`make && ./golex file.l && 6g file.l.go && 6l file.l.6`

The binary `6.out` is now your scanner.

## other golexen

This is not the first attempt at writing a _golex_ utility, though it might be the first with the aim of behaving as similarly to the original flex as possible. Other golexen include (but are not limited to):

 * [Ben Lynn](http://cs.stanford.edu/~blynn/)'s [Nex](http://cs.stanford.edu/~blynn/nex/) tool.
 * [CZ.NIC](http://www.nic.cz)'s [lex](git://git.nic.cz/go/lex) package.
 * [CZ.NIC](http://www.nic.cz)'s [golex](git://git.nic.cz/go/golex) tool.

## license

Copyright 2011 Arlen Cuss

golex is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

golex is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with golex.  If not, see http://www.gnu.org/licenses/.

