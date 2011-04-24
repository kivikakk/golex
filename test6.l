	package main

	import "fmt"
%%
a	|
b	|
c	fmt.Printf("(a or b or c)")
^X	fmt.Printf("(leading X)")

%%
func main() {
	yylex()
}
