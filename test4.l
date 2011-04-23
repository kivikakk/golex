	package main

	import "fmt"
%%
trai.		fmt.Printf("('trai.' match)")
trail/ing	fmt.Printf("('trail/ing' match)")
[ \t]+	fmt.Printf(" ")
[ \t]+$	/* ignore */

%%
func main() {
	yylex()
}
