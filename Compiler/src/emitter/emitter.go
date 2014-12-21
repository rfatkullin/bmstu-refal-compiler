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

	fmt.Fprintf(f, "%sstruct func_result_t funcRes;\n", tabs)
	fmt.Fprintf(f, "%sif (entryPoint == 0)\n", tabs)
	fmt.Fprintf(f, "%s{\n", tabs)
	fmt.Fprintf(f, "%s%senv->locals = (struct l_term*)malloc(%d * sizeof(struct l_term));\n", tabs, tab, 1)
	fmt.Fprintf(f, "%s%sfieldOfView->backups = (struct l_term_chain_t*)malloc(%d * sizeof(struct l_term_chain_t));\n", tabs, tab, 1)
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

func (f *Data) MaxDepth(expr syntax.Expr) int {
	depth := 1
	maxDepth := 1
	exprLen := make([]int, 1, 1)
	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	exprLen[0] = len(terms)

	for len(exprLen) > 0 && exprLen[len(exprLen)-1] > 0 {
		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT, syntax.VAR:
			exprLen[len(exprLen)-1]--
			break

		case syntax.EXPR, syntax.EVAL:
			tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
			tmpTerms = append(tmpTerms, terms...)
			terms = tmpTerms

			exprLen = append(exprLen, len(term.Exprs[0].Terms))

			depth++
			maxDepth = Max(maxDepth, depth)
			break
		}

		for len(exprLen) > 0 && exprLen[len(exprLen)-1] == 0 {
			exprLen = exprLen[0 : len(exprLen)-1]

			if len(exprLen) > 0 {
				exprLen[len(exprLen)-1]--
			}

			depth--
		}
	}

	return maxDepth
}

func (f *Data) ConstructFragmentLTerm(depth int, firstTerm bool, exprIndex int, fragmentOffset int, fragmentLength int) {

	f.PrintLabel(depth, "currTerm = (struct l_term*)malloc(sizeof(struct l_term));\n")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;\n")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));\n")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;\n", fragmentOffset))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;\n", fragmentLength))

	//Самый первый терм в цепочке.
	if firstTerm {
		f.PrintLabel(depth, "currTerm->prev = 0;\n")
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->begin = currTerm;\n", exprIndex))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", exprIndex))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end->next = currTerm;\n", exprIndex))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain->end;\n", exprIndex))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", exprIndex))
	}
}

func (f *Data) ConstructFuncCall(depth int, firstFuncCall bool, funcName string, exprIndex int) {

	f.PrintLabel(depth, fmt.Sprintf("helper[%d]->tag = L_TERM_FUNC_CALL;\n", exprIndex))

	if firstFuncCall {
		f.PrintLabel(depth, fmt.Sprintf("funcCallChain->begin = helper[%d];\n", exprIndex))
		f.PrintLabel(depth, fmt.Sprintf("funcCallChain->end = helper[%d];\n", exprIndex))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("funcCallChain->end->funcCall->next = helper[%d];\n", exprIndex))
	}

	f.PrintLabel(depth, "funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("funcCall->funcName = memMngr.literalTermsHeap[helper[%d]->chain->begin->fragment->offset].str;\n", exprIndex))
	f.PrintLabel(depth, fmt.Sprintf("funcCall->funcPtr = %s;\n", funcName))
	f.PrintLabel(depth, "funcCall->entryPoint = 0;\n")
	f.PrintLabel(depth, "funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("funcCall->fieldOfView->current = helper[%d]->chain;\n", exprIndex))
	f.PrintLabel(depth, fmt.Sprintf("helper[%d]->funcCall = funcCall;\n", exprIndex))
}

func IsLiteral(termTag syntax.TermTag) bool {

	switch termTag {
	case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:
		return true
	}

	return false
}

func (f *Data) ConstructRelationships(depth int, firstTerm bool, prevIndex int, nextIndex int) {

	if firstTerm {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->begin = helper[%d];\n", prevIndex, nextIndex))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = helper[%d];\n", prevIndex, nextIndex))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end->next = helper[%d];\n", prevIndex, nextIndex))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->prev = helper[%d]->chain->end;\n", nextIndex, prevIndex))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = helper[%d];\n", prevIndex, nextIndex))
	}
}

func (f *Data) ConstructResult(depth int, resultExpr syntax.Expr) {

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .mainChain = 0, .callChain = 0};\n")
	} else {

		exprsDepth := f.MaxDepth(resultExpr)

		f.PrintLabel(depth, "struct l_term_chain_t* funcCallChain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));\n")
		f.PrintLabel(depth, "funcCallChain->begin = 0;\n")
		f.PrintLabel(depth, "funcCallChain->end = 0;\n")
		f.PrintLabel(depth, "struct func_call_t* funcCall;\n")
		f.PrintLabel(depth, fmt.Sprintf("struct l_term** helper = (struct l_term**)malloc(%d * sizeof(struct l_term*));\n", exprsDepth))
		f.PrintLabel(depth, "int i;\n")
		f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)\n", exprsDepth))
		f.PrintLabel(depth, "{\n")
		f.PrintLabel(depth+1, "helper[i] = (struct l_term*)malloc(sizeof(struct l_term));\n")
		f.PrintLabel(depth+1, "helper[i]->chain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));\n")
		f.PrintLabel(depth, "}\n")

		f.PrintLabel(depth, "struct l_term* currTerm = 0;\n")

		terms := make([]*syntax.Term, len(resultExpr.Terms))
		copy(terms, resultExpr.Terms)

		isThereEvalTerm := false
		fragmentOffset := 0
		fragmentLength := 0

		exprLen := make([]int, exprsDepth, exprsDepth)
		exprCurrLTermNum := make([]int, exprsDepth, exprsDepth)
		exprType := make([]syntax.TermTag, exprsDepth, exprsDepth)
		funcName := make([]string, exprsDepth, exprsDepth)

		exprLen[0] = len(terms)
		exprCurrLTermNum[0] = 0
		exprIndex := 0

		for exprLen[0] > 0 {

			term := terms[0]
			terms = terms[1:]

			switch term.TermTag {

			case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:

				fragmentLength = 1
				fragmentOffset = term.Index
				for _, val := range terms {
					if IsLiteral(val.TermTag) {
						fragmentLength++
					}
				}

				exprLen[exprIndex] -= fragmentLength
				terms = terms[fragmentLength-1:]
				exprCurrLTermNum[exprIndex]++
				firstTerm := exprCurrLTermNum[exprIndex] == 1
				f.ConstructFragmentLTerm(depth, firstTerm, exprIndex, fragmentOffset, fragmentLength)
				break

			case syntax.EXPR, syntax.EVAL:

				tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
				tmpTerms = append(tmpTerms, terms...)
				terms = tmpTerms

				exprIndex++
				exprType[exprIndex] = term.TermTag
				exprLen[exprIndex] = len(term.Exprs[0].Terms)
				exprCurrLTermNum[exprIndex] = 0

				if term.TermTag == syntax.EVAL {
					funcName[exprIndex] = term.Exprs[0].Terms[0].Value.Name
				}

				break

			case syntax.FUNC, syntax.BRACED_EXPR, syntax.BRACKETED_EXPR, syntax.ANGLED_EXPR,
				syntax.VAR, syntax.L, syntax.R:
				//TO DO
				break

			} //switch

			//Обработали последний элемент в подвыражении(Например: термы в скобках, термы внутри скобок вычисления)
			for exprLen[exprIndex] == 0 {

				switch exprType[exprIndex] {
				case syntax.EVAL:
					f.ConstructFuncCall(depth, !isThereEvalTerm, funcName[exprIndex], exprIndex)
					isThereEvalTerm = true
					break
				case syntax.EXPR:
					f.PrintLabel(depth, fmt.Sprintf("helper[%d]->tag = L_TERM_CHAIN_TAG;\n", exprIndex))
					break
				}

				exprIndex--
				if exprIndex < 0 {
					break
				}

				exprLen[exprIndex]--
				exprCurrLTermNum[exprIndex]++
				firstTerm := exprCurrLTermNum[exprIndex] == 1

				f.ConstructRelationships(depth, firstTerm, exprIndex, exprIndex+1)
			}
		}

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
