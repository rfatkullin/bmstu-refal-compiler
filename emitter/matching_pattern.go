package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func (emitter *EmitterData) matchingPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	emitter.printLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceInfo.index, ctx.sentenceInfo.patternIndex))
	emitter.printLabel(depth, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	emitter.printLabel(depth, fmt.Sprintf("{"))

	emitter.checkAndAssemblyChain(depth+1, ctx)

	emitter.checkFragmentLength(depth+1, -1, false, terms)

	if len(terms) > 0 {
		emitter.printLabel(depth+1, "else")
		emitter.processPattern(depth+2, ctx, terms)
	}

	emitter.processPatternFail(depth+1, ctx)

	ctx.addPrevEntryPoint(ctx.entryPointNumerator, ctx.sentenceInfo.actionIndex)
	ctx.entryPointNumerator++
	ctx.sentenceInfo.patternIndex++
}

func (emitter *EmitterData) processEmptyPattern(depth int, ctx *emitterContext) {
	emitter.printLabel(depth+1, "if (currFrag->length > 0)")
	emitter.printFailBlock(depth+1, -1, false)
	emitter.printLabel(depth+1, "break;")
}

func (emitter *EmitterData) processPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	emitter.printLabel(depth, "while (stretchingVarNumber >= 0)")
	emitter.printLabel(depth, "{")

	emitter.printLabel(depth+1, "//From what stretchable variable start?")
	emitter.printLabel(depth+1, "switch (stretchingVarNumber)")
	emitter.printLabel(depth+1, "{")

	ctx.patternCtx.entryPoint = 0
	ctx.patternCtx.prevEntryPoint = -1

	emitter.printFirstCase(depth, ctx, terms[0])

	emitter.matchingTerms(depth+2, false, ctx, terms)

	emitter.printLabel(depth+1, "} //pattern switch\n")

	emitter.printLabel(depth+1, "if (!stretching)")
	emitter.printLabel(depth+1, "{")
	emitter.printLabel(depth+2, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	emitter.printFailBlock(depth+2, ctx.patternCtx.prevEntryPoint, false)
	emitter.printLabel(depth+2, "else")
	emitter.printLabel(depth+3, "break; // Success!")
	emitter.printLabel(depth+1, "}")

	emitter.printLabel(depth, "} // Pattern while\n")
}

func (emitter *EmitterData) printFirstCase(depth int, ctx *emitterContext, term *syntax.Term) {

	if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
		if _, ok := ctx.fixedVars[term.Name]; !ok {
			return
		}
	}

	ctx.patternCtx.entryPoint = 1
	emitter.printLabel(depth+1, "case 0:")
}

func (emitter *EmitterData) matchingTerms(depth int, inBrackets bool, ctx *emitterContext, terms []*syntax.Term) {
	parentMatchingOrder := ctx.isLeftMatching
	termsCount := len(terms)

	if termsCount == 0 {
		return
	}

	ctx.isLeftMatching = !(terms[0].TermTag == syntax.R)

	for _, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			emitter.matchingVariable(depth, ctx, &term.Value)
			break
		case syntax.STR:
			emitter.matchingStrLiteral(depth, ctx, len(term.Value.Str), term.IndexInLiterals)
			break
		case syntax.COMP:
			emitter.matchingCompLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.INT:
			emitter.matchingIntLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.EXPR:
			emitter.matchingExpr(depth, ctx, term.Exprs[0].Terms)
			break
		case syntax.FLOAT:
			emitter.mathcingDoubleLiteral(depth, ctx, term.IndexInLiterals)
			break
		}
	}

	ctx.isLeftMatching = parentMatchingOrder
}

func (emitter *EmitterData) getMinLengthForTerms(terms []*syntax.Term) int {
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

func (emitter *EmitterData) matchingExpr(depth int, ctx *emitterContext, terms []*syntax.Term) {

	ctx.bracketsNumerator++

	tmpBracketsCurrIndex := ctx.bracketsCurrentIndex
	bracketsIndex := ctx.bracketsNumerator
	ctx.bracketsCurrentIndex = bracketsIndex

	emitter.printLabel(depth, "//Check ().")
	emitter.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, " || _memMngr.vterms[fragmentOffset].tag != V_BRACKETS_TAG")

	emitter.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[%d] = fragmentOffset;", bracketsIndex))
	emitter.printLabel(depth, "rightBound = RIGHT_BOUND(fragmentOffset);")
	emitter.printLabel(depth, "fragmentOffset = VTERM_BRACKETS(fragmentOffset)->offset;")

	emitter.checkFragmentLength(depth, ctx.patternCtx.prevEntryPoint, true, terms)

	emitter.printLabel(depth, "//Start check in () terms.")
	emitter.matchingTerms(depth, true, ctx, terms)

	emitter.checkConsumeAllFragment(depth, ctx.patternCtx.prevEntryPoint)

	emitter.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", tmpBracketsCurrIndex))
	emitter.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FUNC_CALL->env->bracketsOffset[%d] + 1;", bracketsIndex))

	emitter.printLabel(depth, "//End check in () terms.")

	ctx.bracketsCurrentIndex = tmpBracketsCurrIndex
}

func (emitter *EmitterData) processPatternFail(depth int, ctx *emitterContext) {

	emitter.printLabel(depth, "if (stretchingVarNumber < 0)")
	emitter.printLabel(depth, "{")

	prevEntryPoint := ctx.getPrevEntryPoint()
	//First pattern in current sentence
	if ctx.sentenceInfo.patternIndex == 0 || prevEntryPoint == -1 {
		emitter.processFailOfFirstPattern(depth+1, ctx)
	} else {
		emitter.processFailOfCommonPattern(depth+1, prevEntryPoint)
	}

	emitter.printLabel(depth+1, "break;")
	emitter.printLabel(depth, "}")
}

func (emitter *EmitterData) checkFragmentLength(depth, prevStertchingVarNumber int, withBreakStatement bool, terms []*syntax.Term) {

	if len(terms) == 0 {
		emitter.printLabel(depth, "if (rightBound != fragmentOffset)")
	} else {
		emitter.printLabel(depth, fmt.Sprintf("if (rightBound - fragmentOffset < %d)", emitter.getMinLengthForTerms(terms)))
	}

	emitter.printFailBlock(depth, prevStertchingVarNumber, withBreakStatement)
}

func (emitter *EmitterData) processFailOfFirstPattern(depth int, ctx *emitterContext) {
	if ctx.sentenceInfo.isLast {
		emitter.printLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		emitter.printLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		emitter.printLabel(depth, "CURR_FUNC_CALL->entryPoint = -1;")

	} else {
		emitter.printLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		emitter.printLabel(depth, "stretching = 0;")
		emitter.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.nextSentenceEntryPoint))
		emitter.printLabel(depth, "clearCurrFuncEnvData();")
	}
}

func (emitter *EmitterData) processFailOfCommonPattern(depth, prevEntryPoint int) {
	emitter.printLabel(depth, "//Jump to previouse pattern of same sentence!")
	emitter.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", prevEntryPoint))
}

func (emitter *EmitterData) checkAndAssemblyChain(depth int, ctx *emitterContext) {
	patternIndex := ctx.sentenceInfo.patternIndex

	emitter.printLabel(depth, "if (!stretching)")
	emitter.printLabel(depth, "{")

	if ctx.sentenceInfo.actionIndex == 0 {
		if ctx.sentenceInfo.index == 0 {
			emitter.printLabel(depth+1, "ASSEMBLY_FIELD(0, CURR_FUNC_CALL->fieldOfView);")
		} else {
			emitter.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[0];",
				patternIndex))
		}
	} else {
		if ctx.needToAssembly() {
			emitter.printLabel(depth+1, fmt.Sprintf("ASSEMBLY_FIELD(%d, CURR_FUNC_CALL->env->workFieldOfView);", patternIndex))
		} else {
			emitter.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[%d];",
				patternIndex, patternIndex-1))
		}
	}

	emitter.printLabel(depth, "} // !stretching")

	emitter.printLabel(depth, fmt.Sprintf("currFrag = VTERM_BRACKETS(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emitter.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emitter.printLabel(depth+1, "fragmentOffset = currFrag->offset;")
	emitter.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[0] = CURR_FUNC_CALL->env->assembled[%d];", patternIndex))
	emitter.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = CURR_FUNC_CALL->env->stretchVarsNumber[%d];", ctx.sentenceInfo.patternIndex))
}

func (emitter *EmitterData) matchingVariable(depth int, ctx *emitterContext, value *tokens.Value) {

	varInfo, isLocalVar := ctx.sentenceInfo.scope.VarMap[value.Name]
	isFixedVar := true

	if !isLocalVar {
		varInfo = ctx.funcInfo.Env[value.Name]
	} else {
		_, isFixedVar = ctx.fixedVars[value.Name]
	}

	varNumber := varInfo.Number
	emitter.printLabel(depth-1, fmt.Sprintf("//Matching %s variable", value.Name))

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			if isLocalVar {
				emitter.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				emitter.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			emitter.matchingFreeTermVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_S:
		if isFixedVar {
			if isLocalVar {
				emitter.matchingFixedLocalSymbolVar(depth, ctx, varNumber)
			} else {
				emitter.matchingFixedEnvSymbolVar(depth, ctx, varNumber)
			}

		} else {
			emitter.matchingFreeSymbolVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_E, tokens.VT_V:

		if isFixedVar {
			if isLocalVar {
				emitter.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				emitter.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			emitter.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.patternCtx.entryPoint))

			if value.VarType == tokens.VT_E {
				emitter.matchingFreeExprVar(depth, ctx, varNumber)
			} else {
				emitter.matchingFreeVExprVar(depth, ctx, varNumber)
			}

			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
			ctx.patternCtx.prevEntryPoint = ctx.patternCtx.entryPoint
			ctx.patternCtx.entryPoint++
		}
		break
	}
}

func (emitter *EmitterData) printFailBlock(depth, prevStretchVarNumber int, withBreakStatement bool) {

	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, "stretching = 1;")
	emitter.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	if withBreakStatement {
		emitter.printLabel(depth+1, "break;")
	}
	emitter.printLabel(depth, "}")
}

func (emitter *EmitterData) printOffsetCheck(depth, prevStretchVarNumber int, optionalCond string) {

	emitter.printLabel(depth, fmt.Sprintf("if (fragmentOffset >= rightBound%s)", optionalCond))
	emitter.printFailBlock(depth, prevStretchVarNumber, true)
}

func (emitter *EmitterData) checkConsumeAllFragment(depth, prevStretchVarNumber int) {
	emitter.printLabel(depth, "if (fragmentOffset != rightBound)")
	emitter.printFailBlock(depth, prevStretchVarNumber, true)
}
