include $(GOROOT)/src/Make.inc

default: all

TARG=golex
GOFILES=\
	golex.go\
	parser.go\
	regexp.go\

include $(GOROOT)/src/Make.cmd
