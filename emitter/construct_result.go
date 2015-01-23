package emitter

import (
	"fmt"
	//"io"
)

import (
	"syntax"
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

	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;\n")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;\n", fragmentOffset))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;\n", fragmentLength))

	return terms[literalsNumber:]
}

func (f *Data) ConcatToParentChain(depth int, firstTerm bool, chainNumber int) {

	if firstTerm {
		//Самый первый терм в цепочке.
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->begin = currTerm;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", chainNumber))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end->next = currTerm;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain->end;\n", chainNumber))
		f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->end = currTerm;\n", chainNumber))
	}
}

func (f *Data) ConstructFuncCall(depth int, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	funcName := terms[0].Value.Name
	currChainNumber := *chainNumber

	terms = f.ConstructExprInParenthesis(depth, chainNumber, firstFuncCall, terms)

	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FUNC_CALL;\n")
	f.PrintLabel(depth, "currTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->funcCall->funcName = %q;\n", funcName))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->funcCall->funcPtr = %s;\n", funcName))
	f.PrintLabel(depth, "currTerm->funcCall->entryPoint = 0;\n")
	f.PrintLabel(depth, "currTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));\n")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->funcCall->fieldOfView->current = helper[%d]->chain;\n", currChainNumber))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->funcCall->inField = helper[%d];\n", currChainNumber))

	return terms
}

func (f *Data) ConcatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		f.PrintLabel(depth, "funcCallChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));\n")
		f.PrintLabel(depth, "currTerm->prev = 0;\n")
		f.PrintLabel(depth, "currTerm->next = 0;\n")
		f.PrintLabel(depth, "funcCallChain->begin = currTerm;\n")
		f.PrintLabel(depth, "funcCallChain->end = currTerm;\n")
		*firstFuncCall = false
	} else {
		f.PrintLabel(depth, "funcCallChain->end->funcCall->next = currTerm;\n")
		f.PrintLabel(depth, "funcCallChain->end->next = currTerm;\n")
		f.PrintLabel(depth, "currTerm->prev = funcCallChain->end;\n")
		f.PrintLabel(depth, "currTerm->next = 0;\n")
		f.PrintLabel(depth, "funcCallChain->end = currTerm;\n")
	}
}

func (f *Data) ConstructExprInParenthesis(depth int, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber
	isEmptyParenthesis := len(terms) == 0

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
			f.ConstructFuncCall(depth, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			f.ConcatToCallChain(depth, firstFuncCall)
		}

		// Выражение в скобках
		if term.TermTag == syntax.EXPR {
			*chainNumber++
			f.ConstructExprInParenthesis(depth, chainNumber, firstFuncCall, term.Exprs[0].Terms)
		}

		f.ConcatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false

		// Остальные случаи
		//case syntax.FUNC, syntax.BRACED_EXPR, syntax.BRACKETED_EXPR, syntax.ANGLED_EXPR,
		//syntax.VAR, syntax.L, syntax.R:
	}

	f.PrintLabel(depth, fmt.Sprintf("currTerm = helper[%d];\n", currChainNumber))
	f.PrintLabel(depth, "currTerm->tag = L_TERM_CHAIN_TAG;\n")

	if isEmptyParenthesis {
		f.PrintLabel(depth, "currTerm->chain->begin = 0;\n")
		f.PrintLabel(depth, "currTerm->chain->end = 0;\n")
	} else {

		f.PrintLabel(depth, "currTerm->chain->begin->prev = 0;\n")
		f.PrintLabel(depth, "currTerm->chain->end->next = 0;\n")
	}

	return terms
}

func (f *Data) ConstructResult(depth int, resultExpr syntax.Expr) {

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};\n")
	} else {
		chainsCount := calcChainsCount(resultExpr.Terms)

		f.PrintLabel(depth, "struct lterm_chain_t* funcCallChain = 0;\n")
		f.PrintLabel(depth, fmt.Sprintf("struct lterm_t** helper = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));\n", chainsCount))
		f.PrintLabel(depth, "int i;\n")
		f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)\n", chainsCount))
		f.PrintLabel(depth, "{\n")
		f.PrintLabel(depth+1, "helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));\n")
		f.PrintLabel(depth+1, "helper[i]->chain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));\n")
		f.PrintLabel(depth, "}\n")

		f.PrintLabel(depth, "struct lterm_t* currTerm = 0;\n")

		firstFuncCall := true
		chainNumber := 0

		f.ConstructExprInParenthesis(depth, &chainNumber, &firstFuncCall, resultExpr.Terms)

		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};\n")
	}
}
