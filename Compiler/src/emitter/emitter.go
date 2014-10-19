package emitter

import (
	"fmt"
	"io"
)

import (
	"strings"
	"syntax"
	"tokens"
)

const (
	tab = " "
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
}

func (f *Data) Blank() { fmt.Fprint(f, "\n") }

func (f *Data) Comment(s string) { fmt.Fprintf(f, "\t/* %s */\n", s) }

func (f *Data) Header() {
	fmt.Fprintf(f, "// file:%s\n\n", f.Name)
}

func (f *Data) DotAlign(val int) { fmt.Fprintf(f, "\t.balign\t%d\n", val) }

func (f *Data) FuncHeader(name string) {
	fmt.Fprintf(f, "l_term* %s(vec_header* vecData) \n{\n", name)
	fmt.Fprintf(f, "%sstruct v_term* data = vecData.data;\n", tab)
	fmt.Fprintf(f, "%suint32_t length = vecData.size;\n", tab)
}

func (f *Data) FuncEnd(name string) {
	fmt.Fprintf(f, "} // %s\n\n", name)
}

type functionStack []*syntax.Function

func (s *functionStack) Empty() bool { return len(*s) == 0 }

func (s *functionStack) Push(fun *syntax.Function) {
	*s = append(*s, fun)
}

func (s *functionStack) Pop() (fun *syntax.Function) {
	ln := len(*s) - 1
	fun = (*s)[ln]
	*s = (*s)[:ln]
	return
}

func tabulation(depth int) string {
	return strings.Repeat(tab, depth)
}

func (f *Data) processSymbol(termNumber, depth int) {

	var startVar = fmt.Sprintf("start_%d", termNumber)
	var followVar = fmt.Sprintf("follow_%d", termNumber)
	var startValue = "0"
	var tabs = tabulation(depth)

	if termNumber != 0 {
		var prevFollow = fmt.Sprintf("follow_%d", termNumber-1)

		fmt.Fprintf(f, "%sif (%s >= length) /*Откат*/;\n", tabs, prevFollow)

		startValue = fmt.Sprintf("follow_%d", termNumber-1)
	}

	fmt.Fprintf(f, "%sint %s = %s = %s;\n", tabs, startVar, followVar, startValue)

	fmt.Fprintf(f, "%sif (data[%s]->tag == V_TERM_SYMBOL_TAG) %s++;\n", tabs, startVar, followVar)
	fmt.Fprintf(f, "%s\telse /*Откат*/;\n", tabs)

	fmt.Fprintf(f, "//--------------------------------------\n")
}

func (f *Data) processExpr(termNumber, nestedDepth int) {

}

func (f *Data) processPattern(p *syntax.Expr) {

	for termIndex, term := range p.Terms {

		switch term.TermTag {
		case syntax.VAR:
			switch term.Value.VarType {

			case tokens.VT_S:
				f.processSymbol(termIndex, termIndex+1)
				break

			case tokens.VT_E:
				f.processExpr(termIndex, termIndex+1)
				break
			}

			//fmt.Fprintf(f, "%s ", term.Value.VarType.String())
			break
		}
	}

	fmt.Fprintf(f, "\n")
}

func processFile(f Data) {
	unit := f.Ast

	f.Header()

	stack := make(functionStack, 0, len(unit.GlobMap)*2)
	for _, fun := range f.Ast.GlobMap {
		stack.Push(fun)
	}

	for !stack.Empty() {
		fun := stack.Pop()
		f.FuncHeader(fun.FuncName)

		for _, s := range fun.Sentences {
			f.processPattern(&s.Pattern)
			for _, a := range s.Actions {
				switch a.ActionOp {
				case syntax.COMMA: // ','
				case syntax.REPLACE: // '='
				case syntax.TARROW: // '->'
				case syntax.ARROW: // '=>'
				case syntax.COLON: // ':'
				case syntax.DCOLON: // '::'
				}
			}
		}

		f.FuncEnd(fun.FuncName)
	}
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
