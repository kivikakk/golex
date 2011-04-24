package main

import (
	"fmt"
	"io"
	"bufio"
	"strings"
	"regexp"
	"container/list"
	goparser "go/parser"
	gotoken "go/token"
	goast "go/ast"
	goprinter "go/printer"
	"bytes"
)

type Parser struct {
	state     stateFunc
	out       *bufio.Writer

	inComment bool
	wroteEnd  bool

	parseSubs map[string]string
	firstPat  bool
	lastPat   string
	patStack  *list.List
}

func NewParser(out io.Writer) *Parser {
	return &Parser{state:     (*Parser).statePrologue,
	               out:       bufio.NewWriter(out),
		       inComment: false,
		       wroteEnd:  false,
		       parseSubs: make(map[string]string),
		       firstPat:  true,
		       lastPat:   "",
		       patStack:  list.New()}
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

	p.ParseFinish()
}

func (p *Parser) ParseFinish() {
	if !p.wroteEnd {
		p.out.WriteString("}\n")
	}

	p.out.WriteString(`
var yydata string = ""
var yyactionreturn bool = false

var yytext string = ""
var yytextrepl bool = true
func yymore() {
	yytextrepl = false
}

func yylex() int {
	reader := bufio.NewReader(yyin)

	for {
		line, err := reader.ReadString('\n')
		if len(line) == 0 && err == os.EOF {
			break
		}

		yydata += line
	}

	origData := yydata
	dataIndex := 0

	for len(yydata) > 0 {
		longestMatch, longestMatchLen := (func() int)(nil), -1
		longestAdvLen := -1

		for _, v := range yyrules {
			sol := dataIndex == 0 || origData[dataIndex-1] == '\n'

			// I want a real XOR! :(
			if v.sol && !sol {
				continue
			}

			idxs := v.regexp.FindStringIndex(yydata)
			if idxs != nil && idxs[0] == 0 {
				// Check the trailing context, if any.
				checksOk := true
				matchLen := idxs[1]
				subMatchLen := idxs[1]

				if v.trailing != nil {
					tridxs := v.trailing.FindStringIndex(yydata[idxs[1]:])
					if tridxs == nil || tridxs[0] != 0 {
						checksOk = false
					} else {
						matchLen += tridxs[1]
					}
				}

				if checksOk && matchLen > longestMatchLen {
					longestMatch, longestMatchLen = v.action, matchLen
					longestAdvLen = subMatchLen
				}
			}
		}

		if yytextrepl {
			yytext = ""
		}

		if longestMatch == nil {
			yytext += yydata[:1]
			yydata = yydata[1:]
			dataIndex += 1

			yyout.Write([]byte(yytext))
		} else {
			yytext += yydata[:longestAdvLen]
			yydata = yydata[longestAdvLen:]
			dataIndex += longestAdvLen

			yyactionreturn, yytextrepl = true, true
			rv := longestMatch()

			if yyactionreturn {
				return rv
			}
		}
	}

	return 0
}
`)

	p.out.Flush()
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

var (
	hexOrOctal *regexp.Regexp = regexp.MustCompile("\\\\\\\\([0-9][0-9][0-9]|[xX][0-9a-fA-F][0-9a-fA-F])")
	nulEscape  *regexp.Regexp = regexp.MustCompile("\\\\\\\\0($|[^0-9]|[0-9][^0-9])")
)

// quoteRegexp prepares a regular expression for insertion into a Go source
// as a string suitable for use as argument to regexp.(Must)?Compile.
func quoteRegexp(re string) string {
	re = strings.Replace(re, "\\", "\\\\", -1)
	re = strings.Replace(re, "\"", "\\\"", -1)
	re = hexOrOctal.ReplaceAllStringFunc(re, func(s string) string {
		var n int
		fmt.Sscan("0" + s[2:], &n)

		if n < 32 {
			s = fmt.Sprintf("\\x%02x", n)
		} else {
			s = string(n)
			s = strings.Replace(regexp.QuoteMeta(s), "\\", "\\\\", -1)
		}

		return s
	})
	re = nulEscape.ReplaceAllStringFunc(re, func(s string) string {
		s = "\\x00" + s[3:]
		return s
	})
	return re
}

type codeToActionVisitor struct { }
func (ctav *codeToActionVisitor) Visit(node goast.Node) goast.Visitor {
	// Transforms any lone expression statements where the expression is a lone ident
	// to a call of that name prefixed with yy (i.e. 'ECHO' -> 'yyECHO()').
	exprs, ok := node.(*goast.ExprStmt)
	if ok {
		rid, rok := exprs.X.(*goast.Ident)
		if rok {
			rid.Name = "yy" + rid.Name
			exprs.X = &goast.CallExpr{Fun: exprs.X,
						  Args: nil}
		}
	}

	return ctav
}

func codeToAction(code string) string {
	fs := gotoken.NewFileSet()

	expr, _ := goparser.ParseExpr(fs, "", "func() int {" + code + "; yyactionreturn = false; return 0}")

	fexp := expr.(*goast.FuncLit)

	ctav := &codeToActionVisitor{}
	goast.Walk(ctav, fexp)

	result := bytes.NewBuffer(make([]byte, 0, len(code) * 2))
	goprinter.Fprint(result, fs, fexp)

	return result.String()
}

// functions to handle each state

type stateFunc func(p *Parser, line string)

func (p *Parser) statePrologue(line string) {
	if line == "%%" {
		p.state = (*Parser).stateActions

		p.out.WriteString(`
			import (
				"regexp"
				"io"
				"bufio"
				"os"
			)

			var yyin io.Reader = os.Stdin
			var yyout io.Writer = os.Stdout
			type yyrule struct {
				regexp   *regexp.Regexp
				trailing *regexp.Regexp
				sol      bool
				action   func() int
			}
			var yyrules []yyrule = []yyrule{`)
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
		p.out.WriteString(strings.TrimSpace(line) + "\n")
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

		p.parseSubs[line[:smaller]] = strings.TrimSpace(line[smaller:])
	}
}

func (p *Parser) statePrologueLit(line string) {
	if line == "%}" {
		p.state = (*Parser).statePrologue
	} else {
		p.out.WriteString(line + "\n")
	}
}

func (p *Parser) stateActions(line string) {
	if line == "%%" {
		p.state = (*Parser).stateEpilogue
		p.wroteEnd = true
		p.out.WriteString("}\n")
		return
	}

	quotedPattern, trailingContext, remainder := p.ParseFlex(line)

	if trailingContext != "" {
		trailingContext = fmt.Sprintf("regexp.MustCompile(\"%s\")", quoteRegexp(trailingContext))
	} else {
		trailingContext = "nil"
	}

	quotedPattern = quoteRegexp(quotedPattern)

	p.lastPat = strings.TrimSpace(remainder)

	if len(p.lastPat) > 0 {
		if p.lastPat[0] == '{' {
			if p.lastPat[len(p.lastPat)-1] == '}' {
				p.lastPat = p.lastPat[:len(p.lastPat)-1]
			} else {
				p.state = (*Parser).stateActionsCont
			}

			p.lastPat = p.lastPat[1:]
		}
	}

	p.patStack.PushFront([]string{quotedPattern, trailingContext})

	if p.lastPat == "|" {
		return
	}

	for e := p.patStack.Front(); e != nil; e = p.patStack.Front() {
		if p.firstPat {
			p.firstPat = false
		} else {
			p.out.WriteString(",\n")
		}

		saved := e.Value.([]string)

		sol := "false"
		if len(saved[0]) > 0 && saved[0][0] == '^' {
			sol = "true"
		}

		p.out.WriteString(fmt.Sprintf("{regexp.MustCompile(\"%s\"), %s, %s, \n", saved[0], saved[1], sol))

		if p.state == (*Parser).stateActions {
			p.out.WriteString(codeToAction(p.lastPat))
			p.out.WriteString("}")
		}

		p.patStack.Remove(e)
	}
}

func (p *Parser) stateActionsCont(line string) {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == '}' {
		p.lastPat = strings.TrimSpace(p.lastPat + line)
		p.lastPat = p.lastPat[:len(p.lastPat)-1]

		p.out.WriteString(codeToAction(p.lastPat))
		p.out.WriteString("}")

		p.state = (*Parser).stateActions
	} else {
		p.lastPat += line + "\n"
	}
}

func (p *Parser) stateEpilogue(line string) {
	p.out.WriteString(line + "\n")
}

