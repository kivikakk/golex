	package main
%%
"/*"	{
		for {
			c := input()
			for ; c != '*' && c != EOF; c = input() {}
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
%%
func main() { yylex() }
