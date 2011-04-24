	package main

	import "fmt"
%%
a	|
b	|
c	fmt.Printf("(a or b or c)")

%%
func main() {
	yylex()
}
