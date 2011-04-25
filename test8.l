	package main


%s incl
%x excl
%%

	var enter_excl = false
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

<incl,excl>wut ECHO
<*>haa ECHO
<INITIAL,excl>sump ECHO

<*>.	/* munch */

%%
func main() {
	yylex()
}
