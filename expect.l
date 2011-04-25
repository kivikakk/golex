%{
package main
import "fmt"
%}
%s expect

%%
expect-floats	BEGIN(expect)

<expect>[0-9]+"."[0-9]+	{
		var f float32
		fmt.Sscanf(yytext, "%g", &f)
		fmt.Printf("found a float, = %f\n", f)
	}

<expect>\n	BEGIN(INITIAL)

[0-9]+	{
		var i int
		fmt.Sscanf(yytext, "%d", &i)
		fmt.Printf("found an integer, = %d\n", i)
	}

"."	fmt.Printf("found a dot\n")
%%
func main() { yylex() }
