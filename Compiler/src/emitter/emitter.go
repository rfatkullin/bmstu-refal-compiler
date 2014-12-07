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
	}

	fmt.Fprintf(f, "\n")
}

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))

	f.PrintLabel(1, "%s}\n") //end block
}

func (f *Data) ConstructResult(depth int, resultExpr syntax.Expr) {
	fragmentOffset := 0
	fragmentLength := 0
	currChainNum := 0
	stack := make([]int, 0)
	stackSize := 0
	termNumberInChain := 0

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "result.status = OK_RESULT;\n")
		f.PrintLabel(depth, "result.mainChain = 0;\n")
		f.PrintLabel(depth, "result.callChain = 0;\n")
	} else {

		fmt.Println(fmt.Sprintf("Terms count: %d", len(resultExpr.Terms)))
		f.PrintLabel(depth, "struct l_term* currTerm = 0;\n")
		f.PrintLabel(depth, "struct l_term* tmpTerm = 0;\n")
		f.PrintLabel(depth, fmt.Sprintf("struct l_term_chain_t* chain%d = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));\n", currChainNum))

		terms := make([]*syntax.Term, len(resultExpr.Terms))
		copy(terms, resultExpr.Terms)

		for len(terms) > 0 {

			term := terms[0]
			terms = terms[1:]

			switch term.TermTag {

			case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:

				if fragmentLength == 0 {

					fragmentOffset = term.Index
					f.PrintLabel(depth, "tmpTerm = currTerm;\n")
					f.PrintLabel(depth, "currTerm = (struct l_term*)malloc(sizeof(struct l_term));\n")
					f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;\n")
					f.PrintLabel(depth, "currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));\n")

					f.PrintLabel(depth, "if (tmpTerm != 0) {\n")
					f.PrintLabel(depth+1, "tmpTerm->next = currTerm;\n")
					f.PrintLabel(depth+1, "currTerm->prev = tmpTerm;\n")
					f.PrintLabel(depth, "}\n")
				}

				if termNumberInChain == 0 {
					f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->chain->begin = currTerm;\n", currChainNum))
				}

				fragmentLength++
				termNumberInChain++

				break
			case syntax.EXPR:

				if fragmentLength > 0 {
					f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;\n", fragmentOffset))
					f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;\n", fragmentLength))
					fragmentLength = 0
				}

				currChainNum++
				f.PrintLabel(depth, fmt.Sprintf("struct l_term* chainTerm%d = (struct l_term*)malloc(struct l_term);\n", currChainNum))
				f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->tag = L_TERM_CHAIN_TAG;\n", currChainNum))
				f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->chain = (struct l_term_chain_t*)malloc(struct l_term_chain_t);\n", currChainNum))

				if termNumberInChain == 0 {
					f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->chain->begin = chainTerm%d;\n", currChainNum-1, currChainNum))
					f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->prev = 0;\n", currChainNum))
				} else {
					f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->prev = currTerm;\n", currChainNum))
					f.PrintLabel(depth, fmt.Sprintf("currTerm->next = chainTerm%d;\n", currChainNum))
				}

				termNumberInChain = 0

				if len(stack) <= stackSize {
					stack = append(stack, len(term.Exprs[0].Terms))
				} else {
					stack[stackSize] = len(term.Exprs[0].Terms)
				}

				stackSize++

				tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
				tmpTerms = append(tmpTerms, terms...)
				terms = tmpTerms
				break

			case syntax.EVAL:
				//TO DO
				break

			case syntax.FUNC:
				//TO DO
				break

			case syntax.BRACED_EXPR:
			case syntax.BRACKETED_EXPR:
			case syntax.ANGLED_EXPR:
				//Пока считаем, что тут не может быть литералов
				break

			case syntax.VAR:
			case syntax.L:
			case syntax.R:
				//Не литералы
				break
			}

			// If prev term was last term in expr
			for termNumberInChain > 0 && stackSize > 0 {

				stack[stackSize-1]--

				if stack[stackSize-1] > 0 {
					break
				}

				if fragmentLength > 0 {
					f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;\n", fragmentOffset))
					f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;\n", fragmentLength))
					fragmentLength = 0
				}

				f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->chain->end = currTerm;\n", currChainNum))
				f.PrintLabel(depth, fmt.Sprintf("currTerm = chainTerm%d;\n", currChainNum))
				currChainNum--
				stackSize--
			}
		}
	}

	f.PrintLabel(depth, fmt.Sprintf("chainTerm%d->end = currTerm;\n", currChainNum))
}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	currEntryPoint := 0

	f.funcHeader(currFunc.FuncName)
	f.PrintLabel(depth, "struct fresult_t result;\n")
	f.FuncDataMemoryAllocation(depth, currFunc)

	f.PrintLabel(depth, "switch (entryPoint)\n") //case begin
	f.PrintLabel(depth, "{\n")                   //case block begin

	for _, s := range currFunc.Sentences {

		f.PrintLabel(depth+1, fmt.Sprintf("case %d: \n", currEntryPoint))
		f.processPattern(depth+2, &s.Pattern)

		for _, a := range s.Actions {
			switch a.ActionOp {

			case syntax.REPLACE: // '='
				f.ConstructResult(depth+2, a.Expr)
				currEntryPoint++
				break

			case syntax.COLON: // ':'
				currEntryPoint++
				f.PrintLabel(depth+2, fmt.Sprintf("break;\n"))
				f.PrintLabel(depth+1, fmt.Sprintf("case %d: \n", currEntryPoint))
				break

			case syntax.COMMA: // ','
			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.DCOLON: // '::'
			}
		}
	}

	f.PrintLabel(depth+1, fmt.Sprintf("break;\n")) // last case break
	f.PrintLabel(1, "} // case block end\n")       //case block end
	f.FuncDataMemoryFree(1)
	f.PrintLabel(0, fmt.Sprintf("} // %s\n\n", currFunc.FuncName)) // func block end
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.PrintLabel(0, fmt.Sprintf("// file:%s\n\n", currFileData.Name))
	currFileData.PrintHeaders()

	for _, currFunc := range unit.GlobMap {
		currFileData.processFuncSentences(1, currFunc)
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
