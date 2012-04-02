package main

import (
	"container/list"
	"fmt"
	"regexp"
	"strings"
)

// Parses flex-style regular expressions.

type flexParser struct {
	p           *Parser
	line        string
	i           int
	stateFunc   func(fp *flexParser) bool
	qStart      int
	tcStart     int
	rangeStarts *list.List
	lastElement int
}

func (p *Parser) ParseFlex(line string) (startConds []string, expr string, trailing string, remainder string) {
	fp := &flexParser{p: p,
		line:        line,
		i:           0,
		stateFunc:   (*flexParser).stateRoot,
		qStart:      0,
		tcStart:     -1,
		rangeStarts: list.New(),
		lastElement: -1}

	if fp.line[0] == '<' {
		sce := strings.Index(fp.line, ">")
		scs := fp.line[1:sce]
		startConds = strings.Split(scs, ",")
		fp.line = fp.line[sce+1:]
	} else {
		startConds = []string{}
	}

	for ; fp.i < len(fp.line); fp.i++ {
		if fp.line[fp.i] == '\\' {
			fp.i++
			continue
		}

		if fp.stateFunc(fp) {
			break
		}
	}

	if fp.tcStart != -1 {
		expr = fp.line[:fp.tcStart]
		trailing = fp.line[fp.tcStart+1 : fp.i]
	} else {
		expr = fp.line[:fp.i]
		trailing = ""
	}

	remainder = fp.line[fp.i:]
	return
}

func (fp *flexParser) stateRoot() bool {
	switch fp.line[fp.i] {
	case ' ', '\t':
		return true
	case '[':
		fp.stateFunc = (*flexParser).stateClass
		fp.lastElement = fp.i
	case '"':
		fp.stateFunc = (*flexParser).stateQuotes
		fp.qStart = fp.i
	case '{':
		fp.stateFunc = (*flexParser).stateSubst
		fp.qStart = fp.i
	case '/':
		if fp.tcStart != -1 {
			panic("multiple trailing contexts '/'")
		}
		fp.tcStart = fp.i
	case '.':
		repl := "[^\\n]"
		fp.line = fp.line[:fp.i] + repl + fp.line[fp.i+1:]
		fp.lastElement = fp.i
		fp.i += len(repl) - 1
	case '^':
		if fp.i != 0 {
			// ^ to be treated as non-special if not at start
			// of fp.line
			fp.line = fp.line[:fp.i] + "\\^" + fp.line[fp.i+1:]
			fp.lastElement = fp.i
			fp.i += 1
		}
	case '$':
		if fp.tcStart != -1 {
			panic("unescaped '$' in pattern found after trailing context '/'")
		} else if fp.i != len(fp.line)-1 && fp.line[fp.i+1] != ' ' && fp.line[fp.i+1] != '\t' {
			// $ to be treated as non-special if not last char
			fp.line = fp.line[:fp.i] + "\\$" + fp.line[fp.i+1:]
			fp.lastElement = fp.i
			fp.i += 1
		} else {
			// last char.
			fp.tcStart = fp.i
			// fp.line[fp.i+1:] should be empty anyway.
			fp.line = fp.line[:fp.i] + "/\\n|$" + fp.line[fp.i+1:]
			fp.i += 5 - 1
		}
	case '(':
		if len(fp.line) > fp.i+3 && fp.line[fp.i+1:fp.i+3] == "?#" {
			// Regular expression comment.
			end := strings.Index(fp.line[fp.i:], ")")
			fp.line = fp.line[:fp.i] + fp.line[end+fp.i+1:]
			fp.i--
		} else {
			fp.rangeStarts.PushFront(fp.i)
		}
	case ')':
		f := fp.rangeStarts.Front()
		fp.lastElement = f.Value.(int)
		fp.rangeStarts.Remove(f)
	default:
		fp.lastElement = fp.i
	}
	return false
}

func (fp *flexParser) stateClass() bool {
	if fp.line[fp.i] == ']' {
		fp.stateFunc = (*flexParser).stateRoot
	}
	return false
}

func (fp *flexParser) stateQuotes() bool {
	if fp.line[fp.i] != '"' {
		return false
	}

	origQuoted := fp.line[fp.qStart+1 : fp.i]
	quoted := strings.Replace(origQuoted, "\\\"", "\"", -1)
	quoted = regexp.QuoteMeta(quoted)

	fp.line = fp.line[:fp.qStart] + quoted + fp.line[fp.i+1:]
	fp.i += len(quoted) - len(origQuoted) - 2

	fp.stateFunc = (*flexParser).stateRoot
	return false
}

func (fp *flexParser) stateSubst() bool {
	if fp.line[fp.i] != '}' {
		return false
	}

	name := fp.line[fp.qStart+1 : fp.i]
	repl, found := fp.p.parseSubs[name]
	if !found {
		panic(fmt.Sprintf("substitution {%s} found, but no such sub found!", name))
	}

	fp.line = fp.line[:fp.qStart] + "(" + repl + ")" + fp.line[fp.i+1:]
	fp.i += 2 + len(repl) - len(name) - 2

	fp.stateFunc = (*flexParser).stateRoot
	return false
}
