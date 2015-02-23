package emitter

import (
	"fmt"
	//"io"
)

import (
	"bmstu-refal-compiler/syntax"
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

func (f *Data) isLiteral(term *syntax.Term, ctx *emitterContext) bool {

	switch term.TermTag {
	case syntax.STR, syntax.INT, syntax.FLOAT:
		return true

	case syntax.COMP:
		_, yes := f.isFuncName(term.Value.Name, ctx)
		return !yes
	}

	return false
}

func (f *Data) ConstructLiteralsFragment(depth int, ctx *emitterContext, terms []*syntax.Term) []*syntax.Term {
	var term *syntax.Term
	fragmentLength := 0
	fragmentOffset := terms[0].IndexInLiterals
	literalsNumber := 0

	for _, term = range terms {

		if !f.isLiteral(term, ctx) {
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

	f.PrintLabel(depth, "//Adding term to field chain -- Just concat.")
	f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->prev->next = currTerm;", chainNumber))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->prev = helper[%d]->chain->prev;", chainNumber))
	f.PrintLabel(depth, fmt.Sprintf("helper[%d]->chain->prev = currTerm;", chainNumber))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->next = helper[%d]->chain;", chainNumber))
}

func (f *Data) ConstructFuncCall(depth, entryPoint int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	f.PrintLabel(depth, "//Start construction func call.")

	terms = f.ConstructExprInParenthesis(depth, entryPoint, ctx, chainNumber, firstFuncCall, terms)

	f.PrintLabel(depth, "funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth, "funcTerm->tag = L_TERM_FUNC_CALL;")
	f.PrintLabel(depth, "funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));")
	f.PrintLabel(depth, "funcTerm->funcCall->env = (struct env_t*)malloc(sizeof(struct env_t));")

	f.PrintLabel(depth, "funcTerm->funcCall->entryPoint = 0;")
	f.PrintLabel(depth, "funcTerm->funcCall->funcPtr = 0;")
	f.PrintLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")
	f.PrintLabel(depth, "//WARN: Correct free currTerm.")
	f.PrintLabel(depth, "free(currTerm);")

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

func (f *Data) ConstructExprInParenthesis(depth, entryPoint int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	f.PrintLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if f.isLiteral(term, ctx) {
			terms = f.ConstructLiteralsFragment(depth, ctx, terms)
		} else {
			terms = terms[1:]
		}

		// Вызов функции
		if term.TermTag == syntax.EVAL {
			ctx.isFuncCallInConstruct = true
			*chainNumber++
			f.ConstructFuncCall(depth, entryPoint, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			f.ConcatToCallChain(depth, firstFuncCall)
		}

		// Выражение в скобках
		if term.TermTag == syntax.EXPR {
			*chainNumber++
			f.ConstructExprInParenthesis(depth, entryPoint, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
		}

		// Значение переменной
		if term.TermTag == syntax.VAR {
			fixedEntryPoint := ctx.fixedVars[term.Value.Name]
			f.constructVar(depth, fixedEntryPoint, term.Value.Name, ctx)
		}

		//Имя функции. Создаем функциональный vterm.
		if term.TermTag == syntax.COMP {
			//fmt.Printf("Check on functional compound: %s\n", term.Value.Name)
			if funcInfo, yes := f.isFuncName(term.Value.Name, ctx); yes {
				//fmt.Printf("There is functional compound: %s\n", term.Value.Name)
				f.constructFunctionalVTerm(depth, ctx, term.Value.Name, funcInfo)
			}
		}

		//Создание вложенной функции. Создание функционального vterm'a
		if term.TermTag == syntax.FUNC {
			funcInfo := ctx.funcsKeeper.AddFunc(ctx.scopeKeeper, term.Function)
			f.constructFunctionalVTerm(depth, ctx, term.Function.FuncName, funcInfo)
			ctx.nestedNamedFuncs = append(ctx.nestedNamedFuncs, funcInfo)
		}

		f.ConcatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false
	}

	f.PrintLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	f.PrintLabel(depth, fmt.Sprintf("currTerm = helper[%d];", currChainNumber))

	return terms
}

func (f *Data) constructVar(depth, fixedEntryPoint int, varName string, ctx *emitterContext) {

	if scopeVar, ok := ctx.sentenceInfo.scope.VarMap[varName]; ok {
		f.PrintLabel(depth, fmt.Sprintf("currTerm = &env->locals[%d][%d];", fixedEntryPoint, scopeVar.Number))
	} else {
		// Get env var
		needVarInfo, _ := ctx.currFuncInfo.EnvVarMap[varName]
		f.PrintLabel(depth, fmt.Sprintf("currTerm = &env->params[%d];", needVarInfo.Number))
	}
}

func (f *Data) ConstructResult(depth int, ctx *emitterContext, resultExpr syntax.Expr) {
	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {
		chainsCount := calcChainsCount(resultExpr.Terms)

		f.printInitializeConstructVars(depth, chainsCount)

		ctx.isFuncCallInConstruct = false
		firstFuncCall := true
		chainNumber := 0

		f.ConstructExprInParenthesis(depth, ctx.entryPoint-1, ctx, &chainNumber, &firstFuncCall, resultExpr.Terms)

		f.PrintLabel(depth, "fieldOfView = currTerm->chain;")

		if ctx.sentenceInfo.isLastAction() {
			f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else if ctx.isFuncCallInConstruct && ctx.sentenceInfo.isNextMatchingAction() {
			f.PrintLabel(depth, fmt.Sprintf("*entryPoint = %d;", ctx.entryPoint))
			f.PrintLabel(depth, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		}
	}
}

func (f *Data) printInitializeConstructVars(depth, chainsCount int) {
	f.PrintLabel(depth, "//WARN: Correct free funcCallChain.")
	f.PrintLabel(depth, "free(funcCallChain);")
	f.PrintLabel(depth, "funcCallChain = 0;")

	f.PrintLabel(depth, "//WARN: Correct free prev helper.")
	f.PrintLabel(depth, fmt.Sprintf("helper = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));", chainsCount))
	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i)", chainsCount))
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth+1, "helper[i]->tag = L_TERM_CHAIN_TAG;")
	f.PrintLabel(depth+1, "helper[i]->chain = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth+1, "helper[i]->chain->prev = helper[i]->chain;")
	f.PrintLabel(depth+1, "helper[i]->chain->next = helper[i]->chain;")
	f.PrintLabel(depth, "}")
}
