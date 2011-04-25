	package main
%%
a	|
ab	|
abc	|
abcd	ECHO; REJECT;
.|\n	/* eat up unmatched character */
%%
func main() { yylex() }
