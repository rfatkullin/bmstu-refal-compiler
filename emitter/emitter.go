package emitter

import (
	"fmt"
	"io"
)

import (
	"BMSTU-Refal-Compiler/syntax"
	"BMSTU-Refal-Compiler/tokens"
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	CurrTermNum int
}

func (f *Data) mainFunc(depth int) {

	f.PrintLabel(depth, "int main()")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "__initLiteralData();")
	f.PrintLabel(depth+1, "mainLoop(Go);")
	f.PrintLabel(depth+1, "return 0;")
	f.PrintLabel(depth, "}")
}

func (f *Data) FuncDataMemoryAllocation(depth int, funcInfo *syntax.Function) {

	f.PrintLabel(depth, "struct func_result_t funcRes;")
	f.PrintLabel(depth, "if (entryPoint == 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals = (struct lterm_t*)malloc(%d * sizeof(struct lterm_t));", 1))
	f.PrintLabel(depth+1, fmt.Sprintf("fieldOfView->backups = (struct lterm_chain_t*)malloc(%d * sizeof(struct lterm_chain_t));", 1))
	f.PrintLabel(depth, "}")
}

func (f *Data) FuncDataMemoryFree(depth int) {

	f.PrintLabel(depth, "if (funcRes.status != CALL_RESULT)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "free(env->locals);")
	f.PrintLabel(depth+1, "free(fieldOfView->backups);")
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "return funcRes;")
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
	fmt.Fprintf(f, "%s\telse /*Откат*/;", tabs)

	fmt.Fprintf(f, "//--------------------------------------\n")
}

func (f *Data) processExpr(termNumber, nestedDepth int) {

}

func (f *Data) processPattern(depth int, p *syntax.Expr) {

	if len(p.Terms) != 0 {
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
}

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))
	f.PrintLabel(1, "%s}") //end block
}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	currEntryPoint := 0

	f.funcHeader(currFunc.FuncName)
	f.FuncDataMemoryAllocation(depth, currFunc)

	f.PrintLabel(depth, "switch (entryPoint)") //case begin
	f.PrintLabel(depth, "{")                   //case block begin

	for _, s := range currFunc.Sentences {

		f.PrintLabel(depth+1, fmt.Sprintf("case %d:", currEntryPoint))
		f.PrintLabel(depth+1, fmt.Sprintf("{"))
		f.processPattern(depth+2, &s.Pattern)

		for _, a := range s.Actions {
			switch a.ActionOp {

			case syntax.REPLACE: // '='
				f.ConstructResult(depth+2, a.Expr)
				currEntryPoint++
				break

			case syntax.COLON: // ':'
				currEntryPoint++
				f.PrintLabel(depth+2, fmt.Sprintf("break;"))
				f.PrintLabel(depth+1, fmt.Sprintf("case %d: ", currEntryPoint))
				f.PrintLabel(depth+1, fmt.Sprintf("{"))
				break

			case syntax.COMMA: // ','
			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.DCOLON: // '::'
			}
		}
	}

	f.PrintLabel(depth+2, fmt.Sprintf("break;")) // last case break
	f.PrintLabel(depth+1, fmt.Sprintf("}"))      // last case }
	f.PrintLabel(1, "} // switch block end")     //switch block end
	f.FuncDataMemoryFree(1)
	f.PrintLabel(0, fmt.Sprintf("} // %s\n", currFunc.FuncName)) // func block end
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.PrintLabel(0, fmt.Sprintf("// file:%s\n", currFileData.Name))
	currFileData.PrintHeaders()

	currFileData.initLiteralDataFunc(0)

	for _, currFunc := range unit.GlobMap {
		currFileData.processFuncSentences(1, currFunc)
	}

	currFileData.mainFunc(0)
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
