package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func (emt *EmitterData) matchingPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	emt.printLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceInfo.index, ctx.sentenceInfo.patternIndex))
	emt.printLabel(depth, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	emt.printLabel(depth, fmt.Sprintf("{"))

	emt.checkAndAssemblyChain(depth+1, ctx)

	emt.checkFragmentLength(depth+1, -1, false, terms)

	if len(terms) > 0 {
		emt.printLabel(depth+1, "else")
		emt.processPattern(depth+2, ctx, terms)
	}

	emt.processPatternFail(depth+1, ctx)

	ctx.addPrevEntryPoint(ctx.entryPointNumerator, ctx.sentenceInfo.actionIndex)
	ctx.entryPointNumerator++
	ctx.sentenceInfo.patternIndex++
}

func (emt *EmitterData) processEmptyPattern(depth int, ctx *emitterContext) {
	emt.printLabel(depth+1, "if (currFrag->length > 0)")
	emt.printFailBlock(depth+1, -1, false)
	emt.printLabel(depth+1, "break;")
}

func (emt *EmitterData) processPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	emt.printLabel(depth, "while (stretchingVarNumber >= 0)")
	emt.printLabel(depth, "{")

	emt.printLabel(depth+1, "//From what stretchable variable start?")
	emt.printLabel(depth+1, "switch (stretchingVarNumber)")
	emt.printLabel(depth+1, "{")

	ctx.patternCtx.entryPoint = 0
	ctx.patternCtx.prevEntryPoint = -1

	emt.printFirstCase(depth, ctx, terms[0])

	emt.matchingTerms(depth+2, false, ctx, terms)

	emt.printLabel(depth+1, "} //pattern switch\n")

	emt.printLabel(depth+1, "if (!stretching)")
	emt.printLabel(depth+1, "{")
	emt.printLabel(depth+2, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	emt.printFailBlock(depth+2, ctx.patternCtx.prevEntryPoint, false)
	emt.printLabel(depth+2, "else")
	emt.printLabel(depth+3, "break; // Success!")
	emt.printLabel(depth+1, "}")

	emt.printLabel(depth, "} // Pattern while\n")
}

func (emt *EmitterData) printFirstCase(depth int, ctx *emitterContext, term *syntax.Term) {

	if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
		if _, ok := ctx.fixedVars[term.Name]; !ok {
			return
		}
	}

	ctx.patternCtx.entryPoint = 1
	emt.printLabel(depth+1, "case 0:")
}

func (emt *EmitterData) matchingTerms(depth int, inBrackets bool, ctx *emitterContext, terms []*syntax.Term) {
	parentMatchingOrder := ctx.isLeftMatching
	termsCount := len(terms)

	if termsCount == 0 {
		return
	}

	ctx.isLeftMatching = !(terms[0].TermTag == syntax.R)

	for _, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			emt.matchingVariable(depth, ctx, &term.Value)
			break
		case syntax.STR:
			emt.matchingStrLiteral(depth, ctx, len(term.Value.Str), term.IndexInLiterals)
			break
		case syntax.COMP:
			emt.matchingCompLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.INT:
			emt.matchingIntLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.EXPR:
			emt.matchingExpr(depth, ctx, term.Exprs[0].Terms)
			break
		case syntax.FLOAT:
			emt.mathcingDoubleLiteral(depth, ctx, term.IndexInLiterals)
			break
		}
	}

	ctx.isLeftMatching = parentMatchingOrder
}

func (emt *EmitterData) getMinLengthForTerms(terms []*syntax.Term) int {
	length := 0

	for _, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			varType := term.Value.VarType
			if varType == tokens.VT_T || varType == tokens.VT_S ||
				varType == tokens.VT_V {
				length += 1
			}
			break
		case syntax.STR:
			length += len(term.Value.Str)
			break
		case syntax.COMP, syntax.INT, syntax.EXPR, syntax.FLOAT:
			length += 1
			break
		}
	}

	return length
}

func (emt *EmitterData) matchingExpr(depth int, ctx *emitterContext, terms []*syntax.Term) {

	ctx.bracketsNumerator++

	tmpBracketsCurrIndex := ctx.bracketsCurrentIndex
	bracketsIndex := ctx.bracketsNumerator
	ctx.bracketsCurrentIndex = bracketsIndex

	emt.printLabel(depth, "//Check ().")
	emt.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, " || _memMngr.vterms[fragmentOffset].tag != V_BRACKETS_TAG")

	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[%d] = fragmentOffset;", bracketsIndex))
	emt.printLabel(depth, "rightBound = RIGHT_BOUND(fragmentOffset);")
	emt.printLabel(depth, "fragmentOffset = VTERM_BRACKETS(fragmentOffset)->offset;")

	emt.checkFragmentLength(depth, ctx.patternCtx.prevEntryPoint, true, terms)

	emt.printLabel(depth, "//Start check in () terms.")
	emt.matchingTerms(depth, true, ctx, terms)

	emt.checkConsumeAllFragment(depth, ctx.patternCtx.prevEntryPoint)

	emt.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", tmpBracketsCurrIndex))
	emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FUNC_CALL->env->bracketsOffset[%d] + 1;", bracketsIndex))

	emt.printLabel(depth, "//End check in () terms.")

	ctx.bracketsCurrentIndex = tmpBracketsCurrIndex
}

func (emt *EmitterData) processPatternFail(depth int, ctx *emitterContext) {

	emt.printLabel(depth, "if (stretchingVarNumber < 0)")
	emt.printLabel(depth, "{")

	prevEntryPoint := ctx.getPrevEntryPoint()
	//First pattern in current sentence
	if ctx.sentenceInfo.patternIndex == 0 || prevEntryPoint == -1 {
		emt.processFailOfFirstPattern(depth+1, ctx)
	} else {
		emt.processFailOfCommonPattern(depth+1, prevEntryPoint)
	}

	emt.printLabel(depth+1, "break;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) checkFragmentLength(depth, prevStertchingVarNumber int, withBreakStatement bool, terms []*syntax.Term) {

	if len(terms) == 0 {
		emt.printLabel(depth, "if (rightBound != fragmentOffset)")
	} else {
		emt.printLabel(depth, fmt.Sprintf("if (rightBound - fragmentOffset < %d)", emt.getMinLengthForTerms(terms)))
	}

	emt.printFailBlock(depth, prevStertchingVarNumber, withBreakStatement)
}

func (emt *EmitterData) processFailOfFirstPattern(depth int, ctx *emitterContext) {
	if ctx.sentenceInfo.isLast {
		emt.printLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		emt.printLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		emt.printLabel(depth, "CURR_FUNC_CALL->entryPoint = -1;")

	} else {
		emt.printLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		emt.printLabel(depth, "stretching = 0;")
		emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.nextSentenceEntryPoint))
		emt.printLabel(depth, "clearCurrFuncEnvData();")
	}
}

func (emt *EmitterData) processFailOfCommonPattern(depth, prevEntryPoint int) {
	emt.printLabel(depth, "//Jump to previouse pattern of same sentence!")
	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", prevEntryPoint))
}

func (emt *EmitterData) checkAndAssemblyChain(depth int, ctx *emitterContext) {
	patternIndex := ctx.sentenceInfo.patternIndex

	emt.printLabel(depth, "if (!stretching)")
	emt.printLabel(depth, "{")

	if ctx.sentenceInfo.actionIndex == 0 {
		if ctx.sentenceInfo.index == 0 {
			emt.printLabel(depth+1, "ASSEMBLY_FIELD(0, CURR_FUNC_CALL->fieldOfView);")
		} else {
			emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[0];",
				patternIndex))
		}
	} else {
		if ctx.needToAssembly() {
			emt.printLabel(depth+1, fmt.Sprintf("ASSEMBLY_FIELD(%d, CURR_FUNC_CALL->env->workFieldOfView);", patternIndex))
		} else {
			emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[%d];",
				patternIndex, patternIndex-1))
		}
	}

	emt.printLabel(depth, "} // !stretching")

	emt.printLabel(depth, fmt.Sprintf("currFrag = VTERM_BRACKETS(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emt.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emt.printLabel(depth+1, "fragmentOffset = currFrag->offset;")
	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[0] = CURR_FUNC_CALL->env->assembled[%d];", patternIndex))
	emt.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = CURR_FUNC_CALL->env->stretchVarsNumber[%d];", ctx.sentenceInfo.patternIndex))
}

func (emt *EmitterData) matchingVariable(depth int, ctx *emitterContext, value *tokens.Value) {

	varInfo, isLocalVar := ctx.sentenceInfo.scope.VarMap[value.Name]
	isFixedVar := true

	if !isLocalVar {
		varInfo = ctx.funcInfo.Env[value.Name]
	} else {
		_, isFixedVar = ctx.fixedVars[value.Name]
	}

	varNumber := varInfo.Number
	emt.printLabel(depth-1, fmt.Sprintf("//Matching %s variable", value.Name))

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				emt.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			emt.matchingFreeTermVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_S:
		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalSymbolVar(depth, ctx, varNumber)
			} else {
				emt.matchingFixedEnvSymbolVar(depth, ctx, varNumber)
			}

		} else {
			emt.matchingFreeSymbolVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_E, tokens.VT_V:

		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				emt.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			emt.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.patternCtx.entryPoint))

			if value.VarType == tokens.VT_E {
				emt.matchingFreeExprVar(depth, ctx, varNumber)
			} else {
				emt.matchingFreeVExprVar(depth, ctx, varNumber)
			}

			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
			ctx.patternCtx.prevEntryPoint = ctx.patternCtx.entryPoint
			ctx.patternCtx.entryPoint++
		}
		break
	}
}

func (emt *EmitterData) printFailBlock(depth, prevStretchVarNumber int, withBreakStatement bool) {

	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "stretching = 1;")
	emt.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	if withBreakStatement {
		emt.printLabel(depth+1, "break;")
	}
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) printOffsetCheck(depth, prevStretchVarNumber int, optionalCond string) {

	emt.printLabel(depth, fmt.Sprintf("if (fragmentOffset >= rightBound%s)", optionalCond))
	emt.printFailBlock(depth, prevStretchVarNumber, true)
}

func (emt *EmitterData) checkConsumeAllFragment(depth, prevStretchVarNumber int) {
	emt.printLabel(depth, "if (fragmentOffset != rightBound)")
	emt.printFailBlock(depth, prevStretchVarNumber, true)
}
