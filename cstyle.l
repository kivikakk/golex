	package main

	import "fmt"

	var string_buf = ""
%x str
%%
\"	string_buf = ""; BEGIN(str)

<str>\"		{
			/* closing quote - all done */
			BEGIN(INITIAL)
			fmt.Printf("Got full string: \"%s\"\n", string_buf)
		}

<str>\n		panic("unterminated string constant")

<str>\\[0-7]{1,3}	{ 
			/* octal escape sequence */
			var result int
			fmt.Sscanf(yytext[1:], "%o", &result)
			if result > 0xff {
				panic("octal escape sequence out of bounds")
			}
			string_buf += string(result)
		}

<str>\\[0-9]+	panic("bad escape sequence '" + yytext + "'")

<str>\\n	string_buf += "\n"
<str>\\t	string_buf += "\t"
<str>\\r	string_buf += "\r"
<str>\\b	string_buf += "\b"
<str>\\f	string_buf += "\f"

<str>\\(.|\n)	string_buf += yytext[1:2]

<str>[^\\\n\"]+	string_buf += yytext

%%
func main() { yylex() }
