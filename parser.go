package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"strings"
)

type Parser struct {
	state stateFunc

	seenPackage bool
	inComment   bool

	parseSubs map[string]string
	patStack  *list.List

	scNext int

	res              *LexFile
	curRule          LexRule
	curAction        string
	appendNextAction bool
	madeToEpilogue   bool
}

func ParseLexFile(in io.Reader) *LexFile {
	p := NewParser()
	p.ParseInput(in)
	return p.res
}

func NewParser() *Parser {
	return &Parser{state: (*Parser).statePrologue,
		seenPackage:      false,
		inComment:        false,
		parseSubs:        make(map[string]string),
		patStack:         list.New(),
		scNext:           1024,
		res:              NewLexFile(),
		curAction:        "",
		appendNextAction: false,
		madeToEpilogue:   false}
}

func (p *Parser) ParseInput(in io.Reader) {
	buffer := bufio.NewReader(in)

	for {
		line, err := buffer.ReadString('\n')

		if len(line) == 0 && err != nil {
			break
		} else if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		p.state(p, line)
	}

	if !p.madeToEpilogue {
		// Insert an artificial '%%' to ensure the last rule is written.
		p.state(p, "%%")
	}
}

func (p *Parser) trimComments(line string) string {
	if !p.inComment {
		idx := strings.Index(line, "/*")
		if idx != -1 {
			p.inComment = true
			trimmed := p.trimComments(line[idx:])
			return line[:idx] + trimmed
		}
		return line
	}

	// In comment.
	idx := strings.Index(line, "*/")

	if idx == -1 {
		p.inComment = true
		return ""
	}

	p.inComment = false
	return p.trimComments(line[idx+2:])
}

func parsePackage(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if len(line) > 8 && line[:8] == "package " {
		return line[8:], true
	}
	return "", false
}

// functions to handle each state

type stateFunc func(p *Parser, line string)

func (p *Parser) statePrologue(line string) {
	if line == "%%" {
		if !p.seenPackage {
			panic("no package statement seen")
		}

		p.state = (*Parser).stateActions
		return
	}

	line = p.trimComments(line)

	if len(strings.TrimSpace(line)) == 0 {
		return
	}

	if line == "%{" {
		p.state = (*Parser).statePrologueLit
		return
	}

	if line[0] == ' ' || line[0] == '\t' {
		if !p.seenPackage {
			if pack, ok := parsePackage(line); ok {
				p.res.packageName = pack
				p.seenPackage = true
			} else {
				p.res.prologue += strings.TrimSpace(line) + "\n"
			}
		} else {
			p.res.prologue += strings.TrimSpace(line) + "\n"
		}
	} else {
		firstSpace := strings.Index(line, " ")
		firstTab := strings.Index(line, "\t")
		if firstSpace == -1 && firstTab == -1 {
			panic(fmt.Sprintf("don't know what to do with line \"%s\" in PROLOGUE", line))
		}

		smaller := firstSpace
		if smaller == -1 {
			smaller = firstTab
		}
		if firstTab != -1 && firstTab < smaller {
			smaller = firstTab
		}

		key, val := line[:smaller], line[smaller:]

		// Is this an option, or substitution rule?
		if key[0] == '%' {
			switch key {
			case "%s", "%x":
				// Start conditions.
				conds := strings.Split(val, " ")
				for _, v := range conds {
					v = strings.TrimSpace(v)
					if len(v) == 0 {
						continue
					}

					p.res.startConditions[v] = LexStartCondition{p.scNext, key == "%x"}
					p.scNext++
				}
			}
		} else {
			p.parseSubs[key] = strings.TrimSpace(val)
		}
	}
}

func (p *Parser) statePrologueLit(line string) {
	if line == "%}" {
		p.state = (*Parser).statePrologue
	} else {
		if pack, ok := parsePackage(line); ok {
			p.res.packageName = pack
			p.seenPackage = true
		} else {
			p.res.prologue += strings.TrimSpace(line) + "\n"
		}
	}
}

func (p *Parser) stateActions(line string) {
	if line == "%%" {
		p.stateActions_Write()
		p.state = (*Parser).stateEpilogue
		p.madeToEpilogue = true
		return
	}

	if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
		if p.patStack.Len() == 0 {
			p.res.actionInline += line + "\n"
		} else {
			p.curAction += line + "\n"
		}
		return
	}

	if !p.appendNextAction {
		p.stateActions_Write()
	}

	startConditions, pattern, trailingPattern, remainder := p.ParseFlex(line)
	p.curRule = LexRule{startConditions: startConditions,
		pattern:         pattern,
		trailingPattern: trailingPattern}

	p.curAction = strings.TrimSpace(remainder)

	p.patStack.PushFront(p.curRule)

	p.appendNextAction = p.curAction == "|"
}

func (p *Parser) stateActions_Write() {
	for e := p.patStack.Front(); e != nil; e = p.patStack.Front() {
		saved := e.Value.(LexRule)
		saved.code = p.curAction
		p.res.rules = append(p.res.rules, saved)
		p.patStack.Remove(e)
	}
}

func (p *Parser) stateEpilogue(line string) {
	p.res.epilogue += line + "\n"
}
