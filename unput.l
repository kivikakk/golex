	package main
%%
xyz	{
		unput(')')
		for i := len(yytext)-1; i >= 0; i-- {
			unput(yytext[i])
		}
		unput('(')
	}
\(.	ECHO /* so that an unputted "(xyz)" doesn't get re-caught */
%%
func main() { yylex() }
