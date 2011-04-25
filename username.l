	package main

	import "fmt"
%%
username	fmt.Printf("%s", os.Getenv("USER"))
%%
func main() { yylex() }
