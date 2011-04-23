	package main

	import "fmt"
%%
^start	fmt.Printf("(found start)")
end$	fmt.Printf("(found end LOSE)")
e.d\n	fmt.Printf("(found e.d WIN)")
path($|ology\n)	{
		// The deal with this is that 'path$' should count as
		// a match length of 4 (not 5), and 'pathology\n' as 10.
		fmt.Printf("(found path(ology))")
	}
..th\n	fmt.Printf("(..th\\n)")

.	fmt.Printf("<%s>", yytext)
\n	fmt.Printf("(NL)\n")

%%
func main() {
	yylex()
}
