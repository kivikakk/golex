	package main
%x incl
%%
include		BEGIN(incl)

[a-z]+		ECHO
[^a-z\n]*\n?	ECHO

<incl>[ \t]*	/* eat whitespace */
<incl>[^ \t\n]+	{ /* got the include name */
	yyin, _ = os.Open(yytext)
	if yyin == nil {
		panic("could not open file '" + yytext + "'")
	}

	yypush_buffer_state(yy_create_buffer(yyin, YY_BUF_SIZE))
	BEGIN(INITIAL)
	}

<<EOF>> {
		yypop_buffer_state()
		if YY_CURRENT_BUFFER == nil {
			yyterminate()
		}
	}
%%
func main() { yylex() }
