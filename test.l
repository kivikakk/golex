/* scanner for a toy Pascal-like language */

%{
	import "fmt"
	import "os"
%}

DIGIT	[0-9]
ID	[a-z][a-z0-9]*

%%

{DIGIT}+	{
			var v int
			fmt.Sscanf(yytext, "%d", v)
			fmt.Printf("An integer: %s (%d)\n", yytext, v)
		}

{DIGIT}+"."{DIGIT}*	{
			var v float
			fmt.Sscanf(yytext, "%f", v)
			fmt.Printf("A float: %s (%g)\n", yytext, v)
		}

if|then|begin|end|procedure|function	{
			fmt.Printf("A keyword: %s\n", yytext)
		}

{ID}		{
			fmt.Printf("An identifier: %s\n", yytext)
		}

"+"|"-"|"*"|"/"	{
			fmt.Printf("An operator: %s\n", yytext)
		}

"{"[\^{}}\n]*"}"	/* eat up one-line comments */

[ \t\n]+		/* eat up whitespace */

.			fmt.Printf("Unknown character: %s\n", yytext)

%%

func main() {
	yyin = os.Stdin

	yylex()
}
