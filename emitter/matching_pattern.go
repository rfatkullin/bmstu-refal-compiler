package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
	"BMSTU-Refal-Compiler/tokens"
)

func (f *Data) matchingPattern(depth int, ctx *emitterContext) {

	if len(ctx.terms) == 0 {
		return
	}

	f.PrintLabel(depth, fmt.Sprintf("//Sentence: %d, Pattern: %d", ctx.sentenceNumber, ctx.currPatternNumber))
	f.PrintLabel(depth, fmt.Sprintf("case %d:", ctx.currEntryPoint))
	f.PrintLabel(depth, fmt.Sprintf("{"))

	f.checkAndAssemblyChain(depth+1, ctx.currPatternNumber)

	f.PrintLabel(depth+1, "int fragmentOffset = currFrag->offset;")
	f.PrintLabel(depth+1, fmt.Sprintf("int stretchingVarNumber = stretchVarsNumber[%d];", ctx.currPatternNumber))
	f.PrintLabel(depth+1, "int stretching = 0;\n")

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

	f.PrintLabel(depth+2, "if (stretchingVarNumber >= 0)")
	f.PrintLabel(depth+2, "{")
	f.PrintLabel(depth+3, "if (fragmentOffset - fragmentTerm->fragment->offset < fragmentTerm->fragment->length)")
	f.PrintLabel(depth+4, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	f.PrintLabel(depth+3, "else")
	f.PrintLabel(depth+4, "break; // Success!")
	f.PrintLabel(depth+2, "}")

	f.PrintLabel(depth+1, "} // Pattern while\n")

	f.processPatternFail(depth+1, ctx)

	if ctx.isLastPatternInSentence {
		f.setToZeroStretchVarNumbers(depth, ctx)
	}

	ctx.currEntryPoint++
	ctx.currPatternNumber++
}

func (f *Data) processPatternFail(depth int, ctx *emitterContext) {

	f.PrintLabel(depth, "if (stretchingVarNumber < 0)")
	f.PrintLabel(depth, "{")

	//First pattern in current sentence
	if ctx.currPatternNumber == 0 {
		f.processFailOfFirstPattern(depth+1, ctx)
	} else {
		f.processFailOfCommonPattern(depth+1, ctx.currEntryPoint-1)
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

func (f *Data) setToZeroStretchVarNumbers(depth int, ctx *emitterContext) {

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i )", ctx.maxPatternNumber))
	f.PrintLabel(depth+1, "stretchVarsNumber[i] = 0;")
}

func (f *Data) checkAndAssemblyChain(depth, indexInSentence int) {
	f.PrintLabel(depth, fmt.Sprintf("if (assembledFOVs[%d] == 0)", indexInSentence))
	f.PrintLabel(depth+1, fmt.Sprintf("assembledFOVs[%d] = getAssembliedChain(fieldOfView->current);", indexInSentence))
	f.PrintLabel(depth, fmt.Sprintf("currFrag = assembledFOVs[%d]->frag;", indexInSentence))
}

func (f *Data) matchingVariable(depth int, ctx *emitterContext, value *tokens.Value, prevStretchVarNumber *int) {

	varNumber := ctx.sentenceScope.VarMap[value.Name].Number

	switch value.VarType {
	case tokens.VT_T:
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))
		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.printTermCheckFailBlock(depth+1, *prevStretchVarNumber)
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", varNumber))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_S:
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))

		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length || memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.printTermCheckFailBlock(depth+1, *prevStretchVarNumber)

		f.PrintLabel(depth+1, "else")

		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_E:

		f.PrintLabel(depth, fmt.Sprintf("case %d:", varNumber))
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))
		f.PrintLabel(depth+1, "if (!stretching) // Just init values")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 0;", varNumber))
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else // stretching")
		f.PrintLabel(depth+1, "{")

		f.PrintLabel(depth+2, "stretching = 0;")

		f.PrintLabel(depth+2, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.printTermCheckFailBlock(depth+2, *prevStretchVarNumber)

		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;;", varNumber))
		f.PrintLabel(depth+2, "}")

		f.PrintLabel(depth+2, "else")

		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += 1;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length += 1;", varNumber))
		f.PrintLabel(depth+2, "}")

		f.PrintLabel(depth+2, fmt.Sprintf("stretchVarsNumber[%d] = %d;", ctx.currPatternNumber, varNumber))
		f.PrintLabel(depth+1, "}")

		*prevStretchVarNumber = varNumber
		break
	}
}

func (f *Data) printTermCheckFailBlock(depth, prevStretchVarNumber int) {
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	f.PrintLabel(depth+1, "break;")
	f.PrintLabel(depth, "}")
}

func (f *Data) processSymbol(termNumber, depth int) {

}
