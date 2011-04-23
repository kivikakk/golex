	package main

	import "fmt"
%%
^start	fmt.Printf("(found start)")
end$	fmt.Printf("(found end)")
.	fmt.Printf("<%s>", yytext)
\n	fmt.Printf("(NL)\n")

%%
func main() {
	yylex()
}
