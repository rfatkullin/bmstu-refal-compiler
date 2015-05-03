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

func (f *Data) isLiteral(term *syntax.Term, ctx *emitterContext) bool {

	switch term.TermTag {
	case syntax.STR, syntax.INT, syntax.FLOAT:
		return true

	case syntax.COMP:
		if _, _, ok := f.isFuncName(ctx, term.Value.Name); !ok {
			return true
		} else {
			return false
		}

	}

	return false
}

func (f *Data) genFuncName(index int) string {
	return fmt.Sprintf("func_%d", index)
}

func (f *Data) isFuncName(ctx *emitterContext, name string) (generatedName string, index int, ok bool) {
	level := -1
	index = 0

	if index, level = ctx.sentenceInfo.sentence.FindFunc(name); level != -1 {
		return f.genFuncName(index), index, true
	}

	if _, ok := f.Ast.Builtins[name]; ok {
		return name, -1, true
	}

	if currFunc, ok := f.Ast.GlobMap[name]; ok {
		return f.genFuncName(currFunc.Index), currFunc.Index, true
	}

	return "", -1, false
}

func (f *Data) constructLiteralsFragment(depth int, ctx *emitterContext, terms []*syntax.Term) []*syntax.Term {
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

	f.printLabel(depth, "//Start construction fragment term.")
	f.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")
	f.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;", fragmentOffset))
	f.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;", fragmentLength))

	return terms[literalsNumber:]
}

func (f *Data) ConcatToParentChain(depth int, firstTerm bool, chainNumber int) {

	f.printLabel(depth, "//Adding term to field chain -- Just concat.")
	f.printLabel(depth, fmt.Sprintf("ADD_TO_CHAIN(helper[%d].chain, currTerm);", chainNumber))
}

func (f *Data) constructFuncCallTerm(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	f.printLabel(depth, "//Start construction func call.")

	terms = f.constructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, terms)

	f.printCheckGCCondition(depth, "funcTerm", "chAllocateFuncCallLTerm(&status)")
	f.printLabel(depth, fmt.Sprintf("funcTerm->funcCall->failEntryPoint = %d;", ctx.getPrevEntryPoint()))
	f.printLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")

	f.printLabel(depth, "//Finished construction func call")
	return terms
}

func (f *Data) concatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		f.printLabel(depth, "//First call in call chain -- Initialization.")
		f.printCheckGCCondition(depth, "funcCallChain", "chAllocateSimpleChainLTerm(&status)")
		f.printLabel(depth, "funcCallChain->next = funcTerm;")
		f.printLabel(depth, "funcCallChain->prev = funcTerm;")
		*firstFuncCall = false
	} else {
		f.printLabel(depth, "//Adding call to call chain -- Just concat.")
		f.printLabel(depth, "funcCallChain->prev->funcCall->next = funcTerm;")
		f.printLabel(depth, "funcCallChain->prev = funcTerm;")
	}

	f.printLabel(depth, "currTerm = funcTerm;")
}

func (f *Data) constructExprInParenthesis(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	f.printLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if f.isLiteral(term, ctx) {
			terms = f.constructLiteralsFragment(depth, ctx, terms)
		} else {
			terms = terms[1:]

			// Вызов функции
			if term.TermTag == syntax.EVAL {
				ctx.isThereFuncCall = true
				*chainNumber++
				f.constructFuncCallTerm(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
				f.concatToCallChain(depth, firstFuncCall)
			}

			// Выражение в скобках
			if term.TermTag == syntax.EXPR {
				*chainNumber++
				f.constructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			}

			// Значение переменной
			if term.TermTag == syntax.VAR {
				fixedEntryPoint := ctx.fixedVars[term.Value.Name]
				f.constructVar(depth, fixedEntryPoint, term.Value.Name, ctx)
			}

			//Имя функции. Создаем функциональный vterm.
			if term.TermTag == syntax.COMP {
				genFuncName, index, _ := f.isFuncName(ctx, term.Value.Name)
				f.constructFunctionalVTerm(depth, ctx, term, genFuncName, index)
			}

			//Создание вложенной функции. Создание функционального vterm'a
			if term.TermTag == syntax.FUNC {
				f.constructFunctionalVTerm(depth, ctx, term, f.genFuncName(term.Index), term.Index)
			}
		}

		f.ConcatToParentChain(depth, firstTermInParenthesis, currChainNumber)
		firstTermInParenthesis = false
	}

	f.printLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	f.printLabel(depth, fmt.Sprintf("currTerm = &helper[%d];", currChainNumber))

	return terms
}

func (f *Data) constructVar(depth, fixedEntryPoint int, varName string, ctx *emitterContext) {

	f.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")

	if scopeVar, ok := ctx.sentenceInfo.scope.VarMap[varName]; ok {
		f.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;", scopeVar.Number))
		f.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->locals + %d)->length;", scopeVar.Number))
	} else {
		// Get env var
		needVarInfo, _ := ctx.funcInfo.Env[varName]
		f.printLabel(depth, fmt.Sprintf("currTerm->fragment->offset = (CURR_FUNC_CALL->env->params + %d)->offset;", needVarInfo.Number))
		f.printLabel(depth, fmt.Sprintf("currTerm->fragment->length = (CURR_FUNC_CALL->env->params + %d)->length;", needVarInfo.Number))
	}
}

func (f *Data) constructAssembly(depth int, ctx *emitterContext, resultExpr syntax.Expr) {

	f.printLabel(depth, "//Start construction assembly action.")

	if len(resultExpr.Terms) == 0 {
		f.printLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {

		ctx.isThereFuncCall = false
		firstFuncCall := true
		chainNumber := 0

		f.setGCOpenBorder(depth)

		chainsCount := calcChainsCount(resultExpr.Terms)
		f.printInitializeConstructVars(depth+1, chainsCount)

		f.constructExprInParenthesis(depth+1, ctx, &chainNumber, &firstFuncCall, resultExpr.Terms)

		if ctx.isThereFuncCall && ctx.sentenceInfo.needToEval() {
			f.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPointNumerator))
			f.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else {
			f.printLabel(depth+1, "CURR_FUNC_CALL->env->workFieldOfView = currTerm->chain;")
			f.printLabel(depth+1, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		}

		f.setGCCloseBorder(depth)
	}
}

func (f *Data) constructFuncCallAction(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.printLabel(depth-1, fmt.Sprintf("//Sentence: %d, Call Action", ctx.sentenceInfo.index))
	f.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	f.printLabel(depth-1, fmt.Sprintf("{"))

	f.printLabel(depth, "//Start construction func call action.")

	firstFuncCall := true
	chainNumber := 1

	f.setGCOpenBorder(depth)

	f.printInitializeConstructVars(depth+1, 2)

	f.constructFuncCallTerm(depth+1, ctx, &chainNumber, &firstFuncCall, terms)
	f.concatToCallChain(depth+1, &firstFuncCall)
	f.ConcatToParentChain(depth+1, true, 0)

	f.printLabel(depth+1, "if (CURR_FUNC_CALL->env->workFieldOfView)")
	f.printLabel(depth+1, "{")
	f.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->env->workFieldOfView);")
	f.printLabel(depth+2, "CURR_FUNC_CALL->env->workFieldOfView = 0;")
	f.printLabel(depth+1, "}")
	f.printLabel(depth+1, "else")
	f.printLabel(depth+1, "{")
	f.printLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->fieldOfView);")
	f.printLabel(depth+2, "CURR_FUNC_CALL->fieldOfView = 0;")
	f.printLabel(depth+1, "}")

	f.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPointNumerator+1))
	f.printLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = helper[0].chain, .callChain = funcCallChain};")

	f.setGCCloseBorder(depth)

	f.printLabel(depth-1, "} // Pattern or Call Action case end\n")

	ctx.entryPointNumerator++

	f.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	f.printLabel(depth-1, "{")

	f.printLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = CURR_FUNC_CALL->env->workFieldOfView, .callChain = 0};")

	ctx.entryPointNumerator++
}

func (f *Data) printInitializeConstructVars(depth, chainsCount int) {
	f.printLabel(depth, "struct lterm_t* funcCallChain = 0;")
	f.printLabel(depth, "struct lterm_t* helper = 0;")
	f.printLabel(depth, "struct lterm_t* currTerm = 0;")
	f.printLabel(depth, "struct lterm_t* funcTerm = 0;")

	f.printCheckGCCondition(depth, "helper", fmt.Sprintf("chAllocateChainKeeperLTerm(UINT64_C(%d), &status)", chainsCount))
}

func (f *Data) setGCOpenBorder(depth int) {
	f.printLabel(depth, "do { // GC block")
	f.printLabel(depth+1, "if(prevStatus == GC_NEED_CLEAN)")
	f.printLabel(depth+2, "PRINT_AND_EXIT(GC_MEMORY_OVERFLOW_MSG);")

	f.printLabel(depth+1, "if(status == GC_NEED_CLEAN)")
	f.printLabel(depth+1, "{")
	f.printLabel(depth+2, "collectGarbage();")
	f.printLabel(depth+2, "prevStatus = GC_NEED_CLEAN;")
	f.printLabel(depth+2, "status = GC_OK;")
	f.printLabel(depth+1, "}")
}

func (f *Data) setGCCloseBorder(depth int) {
	f.printLabel(depth, "} while (status != GC_OK); // GC block")
}

func (f *Data) printCheckGCCondition(depth int, varStr, funcCallStr string) {

	f.printLabel(depth, fmt.Sprintf("CHECK_ALLOCATION_CONTINUE(%s, %s, status);", varStr, funcCallStr))
}
