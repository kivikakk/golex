	package main
%%

foobar	ECHO; yyless(3)
[a-z]+	ECHO

456	{
		unput(')')
		for i := len(yytext)-1; i >= 0; i-- {
			unput(yytext[i])
		}
		unput('(')
	}

"/*"	{
		for {
			var c int
			for c = input(); c != '*' && c != EOF; c = input() {}
			if c == '*' {
				for c = input(); c == '*'; c = input() {}
				if c == '/' {
					break
				}
			}
			if c == EOF {
				panic("EOF in comment")
			}
		}
	}

\(.*\)	ECHO

%%
func main() {
	yylex()
}
