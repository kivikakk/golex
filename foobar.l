	package main
%%
foobar	ECHO; yyless(3)
[a-z]+	ECHO
%%
func main() { yylex() }
