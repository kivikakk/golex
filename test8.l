	package main

%s incl
%x excl
%%

	if enter_excl {
		BEGIN(excl)
	}

i	BEGIN(incl)
x	BEGIN(excl)

a	ECHO
<incl>b	ECHO
<incl>B	ECHO; BEGIN(INITIAL)
<excl>c	ECHO
<excl>C	ECHO; BEGIN(INITIAL)

%%
func main() {
	yylex()
}
