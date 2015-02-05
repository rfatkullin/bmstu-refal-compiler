package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
	"BMSTU-Refal-Compiler/tokens"
)

func (f *Data) matchingPattern(depth int, ctx *emitterContext) {

	f.PrintLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceNumber, ctx.patternNumber))
	f.PrintLabel(depth, fmt.Sprintf("case %d:", ctx.entryPoint))
	f.PrintLabel(depth, fmt.Sprintf("{"))

	f.checkAndAssemblyChain(depth+1, ctx.patternNumber)

	f.PrintLabel(depth+1, "fragmentOffset = currFrag->offset;")
	f.PrintLabel(depth+1, fmt.Sprintf("stretchingVarNumber = env->stretchVarsNumber[%d];", ctx.patternNumber))
	f.PrintLabel(depth+1, "stretching = 0;\n")

	if len(ctx.terms) == 0 {
		f.processEmptyPattern(depth, ctx)
	} else {
		f.processPattern(depth, ctx)
	}

	f.processPatternFail(depth+1, ctx)

	if ctx.isLastPatternInSentence && !ctx.isLastSentence {
		f.initSretchVarNumbers(depth+1, ctx.maxPatternNumber)
	}

	ctx.entryPoint++
	ctx.patternNumber++
}

func (f *Data) processEmptyPattern(depth int, ctx *emitterContext) {
	f.PrintLabel(depth+1, "if (currFrag->length > 0)")
	f.printFailBlock(depth+1, -1, false)
}

func (f *Data) processPattern(depth int, ctx *emitterContext) {

	f.PrintLabel(depth+1, "while (stretchingVarNumber >= 0)")
	f.PrintLabel(depth+1, "{")

	f.PrintLabel(depth+2, "//From what stretchable variable start?")
	f.PrintLabel(depth+2, "switch (stretchingVarNumber)")
	f.PrintLabel(depth+2, "{")

	prevStretchVarNumber := -1
	for _, term := range ctx.terms {

		switch term.TermTag {
		case syntax.VAR:
			f.matchingVariable(depth+2, ctx, &term.Value, &prevStretchVarNumber)
			break
		}
	}

	f.PrintLabel(depth+2, "} //pattern switch\n")

	f.PrintLabel(depth+2, "if (!stretching)")
	f.PrintLabel(depth+2, "{")
	f.PrintLabel(depth+3, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	f.printFailBlock(depth+3, prevStretchVarNumber, false)
	f.PrintLabel(depth+3, "else")
	f.PrintLabel(depth+4, "break; // Success!")
	f.PrintLabel(depth+2, "}")

	f.PrintLabel(depth+1, "} // Pattern while\n")
}

func (f *Data) processPatternFail(depth int, ctx *emitterContext) {

	f.PrintLabel(depth, "if (stretchingVarNumber < 0)")
	f.PrintLabel(depth, "{")

	//First pattern in current sentence
	if ctx.patternNumber == 0 {
		f.processFailOfFirstPattern(depth+1, ctx)
	} else {
		f.processFailOfCommonPattern(depth+1, ctx.entryPoint-1)
	}

	f.PrintLabel(depth+1, "break;")
	f.PrintLabel(depth, "}")
}

func (f *Data) processFailOfFirstPattern(depth int, ctx *emitterContext) {
	if ctx.isLastSentence {
		f.PrintLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		f.PrintLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		f.PrintLabel(depth, "entryPoint = -1;")

	} else {
		f.PrintLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		f.PrintLabel(depth, fmt.Sprintf("entryPoint = %d;", ctx.nextSentenceEntryPoint))
	}
}

func (f *Data) processFailOfCommonPattern(depth, prevEntryPoint int) {
	f.PrintLabel(depth, "//Jump to previouse pattern of same sentence!")
	f.PrintLabel(depth, fmt.Sprintf("entryPoint = %d;", prevEntryPoint))
}

func (f *Data) initSretchVarNumbers(depth, maxPatternNumber int) {

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i )", maxPatternNumber))
	f.PrintLabel(depth+1, "env->stretchVarsNumber[i] = 0;")
}

func (f *Data) checkAndAssemblyChain(depth, indexInSentence int) {
	f.PrintLabel(depth, fmt.Sprintf("if (env->assembledFOVs[%d] == 0)", indexInSentence))
	f.PrintLabel(depth+1, fmt.Sprintf("env->assembledFOVs[%d] = getAssembliedChain(fieldOfView);", indexInSentence))
	f.PrintLabel(depth, fmt.Sprintf("currFrag = env->assembledFOVs[%d]->fragment;", indexInSentence))
}

func (f *Data) matchingVariable(depth int, ctx *emitterContext, value *tokens.Value, prevStretchVarNumber *int) {

	varNumber := ctx.sentenceScope.VarMap[value.Name].Number
	_, isFixedVar := ctx.fixedVars[value.Name]
	f.PrintLabel(depth, fmt.Sprintf("case %d: //Matching %s variable", varNumber, value.Name))
	f.PrintLabel(depth, "{")

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			f.matchingFixedExprVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
		} else {
			f.matchingFreeTermVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
		}
		break

	case tokens.VT_S:
		if isFixedVar {
			f.matchingFixedSymbolVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
		} else {
			f.matchingFreeSymbolVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
		}
		break

	case tokens.VT_E:
		if isFixedVar {
			f.matchingFixedExprVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
		} else {
			f.matchingFreeExprVar(depth+1, *prevStretchVarNumber, ctx.patternNumber, varNumber)
			*prevStretchVarNumber = varNumber
		}
		break
	}

	f.PrintLabel(depth, fmt.Sprintf("} // Matching %s variable", value.Name))
	ctx.fixedVars[value.Name] = true
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
	f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset >= currFrag->offset + currFrag->length%s)", optionalCond))
	f.printFailBlock(depth, prevStretchVarNumber, true)
}
