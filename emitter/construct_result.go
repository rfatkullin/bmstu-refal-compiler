package emitter

import (
	"fmt"
	//"io"
)

import (
	"BMSTU-Refal-Compiler/syntax"
)

const (
	tab = "\t"
)

func calcChainsCount(terms []*syntax.Term) int {
	chainsCount := 0

	for len(terms) > 0 {
		term := terms[0]
		terms = terms[1:]

		if term.TermTag == syntax.EXPR || term.TermTag == syntax.EVAL {
			chainsCount += calcChainsCount(term.Exprs[0].Terms)
		}
	}

	return chainsCount + 1
}

func isLiteral(termTag syntax.TermTag) bool {

	switch termTag {
	case syntax.STR, syntax.COMP, syntax.INT, syntax.FLOAT:
		return true
	}

	return false
}

func (f *Data) ConstructLiteralsFragment(depth int, terms []*syntax.Term) []*syntax.Term {
	var term *syntax.Term
	fragmentLength := 0
	fragmentOffset := terms[0].Index
	literalsNumber := 0

	for _, term = range terms {

		if !isLiteral(term.TermTag) {
			break
		}

		literalsNumber++

		if term.TermTag == syntax.STR {
			fragmentLength += len(term.Value.Str)
		} else {
			fragmentLength++
		}
	}

	f.PrintLabel(depth, "//Start construction fragment term.")
	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;", fragmentOffset))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;", fragmentLength))

	return terms[literalsNumber:]
}

func (f *Data) ConcatToParentChain(depth int, firstTerm bool, chainNumber int) {

	if firstTerm {
		f.PrintLabel(depth, "//First term in field chain -- Initialization.")
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->next = currTerm;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->prev = currTerm;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->next = helper[%d]->chain;", chainNumber))
	} else {
		f.PrintLabel(depth, "//Adding term to field chain -- Just concat.")
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->prev->next = currTerm;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain->prev;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->prev = currTerm;", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->next = helper[%d]->chain;", chainNumber))
	}
}

func (f *Data) printFuncCallPointer(depth int, terms []*syntax.Term) (string, []*syntax.Term) {
	funcName := ""
	term := terms[0]

	switch term.TermTag {

	case syntax.COMP:
		terms = terms[1:]
		funcName = term.Value.Name
		break

	case syntax.VAR:
		//TO DO:
		break

	default:
		f.printFailBlock(depth, -1, true)
		break
	}

	return funcName, terms
}

func (f *Data) ConstructFuncCall(depth, entryPoint int, sentenceScope *syntax.Scope, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	funcName, terms := f.printFuncCallPointer(depth, terms)

	terms = f.ConstructExprInParenthesis(depth, entryPoint, sentenceScope, chainNumber, firstFuncCall, terms)

	f.PrintLabel(depth, "//Start construction func call.")
	f.PrintLabel(depth, "funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth, "funcTerm->tag = L_TERM_FUNC_CALL;")
	f.PrintLabel(depth, "funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));")
	f.PrintLabel(depth, "funcTerm->funcCall->env = (struct env_t*)malloc(sizeof(struct env_t));")
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->funcName = %q;", funcName))
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->funcPtr = %s;", funcName))
	f.PrintLabel(depth, "funcTerm->funcCall->entryPoint = 0;")
	f.PrintLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")
	f.PrintLabel(depth, "//WARN: Begin")
	f.PrintLabel(depth, "free(currTerm->chain);")
	f.PrintLabel(depth, "//WARN: End")

	f.PrintLabel(depth, "//Finished construction func call")
	return terms
}

func (f *Data) ConcatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		f.PrintLabel(depth, "//First call in call chain -- Initialization.")
		f.PrintLabel(depth, "funcCallChain = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
		f.PrintLabel(depth, "funcCallChain->next = funcTerm;")
		f.PrintLabel(depth, "funcCallChain->prev = funcTerm;")
		*firstFuncCall = false
	} else {
		f.PrintLabel(depth, "//Adding call to call chain -- Just concat.")
		f.PrintLabel(depth, "funcCallChain->prev->funcCall->next = funcTerm;")
		f.PrintLabel(depth, "funcCallChain->prev = funcTerm;")
	}

	f.PrintLabel(depth, "currTerm = funcTerm;")
}

func (f *Data) ConstructExprInParenthesis(depth, entryPoint int, sentenceScope *syntax.Scope, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	f.PrintLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if isLiteral(term.TermTag) {
			terms = f.ConstructLiteralsFragment(depth, terms)
		} else {
			terms = terms[1:]
		}

		// Вызов функции
		if term.TermTag == syntax.EVAL {
			*chainNumber++
			f.ConstructFuncCall(depth, entryPoint, sentenceScope, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			f.ConcatToCallChain(depth, firstFuncCall)
		}

		// Выражение в скобках
		if term.TermTag == syntax.EXPR {
			*chainNumber++
			f.ConstructExprInParenthesis(depth, entryPoint, sentenceScope, chainNumber, firstFuncCall, term.Exprs[0].Terms)
		}

		// Значение переменной
		if term.TermTag == syntax.VAR {
			f.PrintLabel(depth, fmt.Sprintf("currTerm = &env->locals[%d][%d];", entryPoint, sentenceScope.VarMap[term.Value.Name].Number))
		}

		f.ConcatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false

		// Остальные случаи
		//case syntax.FUNC, syntax.BRACED_EXPR, syntax.BRACKETED_EXPR, syntax.ANGLED_EXPR,
		//syntax.L, syntax.R:
	}

	f.PrintLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	f.PrintLabel(depth, fmt.Sprintf("currTerm = helper[%d];", currChainNumber))

	//f.PrintLabel(depth, "currTerm->chain->next = currTerm->chain;")
	//f.PrintLabel(depth, "currTerm->chain->prev = currTerm->chain;")

	return terms
}

func (f *Data) ConstructResult(depth, entryPoint int, sentenceScope *syntax.Scope, resultExpr syntax.Expr) {

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {
		chainsCount := calcChainsCount(resultExpr.Terms)

		f.PrintLabel(depth, "struct lterm_t* funcCallChain = 0;")
		f.PrintLabel(depth, fmt.Sprintf("struct lterm_t** helper = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));", chainsCount))
		f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)", chainsCount))
		f.PrintLabel(depth, "{")
		f.PrintLabel(depth+1, "helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
		f.PrintLabel(depth+1, "helper[i]->tag = L_TERM_CHAIN_TAG;")
		f.PrintLabel(depth+1, "helper[i]->chain = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
		f.PrintLabel(depth+1, "helper[i]->chain->prev = 0;")
		f.PrintLabel(depth+1, "helper[i]->chain->next = 0;")
		f.PrintLabel(depth, "}")

		f.PrintLabel(depth, "struct lterm_t* currTerm = 0;")
		f.PrintLabel(depth, "struct lterm_t* funcTerm = 0;")

		firstFuncCall := true
		chainNumber := 0

		f.ConstructExprInParenthesis(depth, entryPoint, sentenceScope, &chainNumber, &firstFuncCall, resultExpr.Terms)

		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
	}
}
