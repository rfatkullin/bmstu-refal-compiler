package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emt *EmitterData) isLiteral(term *syntax.Term) bool {

	switch term.TermTag {
	case syntax.STR, syntax.INT, syntax.FLOAT:
		return true

	case syntax.COMP:

		if _, _, ok := emt.isFuncName(term.Value.Name); !ok {
			return true
		} else {
			return false
		}

	}

	return false
}

func (emt *EmitterData) genFuncName(name string, index int) string {

	if inNative, ok := syntax.Builtins[emt.dialect][name]; ok {
		if !inNative {
			return name
		}
	}

	return fmt.Sprintf("func_%d", index)
}

func (emt *EmitterData) isFuncName(name string) (string, *syntax.Function, bool) {

	if index, level := emt.ctx.sentenceInfo.sentence.FindFunc(name); level != -1 {
		return emt.genFuncName(name, index), emt.FuncByNumber[index], true
	}

	if _, ok := syntax.Builtins[emt.dialect][name]; ok {
		return name, nil, true
	}

	if gFunc, ok := emt.ctx.ast.GlobMap[name]; ok {
		return emt.genFuncName(name, gFunc.Index), gFunc, true
	}

	if _, ok := emt.ctx.ast.ExtMap[name]; ok {
		if eFunc, ok := emt.AllGlobals[name]; ok {
			return emt.genFuncName(name, eFunc.Index), eFunc, true
		}
	}

	return "", nil, false
}

func (emt *EmitterData) constructLiteralsFragment(depth int, terms []*syntax.Term) []*syntax.Term {
	var term *syntax.Term
	fragmentLength := 0
	fragmentOffset := terms[0].IndexInLiterals
	literalsNumber := 0

	for _, term = range terms {

		if !emt.isLiteral(term) {
			break
		}

		literalsNumber++

		if term.TermTag == syntax.STR {
			fragmentLength += len(term.Value.Str)
		} else {
			fragmentLength++
		}
	}

	emt.printLabel(depth, "//Start construction fragment term.")
	emt.printLabel(depth, "ALLC_FRAG_LTERM(currTerm)")
	emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;", fragmentOffset))
	emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;", fragmentLength))

	return terms[literalsNumber:]
}

func (emt *EmitterData) concatToParentChain(depth int, firstTerm bool, chainNumber int) {

	emt.printLabel(depth, "//Adding term to field chain -- Just concat.")
	emt.printLabel(depth, fmt.Sprintf("ADD_TO_CHAIN(helper[%d].chain, currTerm);", chainNumber))
}

func (emt *EmitterData) constructFuncCallTerm(depth int, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	emt.printLabel(depth, "//Start construction func call.")

	terms = emt.constructExprInParenthesis(depth, chainNumber, firstFuncCall, terms)

	emt.printLabel(depth, "ALLC_FUNC_LTERM(funcTerm)")
	emt.printLabel(depth, fmt.Sprintf("funcTerm->funcCall->failEntryPoint = %d;", emt.ctx.getPrevEntryPoint()))
	emt.printLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")

	emt.printLabel(depth, "//Finished construction func call")
	return terms
}

func (emt *EmitterData) concatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		emt.printLabel(depth, "//First call in call chain -- Initialization.")
		emt.printLabel(depth, "ALLC_SIMPL_CHAIN(funcCallChain)")
		emt.printLabel(depth, "funcCallChain->next = funcTerm;")
		emt.printLabel(depth, "funcCallChain->prev = funcTerm;")
		*firstFuncCall = false
	} else {
		emt.printLabel(depth, "//Adding call to call chain -- Just concat.")
		emt.printLabel(depth, "funcCallChain->prev->funcCall->next = funcTerm;")
		emt.printLabel(depth, "funcCallChain->prev = funcTerm;")
	}

	emt.printLabel(depth, "currTerm = funcTerm;")
}

func (emt *EmitterData) constructExprInParenthesis(depth int, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	emt.printLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if emt.isLiteral(term) {
			terms = emt.constructLiteralsFragment(depth, terms)
		} else {
			terms = terms[1:]

			// Вызов функции
			if term.TermTag == syntax.EVAL {
				emt.ctx.isThereFuncCall = true
				*chainNumber++
				emt.constructFuncCallTerm(depth, chainNumber, firstFuncCall, term.Exprs[0].Terms)
				emt.concatToCallChain(depth, firstFuncCall)
			}

			// Выражение в скобках
			if term.TermTag == syntax.EXPR {
				*chainNumber++
				emt.constructExprInParenthesis(depth, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			}

			// Значение переменной
			if term.TermTag == syntax.VAR {
				emt.constructVar(depth, term.Value.Name)
			}

			//Имя функции. Создаем функциональный vterm.
			if term.TermTag == syntax.COMP {
				genFuncName, index, _ := emt.isFuncName(term.Value.Name)
				emt.constructFunctionalVTerm(depth, term, genFuncName, index)
			}

			//Создание вложенной функции. Создание функционального vterm'a
			if term.TermTag == syntax.FUNC {
				emt.constructFunctionalVTerm(depth, term, emt.genFuncName(term.FuncName, term.Index), term.Function)
			}
		}

		emt.concatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false
	}

	emt.printLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	emt.printLabel(depth, fmt.Sprintf("currTerm = &helper[%d];", currChainNumber))

	return terms
}

func (emt *EmitterData) constructVar(depth int, varName string) {

	emt.printLabel(depth, "ALLC_FRAG_LTERM(currTerm)")

	if scopeVar, ok := emt.ctx.sentenceInfo.scope.VarMap[varName]; ok {
		emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;", scopeVar.Number))
		emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->locals + %d)->length;", scopeVar.Number))
	} else {
		// Get env var
		needVarInfo, _ := emt.ctx.funcInfo.Env[varName]
		emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->params + %d)->offset;", needVarInfo.Number))
		emt.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->params + %d)->length;", needVarInfo.Number))
	}
}

func (emt *EmitterData) constructAssembly(depth int, resultExpr syntax.Expr) {

	if len(resultExpr.Terms) == 0 {
		emt.printLabel(depth+1, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {

		if emt.checkForFailSymbol(resultExpr.Terms) {
			emt.printLabel(depth+1, "return (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
			return
		}

		emt.ctx.isThereFuncCall = false
		firstFuncCall := true
		chainNumber := 0

		emt.printLabel(depth+1, "struct lterm_t* funcCallChain = 0;")
		emt.printLabel(depth+1, "struct lterm_t* helper = 0;")
		emt.printLabel(depth+1, "struct lterm_t* currTerm = 0;")
		emt.printLabel(depth+1, "struct lterm_t* funcTerm = 0;")

		emt.checkNeedDataSize(depth+1, resultExpr.Terms)
		emt.printLabel(depth+1, fmt.Sprintf("ALLC_CHAIN_KEEPER_LTERM_N(helper, %d);", emt.calcChainsCount(resultExpr.Terms)))

		emt.constructExprInParenthesis(depth+1, &chainNumber, &firstFuncCall, resultExpr.Terms)

		if emt.ctx.isThereFuncCall && !emt.ctx.sentenceInfo.isLastAction() {
			emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", emt.ctx.entryPointNumerator+1))
			emt.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else {
			emt.printLabel(depth+1, "CURR_FUNC_CALL->env->workFieldOfView = currTerm->chain;")
			emt.printLabel(depth+1, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		}
	}
}

func (emt *EmitterData) checkForFailSymbol(terms []*syntax.Term) bool {

	for _, term := range terms {
		if term.TermTag == syntax.EVAL && len(term.Exprs[0].Terms) == 0 {
			return true
		}
	}

	return false
}

func (emt *EmitterData) constructFuncCallAction(depth int, actionOp syntax.ActionOp, terms []*syntax.Term) {

	firstFuncCall := true
	chainNumber := 1

	emt.printInitializeConstructVars(depth+1, 2)

	emt.constructFuncCallTerm(depth+1, &chainNumber, &firstFuncCall, terms)
	emt.concatToCallChain(depth+1, &firstFuncCall)
	emt.concatToParentChain(depth+1, true, 0)

	emt.printLabel(depth+1, "if (CURR_FUNC_CALL->env->workFieldOfView)")
	emt.printLabel(depth+1, "{")
	emt.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->env->workFieldOfView);")
	emt.printLabel(depth+2, "CURR_FUNC_CALL->env->workFieldOfView = 0;")
	emt.printLabel(depth+1, "}")
	emt.printLabel(depth+1, "else")
	emt.printLabel(depth+1, "{")
	emt.printLabel(depth+2, "struct lterm_t* copyFieldOfView;")
	emt.printLabel(depth+2, "copyFieldOfView = chCopySimpleExpr(CURR_FUNC_CALL->fieldOfView, &status);")
	emt.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, copyFieldOfView);")
	emt.printLabel(depth+1, "}")

	if emt.ctx.sentenceInfo.isLastAction() && actionOp == syntax.ARROW {
		emt.printLabel(depth+1, "return (struct func_result_t){.status = OK_RESULT, .fieldChain = helper[0].chain, .callChain = funcCallChain};")
	} else {
		emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", emt.ctx.entryPointNumerator+1))
		emt.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = helper[0].chain, .callChain = funcCallChain};")
	}

	if !(emt.ctx.sentenceInfo.isLastAction() && actionOp == syntax.ARROW) {
		emt.printLabel(depth-1, "} // Pattern or Call Action case end\n")

		emt.ctx.entryPointNumerator++

		emt.printLabel(depth-1, fmt.Sprintf("case %d:", emt.ctx.entryPointNumerator))
		emt.printLabel(depth-1, "{")

		emt.printLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = CURR_FUNC_CALL->env->workFieldOfView, .callChain = 0};")
	}
}

func (emt *EmitterData) printInitializeConstructVars(depth, chainsCount int) {
	emt.printLabel(depth, "struct lterm_t* funcCallChain = 0;")
	emt.printLabel(depth, "struct lterm_t* helper = 0;")
	emt.printLabel(depth, "struct lterm_t* currTerm = 0;")
	emt.printLabel(depth, "struct lterm_t* funcTerm = 0;")

	emt.printLabel(depth, fmt.Sprintf("helper = allocateChainKeeperLTerm(UINT64_C(%d));", chainsCount))
}

func (emt *EmitterData) calcChainsCount(terms []*syntax.Term) int {
	chainsCount := 0

	for len(terms) > 0 {
		term := terms[0]
		terms = terms[1:]

		if term.TermTag == syntax.EXPR || term.TermTag == syntax.EVAL {
			chainsCount += emt.calcChainsCount(term.Exprs[0].Terms)
		}
	}

	return chainsCount + 1
}
