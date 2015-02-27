package emitter

import (
	"fmt"
)

func (f *Data) matchingFixedLocalSymbolVar(depth int, ctx *emitterContext, matchedEntryPoint, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("env->locals[%d][%d]", matchedEntryPoint, varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedEnvSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	lterm := fmt.Sprintf("env->params[%d]", varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedLocalExprVar(depth int, ctx *emitterContext, patternNumber, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("env->locals[%d][%d]", patternNumber, varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedEnvExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("env->params[%d]", varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedExprVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {
	f.printOffsetCheck(depth, prevStretchVarNumber, "")

	checkTerm := ""
	patternTerm := ""

	if ctx.isLeftMatching {
		checkTerm = "memMngr.vterms[fragmentOffset + i]"
		patternTerm = fmt.Sprintf("memMngr.vterms[%s.fragment->offset + i]", lterm)
		f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset + %s.fragment->length >= currFrag->offset + currFrag->length)", lterm))
		f.printFailBlock(depth, prevStretchVarNumber, true)
	} else {
		checkTerm = "memMngr.vterms[fragmentOffset - i]"
		patternTerm = fmt.Sprintf("memMngr.vterms[%s.fragment->offset - i]", lterm)
		f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset - %s.fragment->length < leftCheckOffset)", lterm))
		f.printFailBlock(depth, prevStretchVarNumber, true)
	}

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %s.fragment->length; i++)", lterm))
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, fmt.Sprintf("if((%s.tag != %s.tag)", checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_CHAR_TAG && %s.ch != %s.ch)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_IDENT_TAG && !UStrCmp(%s.str, %s.str))", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_INT_NUM_TAG && %s.intNum != %s.intNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_FLOAT_NUM_TAG && %s.floatNum != %s.floatNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| ((%s.tag == V_BRACKET_OPEN_TAG || %s.tag == V_BRACKET_CLOSE_TAG) && %s.inBracketLength != %s.inBracketLength))", checkTerm, checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, "break;")

	f.PrintLabel(depth, "}")

	f.PrintLabel(depth, fmt.Sprintf("if(i < %s.fragment->length)", lterm))
	f.printFailBlock(depth, prevStretchVarNumber, true)

	if ctx.isLeftMatching {
		f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += %s.fragment->length;", lterm))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("fragmentOffset -= %s.fragment->length;", lterm))
	}
}

func (f *Data) matchingFixedSymbolVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {
	checkTerm := "memMngr.vterms[fragmentOffset]"
	patternTerm := fmt.Sprintf("memMngr.vterms[%s.fragment->offset]", lterm)

	if ctx.isLeftMatching {
		f.PrintLabel(depth, "if (fragmentOffset >= currFrag->offset + currFrag->length ")
	} else {
		f.PrintLabel(depth, "if (fragmentOffset < leftCheckOffset ")
	}

	f.PrintLabel(depth, fmt.Sprintf("|| %s.tag == V_BRACKET_OPEN_TAG || %s.tag == V_BRACKET_CLOSE_TAG", checkTerm, checkTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag != %s.tag)", checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_CHAR_TAG && %s.ch != %s.ch)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_IDENT_TAG && !UStrCmp(%s.str, %s.str))", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_INT_NUM_TAG && %s.intNum != %s.intNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_FLOAT_NUM_TAG && %s.floatNum != %s.floatNum))", checkTerm, checkTerm, patternTerm))

	f.printFailBlock(depth, prevStretchVarNumber, true)

	if ctx.isLeftMatching {
		f.PrintLabel(depth, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth, "fragmentOffset--;")
	}
}
