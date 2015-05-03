package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
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

func (emitter *EmitterData) isLiteral(term *syntax.Term, ctx *emitterContext) bool {

	switch term.TermTag {
	case syntax.STR, syntax.INT, syntax.FLOAT:
		return true

	case syntax.COMP:
		if _, _, ok := emitter.isFuncName(ctx, term.Value.Name); !ok {
			return true
		} else {
			return false
		}

	}

	return false
}

func (emitter *EmitterData) genFuncName(index int) string {
	return fmt.Sprintf("func_%d", index)
}

func (emitter *EmitterData) isFuncName(ctx *emitterContext, name string) (generatedName string, index int, ok bool) {
	level := -1
	index = 0

	if index, level = ctx.sentenceInfo.sentence.FindFunc(name); level != -1 {
		return emitter.genFuncName(index), index, true
	}

	if _, ok := emitter.Ast.Builtins[name]; ok {
		return name, -1, true
	}

	if currFunc, ok := emitter.Ast.GlobMap[name]; ok {
		return emitter.genFuncName(currFunc.Index), currFunc.Index, true
	}

	return "", -1, false
}

func (emitter *EmitterData) constructLiteralsFragment(depth int, ctx *emitterContext, terms []*syntax.Term) []*syntax.Term {
	var term *syntax.Term
	fragmentLength := 0
	fragmentOffset := terms[0].IndexInLiterals
	literalsNumber := 0

	for _, term = range terms {

		if !emitter.isLiteral(term, ctx) {
			break
		}

		literalsNumber++

		if term.TermTag == syntax.STR {
			fragmentLength += len(term.Value.Str)
		} else {
			fragmentLength++
		}
	}

	emitter.printLabel(depth, "//Start construction fragment term.")
	emitter.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")
	emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;", fragmentOffset))
	emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;", fragmentLength))

	return terms[literalsNumber:]
}

func (emitter *EmitterData) ConcatToParentChain(depth int, firstTerm bool, chainNumber int) {

	emitter.printLabel(depth, "//Adding term to field chain -- Just concat.")
	emitter.printLabel(depth, fmt.Sprintf("ADD_TO_CHAIN(helper[%d].chain, currTerm);", chainNumber))
}

func (emitter *EmitterData) constructFuncCallTerm(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	emitter.printLabel(depth, "//Start construction func call.")

	terms = emitter.constructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, terms)

	emitter.printCheckGCCondition(depth, "funcTerm", "chAllocateFuncCallLTerm(&status)")
	emitter.printLabel(depth, fmt.Sprintf("funcTerm->funcCall->failEntryPoint = %d;", ctx.getPrevEntryPoint()))
	emitter.printLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")

	emitter.printLabel(depth, "//Finished construction func call")
	return terms
}

func (emitter *EmitterData) concatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		emitter.printLabel(depth, "//First call in call chain -- Initialization.")
		emitter.printCheckGCCondition(depth, "funcCallChain", "chAllocateSimpleChainLTerm(&status)")
		emitter.printLabel(depth, "funcCallChain->next = funcTerm;")
		emitter.printLabel(depth, "funcCallChain->prev = funcTerm;")
		*firstFuncCall = false
	} else {
		emitter.printLabel(depth, "//Adding call to call chain -- Just concat.")
		emitter.printLabel(depth, "funcCallChain->prev->funcCall->next = funcTerm;")
		emitter.printLabel(depth, "funcCallChain->prev = funcTerm;")
	}

	emitter.printLabel(depth, "currTerm = funcTerm;")
}

func (emitter *EmitterData) constructExprInParenthesis(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	emitter.printLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if emitter.isLiteral(term, ctx) {
			terms = emitter.constructLiteralsFragment(depth, ctx, terms)
		} else {
			terms = terms[1:]

			// Вызов функции
			if term.TermTag == syntax.EVAL {
				ctx.isThereFuncCall = true
				*chainNumber++
				emitter.constructFuncCallTerm(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
				emitter.concatToCallChain(depth, firstFuncCall)
			}

			// Выражение в скобках
			if term.TermTag == syntax.EXPR {
				*chainNumber++
				emitter.constructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			}

			// Значение переменной
			if term.TermTag == syntax.VAR {
				fixedEntryPoint := ctx.fixedVars[term.Value.Name]
				emitter.constructVar(depth, fixedEntryPoint, term.Value.Name, ctx)
			}

			//Имя функции. Создаем функциональный vterm.
			if term.TermTag == syntax.COMP {
				genFuncName, index, _ := emitter.isFuncName(ctx, term.Value.Name)
				emitter.constructFunctionalVTerm(depth, ctx, term, genFuncName, index)
			}

			//Создание вложенной функции. Создание функционального vterm'a
			if term.TermTag == syntax.FUNC {
				emitter.constructFunctionalVTerm(depth, ctx, term, emitter.genFuncName(term.Index), term.Index)
			}
		}

		emitter.ConcatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false
	}

	emitter.printLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	emitter.printLabel(depth, fmt.Sprintf("currTerm = &helper[%d];", currChainNumber))

	return terms
}

func (emitter *EmitterData) constructVar(depth, fixedEntryPoint int, varName string, ctx *emitterContext) {

	emitter.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")

	if scopeVar, ok := ctx.sentenceInfo.scope.VarMap[varName]; ok {
		emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;", scopeVar.Number))
		emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->locals + %d)->length;", scopeVar.Number))
	} else {
		// Get env var
		needVarInfo, _ := ctx.funcInfo.Env[varName]
		emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->params + %d)->offset;", needVarInfo.Number))
		emitter.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->params + %d)->length;", needVarInfo.Number))
	}
}

func (emitter *EmitterData) constructAssembly(depth int, ctx *emitterContext, resultExpr syntax.Expr) {

	emitter.printLabel(depth, "//Start construction assembly action.")

	if len(resultExpr.Terms) == 0 {
		emitter.printLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {

		ctx.isThereFuncCall = false
		firstFuncCall := true
		chainNumber := 0

		emitter.setGCOpenBorder(depth)

		chainsCount := calcChainsCount(resultExpr.Terms)
		emitter.printInitializeConstructVars(depth+1, chainsCount)

		emitter.constructExprInParenthesis(depth+1, ctx, &chainNumber, &firstFuncCall, resultExpr.Terms)

		if ctx.isThereFuncCall && ctx.sentenceInfo.needToEval() {
			emitter.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPointNumerator))
			emitter.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else {
			emitter.printLabel(depth+1, "CURR_FUNC_CALL->env->workFieldOfView = currTerm->chain;")
			emitter.printLabel(depth+1, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		}

		emitter.setGCCloseBorder(depth)
	}
}

func (emitter *EmitterData) constructFuncCallAction(depth int, ctx *emitterContext, terms []*syntax.Term) {

	emitter.printLabel(depth-1, fmt.Sprintf("//Sentence: %d, Call Action", ctx.sentenceInfo.index))
	emitter.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	emitter.printLabel(depth-1, fmt.Sprintf("{"))

	emitter.printLabel(depth, "//Start construction func call action.")

	firstFuncCall := true
	chainNumber := 1

	emitter.setGCOpenBorder(depth)

	emitter.printInitializeConstructVars(depth+1, 2)

	emitter.constructFuncCallTerm(depth+1, ctx, &chainNumber, &firstFuncCall, terms)
	emitter.concatToCallChain(depth+1, &firstFuncCall)
	emitter.ConcatToParentChain(depth+1, true, 0)

	emitter.printLabel(depth+1, "if (CURR_FUNC_CALL->env->workFieldOfView)")
	emitter.printLabel(depth+1, "{")
	emitter.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->env->workFieldOfView);")
	emitter.printLabel(depth+2, "CURR_FUNC_CALL->env->workFieldOfView = 0;")
	emitter.printLabel(depth+1, "}")
	emitter.printLabel(depth+1, "else")
	emitter.printLabel(depth+1, "{")
	emitter.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->fieldOfView);")
	emitter.printLabel(depth+2, "CURR_FUNC_CALL->fieldOfView = 0;")
	emitter.printLabel(depth+1, "}")

	emitter.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPointNumerator+1))
	emitter.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = helper[0].chain, .callChain = funcCallChain};")

	emitter.setGCCloseBorder(depth)

	emitter.printLabel(depth-1, "} // Pattern or Call Action case end\n")

	ctx.entryPointNumerator++

	emitter.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	emitter.printLabel(depth-1, "{")

	emitter.printLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = CURR_FUNC_CALL->env->workFieldOfView, .callChain = 0};")

	ctx.entryPointNumerator++
}

func (emitter *EmitterData) printInitializeConstructVars(depth, chainsCount int) {
	emitter.printLabel(depth, "struct lterm_t* funcCallChain = 0;")
	emitter.printLabel(depth, "struct lterm_t* helper = 0;")
	emitter.printLabel(depth, "struct lterm_t* currTerm = 0;")
	emitter.printLabel(depth, "struct lterm_t* funcTerm = 0;")

	emitter.printCheckGCCondition(depth, "helper", fmt.Sprintf("chAllocateChainKeeperLTerm(UINT64_C(%d), &status)", chainsCount))
}

func (emitter *EmitterData) setGCOpenBorder(depth int) {
	emitter.printLabel(depth, "do { // GC block")
	emitter.printLabel(depth+1, "if(prevStatus == GC_NEED_CLEAN)")
	emitter.printLabel(depth+2, "PRINT_AND_EXIT(GC_MEMORY_OVERFLOW_MSG);")

	emitter.printLabel(depth+1, "if(status == GC_NEED_CLEAN)")
	emitter.printLabel(depth+1, "{")
	emitter.printLabel(depth+2, "collectGarbage();")
	emitter.printLabel(depth+2, "prevStatus = GC_NEED_CLEAN;")
	emitter.printLabel(depth+2, "status = GC_OK;")
	emitter.printLabel(depth+1, "}")
}

func (emitter *EmitterData) setGCCloseBorder(depth int) {
	emitter.printLabel(depth, "} while (status != GC_OK); // GC block")
}

func (emitter *EmitterData) printCheckGCCondition(depth int, varStr, funcCallStr string) {

	emitter.printLabel(depth, fmt.Sprintf("CHECK_ALLOCATION_CONTINUE(%s, %s, status);", varStr, funcCallStr))
}
