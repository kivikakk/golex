include $(GOROOT)/src/Make.inc

default: all

TARG=golex
GOFILES=\
	golex.go\

include $(GOROOT)/src/Make.cmd
