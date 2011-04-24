	package main

	import "fmt"
%%
\0	fmt.Printf("(got NUL)")
\0xyz	fmt.Printf("(got NUL xyz)")
\01xyz	fmt.Printf("(got NUL 1xyz)")
\011	fmt.Printf("(got 011)")
\x0a	fmt.Printf("(got x0a)")
\x0d	fmt.Printf("(got x0d)")
\123	fmt.Printf("(got S \\123)")
\133	fmt.Printf("(got [ \\133)")
\x(?# This is a regexp comment.)2a	fmt.Printf("(got * \\x2a)")

%%
func main() {
	yylex()
}
