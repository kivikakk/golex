/* scanner for a toy Pascal-like language */

%{
	package main

	import "fmt"
	import "flag"
%}

DIGIT    [0-9]
ID       [a-z][a-z0-9]*

%%

{DIGIT}+	{
			var i int
			fmt.Sscanf(yytext, "%d", &i)
			fmt.Printf("An integer: %s (%d)\n", yytext, i)
		}

{DIGIT}+"."{DIGIT}*	{
			var i float32
			fmt.Sscanf(yytext, "%g", &i)
			fmt.Printf("A float: %s (%g)\n", yytext, i)
		}

if|then|begin|end|procedure|function	{
			fmt.Printf("A keyword: %s\n", yytext)
		}

{ID}			fmt.Printf("An identifier: %s\n", yytext)

"+"|"-"|"*"|"/"		fmt.Printf("An operator: %s\n", yytext)

"{"[^{}\n]*"}"		/* eat up one-line comments */

[ \t\n]+		/* eat up whitespace */

.			fmt.Printf("Unrecognized character: %s\n", yytext)

%%

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		yyin, _ = os.Open(flag.Arg(0))
	} else {
		yyin = os.Stdin
	}

	yylex()
}
