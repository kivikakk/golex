	package main

	import "fmt"

	var (
		num_lines = 0
		num_chars = 0
	)
%%
\n	num_lines++; num_chars++
.	num_chars++

%%
func main() {
	yylex()
	fmt.Printf("# of lines = %d, # of chars = %d\n",
		num_lines, num_chars)
}
