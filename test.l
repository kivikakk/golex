	package main

	import "fmt"
%%
[a-z]+		fmt.Printf("(X)")
foo[a-z]+	fmt.Printf("(foo)")
bar[a-z]+	fmt.Printf("(bar)")
.		
%%
func main() { yylex() }
