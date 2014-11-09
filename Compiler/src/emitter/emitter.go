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

func (f *Data) Comment(s string) { fmt.Fprintf(f, "\t/* %s */\n", s) }

func (f *Data) Header() {
	fmt.Fprintf(f, "// file:%s\n\n", f.Name)
}

func (f *Data) funcHeader(name string) {
	fmt.Fprintf(f, "l_term* %s(vec_header* vecData) \n{\n", name)
	fmt.Fprintf(f, "%sstruct v_term* data = vecData.data;\n", tab)
	fmt.Fprintf(f, "%suint32_t length = vecData.size;\n", tab)
	fmt.Fprintf(f, "%suint32_t ok = 0;\n", tab)
}

func (f *Data) releaseOkVar(tabs string) {
	fmt.Fprintf(f, "%sok = 0;\n", tabs)
}

func (f *Data) setOkVar(tabs string) {
	fmt.Fprintf(f, "%sok = 1;\n", tabs)
}

func (f *Data) checkOKVar(tabs string) {
	fmt.Fprintf(f, "%sif (ok == 1)\n%s{", tabs)
}

func (f *Data) endBlock(tabs string)
{
	fmt.Fprintf(f, "%s}\n");
}

func (f *Data) FuncEnd(name string) {
	fmt.Fprintf(f, "} // %s\n\n", name)
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

	if len(p.Terms) == 0 {
		f.setOkVar(tabulation(1))
	} else {
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
	}

	fmt.Fprintf(f, "\n")
}

func (f *Data) processAction(*Action act ) {

	f.checkOKVar(tabulation(1))
	
	
	f.endBlock(tabulation(1));
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.Header()

	for _, fun := range unit.GlobMap {
		currFileData.funcHeader(fun.FuncName)

		for _, s := range fun.Sentences {
			currFileData.releaseOkVar(tabulation(1))
			currFileData.processPattern(&s.Pattern)
			for _, a := range s.Actions {
				switch a.ActionOp {
				case syntax.COMMA: // ','
				case syntax.REPLACE: // '='
					currFileData.processAction(a)
					break
				case syntax.TARROW: // '->'
				case syntax.ARROW: // '=>'
				case syntax.COLON: // ':'
				case syntax.DCOLON: // '::'
				}
			}
		}

		currFileData.FuncEnd(fun.FuncName)
	}
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
