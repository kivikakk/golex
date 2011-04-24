	package main

	import "fmt"
%%
a	|
b	|
c	fmt.Printf("(a or b or c)")
^X	fmt.Printf("(leading X)")

mega-	ECHO; yymore()
kludge	ECHO

w	|
wx	|
wxy	|
wxyz	if 4 == 4 { ECHO }; REJECT
.|\n	/* eat */

%%
func main() {
	yylex()
}
