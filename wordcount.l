	package main

	import "fmt"

	var word_count = 0
%%

frob		special(); REJECT
[^ \t\n]+	word_count++

%%

func special() {
	fmt.Printf("special() called\n")
}

func main() {
	yylex()
	fmt.Printf("word_count is %d\n", word_count)
}
