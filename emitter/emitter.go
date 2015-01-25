package emitter

import (
	"fmt"
	"io"
)

import (
	"BMSTU-Refal-Compiler/syntax"
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

func (f *Data) funcDataMemoryAllocation(depth int, funcInfo *syntax.Function) {

	f.PrintLabel(depth, "struct func_result_t funcRes;")
	//f.PrintLabel(depth, "if (entryPoint == 0)")
	//f.PrintLabel(depth, "{")
	//f.PrintLabel(depth+1, fmt.Sprintf("env->locals = (struct lterm_t*)malloc(%d * sizeof(struct lterm_t));", 1))
	//f.PrintLabel(depth+1, fmt.Sprintf("fieldOfView->backups = (struct lterm_chain_t*)malloc(%d * sizeof(struct lterm_chain_t));", 1))
	//f.PrintLabel(depth, "}")
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

func (f *Data) processExpr(termNumber, nestedDepth int) {

}

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))
	f.PrintLabel(1, "%s}") //end block
}

func (f *Data) initSentenceLocalVariables(depth, varsNumber int) {

	f.PrintLabel(depth, fmt.Sprintf("env->locals = (struct lterm_t*)malloc(%d * sizeof(struct lterm_t));", varsNumber))
	f.PrintLabel(depth, "int i = 0;")
	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; i++)", varsNumber))
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "env->locals[i].tag = L_TERM_FRAGMENT_TAG;")
	f.PrintLabel(depth, "}")
}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	currEntryPoint := 0

	f.funcHeader(currFunc.FuncName)
	f.funcDataMemoryAllocation(depth, currFunc)

	f.PrintLabel(depth, "while(1)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (entryPoint)") //case begin
	f.PrintLabel(depth+1, "{")                   //case block begin

	for _, s := range currFunc.Sentences {

		f.PrintLabel(depth+2, fmt.Sprintf("case %d:", currEntryPoint))
		f.PrintLabel(depth+2, fmt.Sprintf("{"))
		f.initSentenceLocalVariables(depth+3, s.VarsNumber)

		f.matchingPattern(depth+3, &s.Pattern, &s.Scope)

		for _, a := range s.Actions {
			switch a.ActionOp {

			case syntax.REPLACE: // '='
				f.ConstructResult(depth+3, a.Expr)
				currEntryPoint++
				break

			case syntax.COLON: // ':'
				f.matchingPattern(depth+3, &a.Expr, &s.Scope)
				currEntryPoint++
				f.PrintLabel(depth+3, fmt.Sprintf("break;"))
				f.PrintLabel(depth+2, fmt.Sprintf("case %d: ", currEntryPoint))
				f.PrintLabel(depth+2, fmt.Sprintf("{"))
				break

			case syntax.COMMA: // ','
			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.DCOLON: // '::'
			}
		}
	}

	f.PrintLabel(depth+3, fmt.Sprintf("break;")) // last case break
	f.PrintLabel(depth+2, fmt.Sprintf("}"))      // last case }
	f.PrintLabel(depth+1, "} // switch end")
	f.PrintLabel(depth, "} // while end")
	f.FuncDataMemoryFree(depth)
	f.PrintLabel(depth, fmt.Sprintf("} // %s\n", currFunc.FuncName)) // func block end
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
