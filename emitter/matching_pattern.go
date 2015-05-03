package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func (f *Data) matchingPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.printLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceInfo.index, ctx.sentenceInfo.patternIndex))
	f.printLabel(depth, fmt.Sprintf("case %d:", ctx.entryPointNumerator))
	f.printLabel(depth, fmt.Sprintf("{"))

	f.checkAndAssemblyChain(depth+1, ctx)

	f.checkFragmentLength(depth+1, -1, false, terms)

	if len(terms) > 0 {
		f.printLabel(depth+1, "else")
		f.processPattern(depth+2, ctx, terms)
	}

	f.processPatternFail(depth+1, ctx)

	ctx.addPrevEntryPoint(ctx.entryPointNumerator, ctx.sentenceInfo.actionIndex)
	ctx.entryPointNumerator++
	ctx.sentenceInfo.patternIndex++
}

func (f *Data) processEmptyPattern(depth int, ctx *emitterContext) {
	f.printLabel(depth+1, "if (currFrag->length > 0)")
	f.printFailBlock(depth+1, -1, false)
	f.printLabel(depth+1, "break;")
}

func (f *Data) processPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.printLabel(depth, "while (stretchingVarNumber >= 0)")
	f.printLabel(depth, "{")

	f.printLabel(depth+1, "//From what stretchable variable start?")
	f.printLabel(depth+1, "switch (stretchingVarNumber)")
	f.printLabel(depth+1, "{")

	ctx.patternCtx.entryPoint = 0
	ctx.patternCtx.prevEntryPoint = -1

	f.printFirstCase(depth, ctx, terms[0])

	f.matchingTerms(depth+2, false, ctx, terms)

	f.printLabel(depth+1, "} //pattern switch\n")

	f.printLabel(depth+1, "if (!stretching)")
	f.printLabel(depth+1, "{")
	f.printLabel(depth+2, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	f.printFailBlock(depth+2, ctx.patternCtx.prevEntryPoint, false)
	f.printLabel(depth+2, "else")
	f.printLabel(depth+3, "break; // Success!")
	f.printLabel(depth+1, "}")

	f.printLabel(depth, "} // Pattern while\n")
}

func (f *Data) printFirstCase(depth int, ctx *emitterContext, term *syntax.Term) {

	if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
		if _, ok := ctx.fixedVars[term.Name]; !ok {
			return
		}
	}

	ctx.patternCtx.entryPoint = 1
	f.printLabel(depth+1, "case 0:")
}

func (f *Data) matchingTerms(depth int, inBrackets bool, ctx *emitterContext, terms []*syntax.Term) {
	parentMatchingOrder := ctx.isLeftMatching
	termsCount := len(terms)

	if termsCount == 0 {
		return
	}

	ctx.isLeftMatching = !(terms[0].TermTag == syntax.R)

	for _, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			f.matchingVariable(depth, ctx, &term.Value)
			break
		case syntax.STR:
			f.matchingStrLiteral(depth, ctx, len(term.Value.Str), term.IndexInLiterals)
			break
		case syntax.COMP:
			f.matchingCompLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.INT:
			f.matchingIntLiteral(depth, ctx, term.IndexInLiterals)
			break
		case syntax.EXPR:
			f.matchingExpr(depth, ctx, term.Exprs[0].Terms)
			break
		case syntax.FLOAT:
			f.mathcingDoubleLiteral(depth, ctx, term.IndexInLiterals)
			break
		}
	}

	ctx.isLeftMatching = parentMatchingOrder
}

func (f *Data) getMinLengthForTerms(terms []*syntax.Term) int {
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

func (f *Data) matchingExpr(depth int, ctx *emitterContext, terms []*syntax.Term) {

	ctx.bracketsNumerator++

	tmpBracketsCurrIndex := ctx.bracketsCurrentIndex
	bracketsIndex := ctx.bracketsNumerator
	ctx.bracketsCurrentIndex = bracketsIndex

	f.printLabel(depth, "//Check ().")
	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, " || _memMngr.vterms[fragmentOffset].tag != V_BRACKETS_TAG")

	f.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[%d] = fragmentOffset;", bracketsIndex))
	f.printLabel(depth, "rightBound = RIGHT_BOUND(fragmentOffset);")
	f.printLabel(depth, "fragmentOffset = VTERM_BRACKETS(fragmentOffset)->offset;")

	f.checkFragmentLength(depth, ctx.patternCtx.prevEntryPoint, true, terms)

	f.printLabel(depth, "//Start check in () terms.")
	f.matchingTerms(depth, true, ctx, terms)

	f.checkConsumeAllFragment(depth, ctx.patternCtx.prevEntryPoint)

	f.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", tmpBracketsCurrIndex))
	f.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FUNC_CALL->env->bracketsOffset[%d] + 1;", bracketsIndex))

	f.printLabel(depth, "//End check in () terms.")

	ctx.bracketsCurrentIndex = tmpBracketsCurrIndex
}

func (f *Data) processPatternFail(depth int, ctx *emitterContext) {

	f.printLabel(depth, "if (stretchingVarNumber < 0)")
	f.printLabel(depth, "{")

	prevEntryPoint := ctx.getPrevEntryPoint()
	//First pattern in current sentence
	if ctx.sentenceInfo.patternIndex == 0 || prevEntryPoint == -1 {
		f.processFailOfFirstPattern(depth+1, ctx)
	} else {
		f.processFailOfCommonPattern(depth+1, prevEntryPoint)
	}

	f.printLabel(depth+1, "break;")
	f.printLabel(depth, "}")
}

func (f *Data) checkFragmentLength(depth, prevStertchingVarNumber int, withBreakStatement bool, terms []*syntax.Term) {

	if len(terms) == 0 {
		f.printLabel(depth, "if (rightBound != fragmentOffset)")
	} else {
		f.printLabel(depth, fmt.Sprintf("if (rightBound - fragmentOffset < %d)", f.getMinLengthForTerms(terms)))
	}

	f.printFailBlock(depth, prevStertchingVarNumber, withBreakStatement)
}

func (f *Data) processFailOfFirstPattern(depth int, ctx *emitterContext) {
	if ctx.sentenceInfo.isLast {
		f.printLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		f.printLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		f.printLabel(depth, "CURR_FUNC_CALL->entryPoint = -1;")

	} else {
		f.printLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		f.printLabel(depth, "stretching = 0;")
		f.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", ctx.nextSentenceEntryPoint))
		f.printLabel(depth, "clearCurrFuncEnvData();")
	}
}

func (f *Data) processFailOfCommonPattern(depth, prevEntryPoint int) {
	f.printLabel(depth, "//Jump to previouse pattern of same sentence!")
	f.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", prevEntryPoint))
}

func (f *Data) checkAndAssemblyChain(depth int, ctx *emitterContext) {
	patternIndex := ctx.sentenceInfo.patternIndex

	f.printLabel(depth, "if (!stretching)")
	f.printLabel(depth, "{")

	if ctx.sentenceInfo.actionIndex == 0 {
		if ctx.sentenceInfo.index == 0 {
			f.printLabel(depth+1, "ASSEMBLY_FIELD(0, CURR_FUNC_CALL->fieldOfView);")
		} else {
			f.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[0];",
				patternIndex))
		}
	} else {
		if ctx.needToAssembly() {
			f.printLabel(depth+1, fmt.Sprintf("ASSEMBLY_FIELD(%d, CURR_FUNC_CALL->env->workFieldOfView);", patternIndex))
		} else {
			f.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[%d];",
				patternIndex, patternIndex-1))
		}
	}

	f.printLabel(depth, "} // !stretching")

	f.printLabel(depth, fmt.Sprintf("currFrag = VTERM_BRACKETS(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	f.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	f.printLabel(depth+1, "fragmentOffset = currFrag->offset;")
	f.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[0] = CURR_FUNC_CALL->env->assembled[%d];", patternIndex))
	f.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = CURR_FUNC_CALL->env->stretchVarsNumber[%d];", ctx.sentenceInfo.patternIndex))
}

func (f *Data) matchingVariable(depth int, ctx *emitterContext, value *tokens.Value) {

	varInfo, isLocalVar := ctx.sentenceInfo.scope.VarMap[value.Name]
	isFixedVar := true

	if !isLocalVar {
		varInfo = ctx.funcInfo.Env[value.Name]
	} else {
		_, isFixedVar = ctx.fixedVars[value.Name]
	}

	varNumber := varInfo.Number
	f.printLabel(depth-1, fmt.Sprintf("//Matching %s variable", value.Name))

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			if isLocalVar {
				f.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				f.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			f.matchingFreeTermVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_S:
		if isFixedVar {
			if isLocalVar {
				f.matchingFixedLocalSymbolVar(depth, ctx, varNumber)
			} else {
				f.matchingFixedEnvSymbolVar(depth, ctx, varNumber)
			}

		} else {
			f.matchingFreeSymbolVar(depth, ctx, varNumber)
			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_E, tokens.VT_V:

		if isFixedVar {
			if isLocalVar {
				f.matchingFixedLocalExprVar(depth, ctx, varNumber)
			} else {
				f.matchingFixedEnvExprVar(depth, ctx, varNumber)
			}
		} else {
			f.printLabel(depth-1, fmt.Sprintf("case %d:", ctx.patternCtx.entryPoint))

			if value.VarType == tokens.VT_E {
				f.matchingFreeExprVar(depth, ctx, varNumber)
			} else {
				f.matchingFreeVExprVar(depth, ctx, varNumber)
			}

			ctx.fixedVars[value.Name] = ctx.sentenceInfo.patternIndex
			ctx.patternCtx.prevEntryPoint = ctx.patternCtx.entryPoint
			ctx.patternCtx.entryPoint++
		}
		break
	}
}

func (f *Data) printFailBlock(depth, prevStretchVarNumber int, withBreakStatement bool) {

	f.printLabel(depth, "{")
	f.printLabel(depth+1, "stretching = 1;")
	f.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	if withBreakStatement {
		f.printLabel(depth+1, "break;")
	}
	f.printLabel(depth, "}")
}

func (f *Data) printOffsetCheck(depth, prevStretchVarNumber int, optionalCond string) {

	f.printLabel(depth, fmt.Sprintf("if (fragmentOffset >= rightBound%s)", optionalCond))
	f.printFailBlock(depth, prevStretchVarNumber, true)
}

func (f *Data) checkConsumeAllFragment(depth, prevStretchVarNumber int) {
	f.printLabel(depth, "if (fragmentOffset != rightBound)")
	f.printFailBlock(depth, prevStretchVarNumber, true)
}
