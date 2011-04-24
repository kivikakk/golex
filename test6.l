	package main

	import "fmt"
%%
a	|
b	|
c	fmt.Printf("(a or b or c)")
^X	fmt.Printf("(leading X)")

^F	return 1

mega-	ECHO; yymore()
kludge	ECHO

w	|
wx	|
wxy	|
wxyz	if 4 == 4 { ECHO }; REJECT; fmt.Printf("reached after REJECT\n")
.|\n	/* eat */

%%
func main() {
	yylex()
}
