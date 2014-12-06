package emitter

import (
	"fmt"
	"io"
)

import (
	//	"strings"
	"syntax"
	"tokens"
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	CurrTermNum int
}

func (f *Data) mainFunc(depth int) {
	tabs := genTabs(depth + 1)

	fmt.Fprintf(f, "int main()\n{\n")
	fmt.Fprintf(f, "%s__initLiteralData();\n", tabs)
	fmt.Fprintf(f, "%smainLoop(Go);\n", tabs)
	fmt.Fprintf(f, "%sreturn 0;\n}\n", tabs)
}

func (f *Data) FuncDataMemoryAllocation(depth int, funcInfo *syntax.Function) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%sif (entryPoint == 0)\n", tabs)
	fmt.Fprintf(f, "%s{\n", tabs)
	fmt.Fprintf(f, "%s%senv.locals = (struct l_term*)malloc(%d * sizeof(struct l_term));\n", tabs, tab, 1)
	fmt.Fprintf(f, "%s%sfieldOfView.backups = (struct l_term_chain_t*)malloc(%d * sizeof(struct l_term_chain_t));\n", tabs, tab, 1)
	fmt.Fprintf(f, "%s}\n", tabs)
}

func (f *Data) FuncDataMemoryFree(depth int) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%sif (res != CALL_RESULT)\n", tabs)
	fmt.Fprintf(f, "%s{\n", tabs)
	fmt.Fprintf(f, "%s%sfree(env.locals);\n", tabs, tab)
	fmt.Fprintf(f, "%s%sfree(fieldOfView.backups);\n", tabs, tab)
	fmt.Fprintf(f, "%s}\n", tabs)
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

func (f *Data) processSymbol(termNumber, depth int) {

	var startVar = fmt.Sprintf("start_%d", termNumber)
	var followVar = fmt.Sprintf("follow_%d", termNumber)
	var startValue = "0"
	var tabs = genTabs(depth)

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
		f.setOkVar(genTabs(1))
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

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))

	f.PrintLabel(1, "%s}\n") //end block
}

func (f *Data) processFuncSentences(currFunc *syntax.Function) {
	currEntryPoint := 0

	f.funcHeader(currFunc.FuncName)
	f.FuncDataMemoryAllocation(1, currFunc)

	f.PrintLabel(1, "switch (entryPoint)\n") //case begin
	f.PrintLabel(1, "{\n")                   //case block begin
	f.PrintLabel(2, fmt.Sprintf("case %d: \n", currEntryPoint))

	for _, s := range currFunc.Sentences {
		f.releaseOkVar(genTabs(1))
		f.processPattern(&s.Pattern)
		for _, a := range s.Actions {
			switch a.ActionOp {

			case syntax.COMMA: // ','

			case syntax.REPLACE: // '='
				f.processAction(a)
				break

			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.COLON: // ':'
			case syntax.DCOLON: // '::'
			}
		}
	}

	f.PrintLabel(1, "} // case block end\n") //case block end
	f.FuncDataMemoryFree(1)
	f.PrintLabel(0, fmt.Sprintf("} // %s\n\n", currFunc.FuncName)) // func block end
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.PrintLabel(0, fmt.Sprintf("// file:%s\n\n", currFileData.Name))
	currFileData.PrintHeaders()

	for _, currFunc := range unit.GlobMap {
		currFileData.processFuncSentences(currFunc)
	}

	currFileData.initLiteralDataFunc(0)
	currFileData.mainFunc(0)
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
