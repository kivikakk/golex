	package main

	var line_num = 1

%x comment
%%

"/*"	BEGIN(comment)
<comment>[^*\n]*	/* eat up anything that's not a '*' */
<comment>"*"+[^*/\n]*	/* eat up '*'s not followed by '/'s */
<comment>\n		line_num++
<comment>"*"+"/"	BEGIN(INITIAL)
%%
func main() { yylex() }
