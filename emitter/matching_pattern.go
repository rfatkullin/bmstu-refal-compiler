package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func (f *Data) matchingPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.PrintLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceInfo.index, ctx.sentenceInfo.patternIndex))
	f.PrintLabel(depth, fmt.Sprintf("case %d:", ctx.entryPoint))
	f.PrintLabel(depth, fmt.Sprintf("{"))

	f.checkAndAssemblyChain(depth+1, ctx)

	f.PrintLabel(depth+1, "fragmentOffset = currFrag->offset;")
	f.PrintLabel(depth+1, fmt.Sprintf("stretchingVarNumber = _currFuncCall->env->stretchVarsNumber[%d];", ctx.sentenceInfo.patternIndex))

	f.PrintLabel(depth+1, "while (stretchingVarNumber >= 0)")
	f.PrintLabel(depth+1, "{")

	if len(terms) == 0 {
		f.processEmptyPattern(depth+1, ctx)
	} else {
		f.processPattern(depth+1, ctx, terms)
	}

	f.PrintLabel(depth+1, "} // Pattern while\n")

	f.processPatternFail(depth+1, ctx)

	ctx.prevEntryPoint = ctx.entryPoint
	ctx.entryPoint++
	ctx.sentenceInfo.patternIndex++
}

func (f *Data) processEmptyPattern(depth int, ctx *emitterContext) {
	f.PrintLabel(depth+1, "if (currFrag->length > 0)")
	f.printFailBlock(depth+1, -1, false)
	f.PrintLabel(depth+1, "break;")
}

func (f *Data) processPattern(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.PrintLabel(depth+1, "//From what stretchable variable start?")
	f.PrintLabel(depth+1, "switch (stretchingVarNumber)")
	f.PrintLabel(depth+1, "{")

	ctx.patternCtx.entryPoint = 0
	ctx.patternCtx.prevEntryPoint = -1

	f.printFirstCase(depth, ctx, terms[0])

	f.matchingTerms(depth+2, false, ctx, terms)

	f.PrintLabel(depth+1, "} //pattern switch\n")

	f.PrintLabel(depth+1, "if (!stretching)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	f.printFailBlock(depth+2, ctx.patternCtx.prevEntryPoint, false)
	f.PrintLabel(depth+2, "else")
	f.PrintLabel(depth+3, "break; // Success!")
	f.PrintLabel(depth+1, "}")
}

func (f *Data) printFirstCase(depth int, ctx *emitterContext, term *syntax.Term) {

	if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
		if _, ok := ctx.fixedVars[term.Name]; !ok {
			return
		}
	}

	ctx.patternCtx.entryPoint = 1
	f.PrintLabel(depth+1, "case 0:")
}

func (f *Data) matchingTerms(depth int, inBrackets bool, ctx *emitterContext, terms []*syntax.Term) {
	parentMatchingOrder := ctx.isLeftMatching
	termsCount := len(terms)

	if termsCount == 0 {
		return
	}

	if terms[0].TermTag == syntax.R {
		ctx.isLeftMatching = false

		if inBrackets {
			f.PrintLabel(depth, "rightCheckOffset = fragmentOffset + memMngr.vterms[fragmentOffset-1].inBracketLength - 2;")
		}
	}

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

	if !ctx.isLeftMatching && inBrackets {
		f.PrintLabel(depth, "if (fragmentOffset < rightCheckOffset)")
		f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)
		f.PrintLabel(depth, "rightCheckOffset = fragmentOffset + currFrag->length;")
	}

	ctx.isLeftMatching = parentMatchingOrder
}

func (f *Data) matchingExpr(depth int, ctx *emitterContext, terms []*syntax.Term) {

	f.PrintLabel(depth, "//Check (")
	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, " || memMngr.vterms[fragmentOffset].tag != V_BRACKET_OPEN_TAG")

	f.PrintLabel(depth, "fragmentOffset++;")

	f.matchingTerms(depth, true, ctx, terms)

	f.PrintLabel(depth, "//Check )")
	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, " || memMngr.vterms[fragmentOffset].tag != V_BRACKET_CLOSE_TAG")

	f.PrintLabel(depth, "fragmentOffset++;")
}

func (f *Data) processPatternFail(depth int, ctx *emitterContext) {

	f.PrintLabel(depth, "if (stretchingVarNumber < 0)")
	f.PrintLabel(depth, "{")

	//First pattern in current sentence
	if ctx.sentenceInfo.patternIndex == 0 || ctx.prevEntryPoint == -1 {
		f.processFailOfFirstPattern(depth+1, ctx)
	} else {
		f.processFailOfCommonPattern(depth+1, ctx.prevEntryPoint)
	}

	f.PrintLabel(depth+1, "break;")
	f.PrintLabel(depth, "}")
}

func (f *Data) processFailOfFirstPattern(depth int, ctx *emitterContext) {
	if ctx.sentenceInfo.isLast {
		f.PrintLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		f.PrintLabel(depth, "_currFuncCall->entryPoint = -1;")

	} else {
		f.PrintLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		f.PrintLabel(depth, "stretching = 0;")
		f.PrintLabel(depth, fmt.Sprintf("_currFuncCall->entryPoint = %d;", ctx.nextSentenceEntryPoint))
		f.clearHelpers(depth, ctx.maxPatternNumber)
	}
}

func (f *Data) processFailOfCommonPattern(depth, prevEntryPoint int) {
	f.PrintLabel(depth, "//Jump to previouse pattern of same sentence!")
	f.PrintLabel(depth, fmt.Sprintf("_currFuncCall->entryPoint = %d;", prevEntryPoint))
}

func (f *Data) clearHelpers(depth, maxPatternNumber int) {

	if maxPatternNumber > 0 {
		f.PrintLabel(depth, "_currFuncCall->env->stretchVarsNumber[0] = 0;")

		f.PrintLabel(depth, fmt.Sprintf("for (i = 1; i < %d; ++i )", maxPatternNumber))
		f.PrintLabel(depth, "{")
		f.PrintLabel(depth+1, "_currFuncCall->env->stretchVarsNumber[i] = 0;")
		f.PrintLabel(depth+1, "_currFuncCall->env->fovs[i] = 0;")
		f.PrintLabel(depth+1, "_currFuncCall->env->assembledFOVs[i] = 0;")
		f.PrintLabel(depth, "}")
	}
}

func (f *Data) checkAndAssemblyChain(depth int, ctx *emitterContext) {
	patternIndex := ctx.sentenceInfo.patternIndex

	f.PrintLabel(depth, "if (!stretching)")
	f.PrintLabel(depth, "{")

	if patternIndex == 0 {
		if ctx.sentenceInfo.index == 0 {
			f.PrintLabel(depth+1, "// First sentence in func, first pattern in sentence.")
			f.PrintLabel(depth+1, "_currFuncCall->env->fovs[0] = _currFuncCall->fieldOfView;")
			f.PrintLabel(depth+1, "_currFuncCall->env->assembledFOVs[0] = gcGetAssembliedChain(_currFuncCall->fieldOfView);")
			f.PrintLabel(depth+1, "_currFuncCall->fieldOfView = 0;")
		} else {
			f.PrintLabel(depth+1, "// First pattern in sentence. Get fov, assembledFOVs from first pattern of prev sentence!")
		}
	} else {
		f.PrintLabel(depth+1, "// Pattern in middle of sentence.")
		f.PrintLabel(depth+1, "if (workFieldOfView != 0)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, "// There is assembly action in previous actions -> get this result.")
		f.PrintLabel(depth+2, fmt.Sprintf("_currFuncCall->env->fovs[%d] = workFieldOfView;", patternIndex))
		f.PrintLabel(depth+2, fmt.Sprintf("_currFuncCall->env->assembledFOVs[%d] = gcGetAssembliedChain(workFieldOfView);", patternIndex))
		f.PrintLabel(depth+2, "workFieldOfView = 0;")
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, "if (_currFuncCall->fieldOfView)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, fmt.Sprintf("_currFuncCall->env->fovs[%d] = _currFuncCall->fieldOfView;", patternIndex))
		f.PrintLabel(depth+3, fmt.Sprintf("_currFuncCall->env->assembledFOVs[%d] = gcGetAssembliedChain(_currFuncCall->fieldOfView);",
			patternIndex))
		f.PrintLabel(depth+3, "_currFuncCall->fieldOfView = 0;")
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "// There are no assemblies in previous actions => use prev pattern fieldOfView.")
		f.PrintLabel(depth+3, fmt.Sprintf("_currFuncCall->env->fovs[%d] = _currFuncCall->env->fovs[%d];",
			patternIndex, patternIndex-1))
		f.PrintLabel(depth+3, fmt.Sprintf("_currFuncCall->env->assembledFOVs[%d] = _currFuncCall->env->assembledFOVs[%d];",
			patternIndex, patternIndex-1))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+1, "}")
	}

	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, fmt.Sprintf("currFrag = _currFuncCall->env->assembledFOVs[%d]->fragment;", patternIndex))
	f.PrintLabel(depth, "rightCheckOffset = currFrag->offset + currFrag->length;")
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
	f.PrintLabel(depth-1, fmt.Sprintf("//Matching %s variable", value.Name))

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
			f.PrintLabel(depth-1, fmt.Sprintf("case %d:", ctx.patternCtx.entryPoint))

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

	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "stretching = 1;")
	f.PrintLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	if withBreakStatement {
		f.PrintLabel(depth+1, "break;")
	}
	f.PrintLabel(depth, "}")
}

func (f *Data) printOffsetCheck(depth, prevStretchVarNumber int, optionalCond string) {

	f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset >= rightCheckOffset%s)", optionalCond))
	f.printFailBlock(depth, prevStretchVarNumber, true)
}
