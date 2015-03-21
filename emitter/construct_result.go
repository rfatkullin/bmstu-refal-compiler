package emitter

import (
	"fmt"
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
		if _, _, ok := f.IsFuncName(ctx, term.Value.Name); !ok {
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

func (f *Data) IsFuncName(ctx *emitterContext, name string) (generatedName string, index int, ok bool) {
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
	f.PrintLabel(depth, "currTerm = chAllocateFragmentLTerm(1, &status);")
	f.printCheckGCCondition(depth)
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = %d;", fragmentOffset))
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = %d;", fragmentLength))

	return terms[literalsNumber:]
}

func (f *Data) ConcatToParentChain(depth int, firstTerm bool, chainNumber int) {

	f.PrintLabel(depth, "//Adding term to field chain -- Just concat.")
	f.PrintLabel(depth, fmt.Sprintf("ADD_TO_CHAIN(helper[%d].chain, currTerm);", chainNumber))
}

func (f *Data) ConstructFuncCallTerm(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	f.PrintLabel(depth, "//Start construction func call.")

	terms = f.ConstructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, terms)

	f.PrintLabel(depth, "funcTerm = chAllocateFuncCallLTerm(&status);")
	f.printCheckGCCondition(depth)
	f.PrintLabel(depth, fmt.Sprintf("funcTerm->funcCall->failEntryPoint = %d;", ctx.prevEntryPoint))
	f.PrintLabel(depth, "funcTerm->funcCall->fieldOfView = currTerm->chain;")

	f.PrintLabel(depth, "//Finished construction func call")
	return terms
}

func (f *Data) ConcatToCallChain(depth int, firstFuncCall *bool) {

	if *firstFuncCall {
		f.PrintLabel(depth, "//First call in call chain -- Initialization.")
		f.PrintLabel(depth, "funcCallChain = chAllocateSimpleChainLTerm(&status);")
		f.printCheckGCCondition(depth)
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

func (f *Data) ConstructExprInParenthesis(depth int, ctx *emitterContext, chainNumber *int, firstFuncCall *bool, terms []*syntax.Term) []*syntax.Term {

	firstTermInParenthesis := true
	currChainNumber := *chainNumber

	f.PrintLabel(depth, "//Start construction parenthesis.")

	for 0 < len(terms) {

		term := terms[0]

		if f.isLiteral(term, ctx) {
			terms = f.ConstructLiteralsFragment(depth, ctx, terms)
		} else {
			terms = terms[1:]

			// Вызов функции
			if term.TermTag == syntax.EVAL {
				ctx.isThereFuncCall = true
				*chainNumber++
				f.ConstructFuncCallTerm(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
				f.ConcatToCallChain(depth, firstFuncCall)
			}

			// Выражение в скобках
			if term.TermTag == syntax.EXPR {
				*chainNumber++
				f.ConstructExprInParenthesis(depth, ctx, chainNumber, firstFuncCall, term.Exprs[0].Terms)
			}

			// Значение переменной
			if term.TermTag == syntax.VAR {
				fixedEntryPoint := ctx.fixedVars[term.Value.Name]
				f.constructVar(depth, fixedEntryPoint, term.Value.Name, ctx)
			}

			//Имя функции. Создаем функциональный vterm.
			if term.TermTag == syntax.COMP {
				genFuncName, index, _ := f.IsFuncName(ctx, term.Value.Name)
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

	f.PrintLabel(depth, "//Finished construction parenthesis. Save in currTerm.")
	f.PrintLabel(depth, fmt.Sprintf("currTerm = &helper[%d];", currChainNumber))

	return terms
}

func (f *Data) constructVar(depth, fixedEntryPoint int, varName string, ctx *emitterContext) {

	f.PrintLabel(depth, "currTerm = chAllocateFragmentLTerm(1, &status);")
	f.printCheckGCCondition(depth)

	if scopeVar, ok := ctx.sentenceInfo.scope.VarMap[varName]; ok {
		f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = CURR_FUNC_CALL->env->locals[%d].fragment->offset;", scopeVar.Number))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = CURR_FUNC_CALL->env->locals[%d].fragment->length;", scopeVar.Number))
	} else {
		// Get env var
		needVarInfo, _ := ctx.funcInfo.Env[varName]
		f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = CURR_FUNC_CALL->env->params[%d].fragment->offset;", needVarInfo.Number))
		f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->length = CURR_FUNC_CALL->env->params[%d].fragment->length;", needVarInfo.Number))
	}
}

func (f *Data) ConstructAssembly(depth int, ctx *emitterContext, resultExpr syntax.Expr) {

	f.PrintLabel(depth, "//Start construction assembly action.")

	if len(resultExpr.Terms) == 0 {
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};")
	} else {

		ctx.isThereFuncCall = false
		firstFuncCall := true
		chainNumber := 0

		f.setGCOpenBorder(depth)

		chainsCount := calcChainsCount(resultExpr.Terms)
		f.printInitializeConstructVars(depth+1, chainsCount)

		f.ConstructExprInParenthesis(depth+1, ctx, &chainNumber, &firstFuncCall, resultExpr.Terms)

		//TO CHECK: Always set funcRes.
		if ctx.sentenceInfo.isLastAction() {
			f.PrintLabel(depth+1, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else if ctx.isThereFuncCall && ctx.sentenceInfo.needToEval() {
			f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPoint))
			f.PrintLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};")
		} else {
			f.PrintLabel(depth+1, "workFieldOfView = currTerm->chain;")
		}

		f.setGCCloseBorder(depth)
	}
}

func (f *Data) ConstructFuncCallAction(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.PrintLabel(depth-1, fmt.Sprintf("//Sentence: %d, Call Action", ctx.sentenceInfo.index))
	f.PrintLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPoint))
	f.PrintLabel(depth-1, fmt.Sprintf("{"))

	f.PrintLabel(depth, "//Start construction func call action.")

	firstFuncCall := true
	chainNumber := 1

	f.setGCOpenBorder(depth)

	f.printInitializeConstructVars(depth+1, 2)

	f.ConstructFuncCallTerm(depth+1, ctx, &chainNumber, &firstFuncCall, terms)
	f.ConcatToCallChain(depth+1, &firstFuncCall)
	f.ConcatToParentChain(depth+1, true, 0)

	f.PrintLabel(depth+1, "if (CURR_FUNC_CALL->fieldOfView)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "// There is assembly with func calls -> get this result.")
	f.PrintLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, CURR_FUNC_CALL->fieldOfView);")
	f.PrintLabel(depth+2, "CURR_FUNC_CALL->fieldOfView = 0;")
	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth+1, "else if (workFieldOfView != 0)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "// There is assembly action in previous actions -> get this result.")
	f.PrintLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, workFieldOfView);")
	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth+1, "else")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, fmt.Sprintf("CHECK_ALLOCATION_CONTINUE(workFieldOfView,"+
		" chCopyFieldOfView(CURR_FUNC_CALL->env->fovs[%d], &status), status);", ctx.sentenceInfo.patternIndex-1))
	f.PrintLabel(depth+2, "CONCAT_CHAINS(funcTerm->funcCall->fieldOfView, workFieldOfView);")
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.entryPoint+1))
	f.PrintLabel(depth+1, "return (struct func_result_t){.status = CALL_RESULT, .fieldChain = helper[0].chain, .callChain = funcCallChain};")

	f.setGCCloseBorder(depth)

	f.PrintLabel(depth-1, "} // Pattern or Call Action case end\n")

	ctx.entryPoint++

	f.PrintLabel(depth-1, fmt.Sprintf("case %d:", ctx.entryPoint))
	f.PrintLabel(depth-1, "{")

	f.PrintLabel(depth, "workFieldOfView = CURR_FUNC_CALL->fieldOfView;")
	f.PrintLabel(depth, "CURR_FUNC_CALL->fieldOfView = 0;")
	f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = workFieldOfView, .callChain = 0};")

	ctx.entryPoint++
}

func (f *Data) printInitializeConstructVars(depth, chainsCount int) {
	f.PrintLabel(depth, "allocate_result status = GC_OK;")
	f.PrintLabel(depth, "struct lterm_t* funcCallChain = 0;")
	f.PrintLabel(depth, "struct lterm_t* helper = 0;")
	f.PrintLabel(depth, "struct lterm_t* currTerm = 0;")
	f.PrintLabel(depth, "struct lterm_t* funcTerm = 0;")

	f.PrintLabel(depth, fmt.Sprintf("helper = chAllocateChainKeeperLTerm(UINT64_C(%d), &status);", chainsCount))
	f.printCheckGCCondition(depth)
}

func (f *Data) setGCOpenBorder(depth int) {
	f.PrintLabel(depth, "do { // GC block")
	f.PrintLabel(depth+1, "if(!success)")
	f.PrintLabel(depth+2, "collectGarbage();")
	f.PrintLabel(depth+1, "success = 1;")
}

func (f *Data) setGCCloseBorder(depth int) {
	f.PrintLabel(depth, "} while (!success); // GC block")
}

func (f *Data) printCheckGCCondition(depth int) {
	f.PrintLabel(depth, "if (status == GC_NEED_CLEAN)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "success = 0;")
	f.PrintLabel(depth+1, "continue;")
	f.PrintLabel(depth, "}")
}
