package emitter

import (
	"fmt"
	"io"
)

import (
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

	fmt.Fprintf(f, "%sstruct func_result_t funcRes;\n", tabs)
	fmt.Fprintf(f, "%sif (entryPoint == 0)\n", tabs)
	fmt.Fprintf(f, "%s{\n", tabs)
	fmt.Fprintf(f, "%s%senv->locals = (struct lterm_t*)malloc(%d * sizeof(struct lterm_t));\n", tabs, tab, 1)
	fmt.Fprintf(f, "%s%sfieldOfView->backups = (struct lterm_chain_t*)malloc(%d * sizeof(struct lterm_chain_t));\n", tabs, tab, 1)
	fmt.Fprintf(f, "%s}\n", tabs)
}

func (f *Data) FuncDataMemoryFree(depth int) {

	f.PrintLabel(depth, "if (funcRes.status != CALL_RESULT)\n")
	f.PrintLabel(depth, "{\n")
	f.PrintLabel(depth+1, "free(env->locals);\n")
	f.PrintLabel(depth+1, "free(fieldOfView->backups);\n")
	f.PrintLabel(depth, "}\n")
	f.PrintLabel(depth, "return funcRes;\n")
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

		fmt.Fprintf(f, "\n")
	}
}

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))

	f.PrintLabel(1, "%s}\n") //end block
}

func Max(a int, b int) int {

	if a > b {
		return a
	}

	return b
}

func (f *Data) CalcChainsCount(expr syntax.Expr) int {
	chainsCount := 1
	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	for len(terms) > 0 {
		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.EXPR, syntax.EVAL:
			tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
			tmpTerms = append(tmpTerms, terms...)
			terms = tmpTerms

			chainsCount++
			break
		}
	}

	return chainsCount
}

func (f *Data) ConstructFragmentLTerm(depth int, firstTerm bool, chainNumber int, fragmentOffset int, fragmentLength int) {

	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;\n")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;\n", fragmentOffset))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;\n", fragmentLength))

	//Самый первый терм в цепочке.
	if firstTerm {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->begin = currTerm;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", chainNumber))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end->next = currTerm;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain->end;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", chainNumber))
	}
}

func (f *Data) ConstructFuncCall(depth int, firstFuncCall bool, funcName string, chainNumber int) {

	f.PrintLabel(depth, "funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
	f.PrintLabel(depth, "funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->funcName = memMngr.termsHeap[helper[%d]->chain->begin->fragment->offset].str;\n", chainNumber))
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->funcPtr = %s;\n", funcName))
	f.PrintLabel(depth, "funcTerm->funcCall->entryPoint = 0;\n")
	f.PrintLabel(depth, "funcTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->fieldOfView->current = helper[%d]->chain;\n", chainNumber))
	f.PrintLabel(depth, "funcTerm->tag = L_TERM_FUNC_CALL;\n")

	if firstFuncCall {
		f.PrintLabel(depth, "funcCallChain->begin = funcTerm;\n")
		f.PrintLabel(depth, "funcCallChain->end = funcTerm;\n")
	} else {
		f.PrintLabel(depth, fmt.Sprintf("funcCallChain->end->funcCall->next = helper[%d];\n", chainNumber))
		f.PrintLabel(depth, "funcCallChain->end->next = funcTerm;\n")
		f.PrintLabel(depth, "funcTerm->prev = funcCallChain->end;\n")
		f.PrintLabel(depth, "funcCallChain->end = funcTerm;\n")
	}
}

func IsLiteral(termTag syntax.TermTag) bool {

	switch termTag {
	case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:
		return true
	}

	return false
}

func (f *Data) ConstructRelationships(depth int, firstTerm bool, prevChainNumber int, nextChainNumber int) {

	if firstTerm {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->begin = helper[%d];\n", prevChainNumber, nextChainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = helper[%d];\n", prevChainNumber, nextChainNumber))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end->next = helper[%d];\n", prevChainNumber, nextChainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->prev = helper[%d]->chain->end;\n", nextChainNumber, prevChainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = helper[%d];\n", prevChainNumber, nextChainNumber))
	}
}

type ChainInfo struct {
	length   int
	termNum  int
	orderNum int
	syntax.TermTag
	funcName string
}

func concatTerms(a []*syntax.Term, b []*syntax.Term) []*syntax.Term {

	aLen := len(a)
	bLen := len(b)
	newTerms := make([]*syntax.Term, 0, aLen+bLen)

	newTerms = append(newTerms, a...)
	newTerms = append(newTerms, b...)

	return newTerms
}

func (f *Data) ConstructResult(depth int, resultExpr syntax.Expr) {

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .mainChain = 0, .callChain = 0};\n")
	} else {
		chainsCount := f.CalcChainsCount(resultExpr)

		f.PrintLabel(depth, "struct lterm_chain_t* funcCallChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));\n")
		f.PrintLabel(depth, "funcCallChain->begin = 0;\n")
		f.PrintLabel(depth, "funcCallChain->end = 0;\n")
		f.PrintLabel(depth, "struct lterm_t* funcTerm;\n")
		f.PrintLabel(depth, fmt.Sprintf("struct lterm_t** helper = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));\n", chainsCount))
		f.PrintLabel(depth, "int i;\n")
		f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)\n", chainsCount))
		f.PrintLabel(depth, "{\n")
		f.PrintLabel(depth+1, "helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
		f.PrintLabel(depth+1, "helper[i]->chain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));\n")
		f.PrintLabel(depth, "}\n")

		f.PrintLabel(depth, "struct lterm_t* currTerm = 0;\n")

		terms := make([]*syntax.Term, len(resultExpr.Terms))
		copy(terms, resultExpr.Terms)

		isThereEvalTerm := false

		chainInfo := make([]ChainInfo, chainsCount, chainsCount)

		chainIndex := 0
		chainOrder := 0

		chainInfo[chainIndex].length = len(terms)
		chainInfo[chainIndex].termNum = 0
		chainInfo[chainIndex].orderNum = chainOrder

		for chainInfo[0].length > 0 {

			switch terms[0].TermTag {

			case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:

				termsNumber := 0
				fragmentLength := 0
				fragmentOffset := terms[0].Index

				for _, val := range terms {
					if IsLiteral(val.TermTag) && termsNumber < chainInfo[chainIndex].length {

						termsNumber++

						if val.TermTag == syntax.STR {
							fragmentLength += len(val.Value.Str)
						} else {
							fragmentLength++
						}
					} else {
						break
					}
				}

				terms = terms[termsNumber:]
				chainInfo[chainIndex].length -= termsNumber
				chainInfo[chainIndex].termNum++
				firstTerm := chainInfo[chainIndex].termNum == 1
				f.ConstructFragmentLTerm(depth, firstTerm, chainInfo[chainIndex].orderNum, fragmentOffset, fragmentLength)

				break

			case syntax.EXPR, syntax.EVAL:

				chainIndex++
				chainOrder++
				chainInfo[chainIndex].TermTag = terms[0].TermTag
				chainInfo[chainIndex].termNum = 0
				chainInfo[chainIndex].orderNum = chainOrder

				if terms[0].TermTag == syntax.EVAL {
					chainInfo[chainIndex].funcName = terms[0].Exprs[0].Terms[0].Value.Name
				}

				chainInfo[chainIndex].length = len(terms[0].Exprs[0].Terms)
				terms = concatTerms(terms[0].Exprs[0].Terms, terms[1:])

				break

			case syntax.FUNC, syntax.BRACED_EXPR, syntax.BRACKETED_EXPR, syntax.ANGLED_EXPR,
				syntax.VAR, syntax.L, syntax.R:
				//TO DO
				break

			} //switch

			//Обработали последний элемент в подвыражении(Например: термы в скобках, термы внутри скобок вычисления)
			for chainInfo[chainIndex].length == 0 {

				switch chainInfo[chainIndex].TermTag {
				case syntax.EVAL:
					f.ConstructFuncCall(depth, !isThereEvalTerm, chainInfo[chainIndex].funcName, chainInfo[chainIndex].orderNum)
					isThereEvalTerm = true
					break
				case syntax.EXPR:
					f.PrintLabel(depth, fmt.Sprintf("helper[%d]->tag = L_TERM_CHAIN_TAG;\n", chainInfo[chainIndex].orderNum))
					break
				}

				chainIndex--
				if chainIndex < 0 {
					break
				}

				chainInfo[chainIndex].length--
				chainInfo[chainIndex].termNum++
				firstTerm := chainInfo[chainIndex].termNum == 1

				f.ConstructRelationships(depth, firstTerm, chainInfo[chainIndex].orderNum, chainInfo[chainIndex+1].orderNum)
			}
		}

		f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)\n", chainsCount))
		f.PrintLabel(depth, "{\n")
		f.PrintLabel(depth+1, "if(helper[i]->chain->begin)\n")
		f.PrintLabel(depth+1, "{\n")
		f.PrintLabel(depth+2, "helper[i]->chain->begin->prev = 0;\n")
		f.PrintLabel(depth+2, "helper[i]->chain->end->next = 0;\n")
		f.PrintLabel(depth+1, "}\n")
		f.PrintLabel(depth, "}\n")

		f.PrintLabel(depth, "funcCallChain->begin->prev = 0;\n")
		f.PrintLabel(depth, "funcCallChain->end->next = 0;\n")

		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .mainChain = helper[0]->chain, .callChain = funcCallChain};\n")
	}
}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	currEntryPoint := 0

	f.funcHeader(currFunc.FuncName)
	f.FuncDataMemoryAllocation(depth, currFunc)

	f.PrintLabel(depth, "switch (entryPoint)\n") //case begin
	f.PrintLabel(depth, "{\n")                   //case block begin

	for _, s := range currFunc.Sentences {

		f.PrintLabel(depth+1, fmt.Sprintf("case %d: \n", currEntryPoint))
		f.PrintLabel(depth+1, fmt.Sprintf("{\n"))
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
				f.PrintLabel(depth+1, fmt.Sprintf("{\n"))
				break

			case syntax.COMMA: // ','
			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.DCOLON: // '::'
			}
		}
	}

	f.PrintLabel(depth+2, fmt.Sprintf("break;\n")) // last case break
	f.PrintLabel(depth+1, fmt.Sprintf("}\n"))      // last case }
	f.PrintLabel(1, "} // switch block end\n")     //switch block end
	f.FuncDataMemoryFree(1)
	f.PrintLabel(0, fmt.Sprintf("} // %s\n\n", currFunc.FuncName)) // func block end
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.PrintLabel(0, fmt.Sprintf("// file:%s\n\n", currFileData.Name))
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
