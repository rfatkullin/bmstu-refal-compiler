package emitter

import (
	"fmt"
)

func (f *Data) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, "")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length = 1;", varNumber))
	f.PrintLabel(depth+1, "fragmentOffset++;")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, " || memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG")

	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length = 1;", varNumber))
	f.PrintLabel(depth+1, "fragmentOffset++;")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {

	f.PrintLabel(depth, "if (!stretching) // Just init values")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->offset = fragmentOffset;", varNumber))

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length = 0;", varNumber))
	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length = rightCheckOffset - fragmentOffset;", varNumber))
		f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset += CURR_FUNC_CALL->env->locals[%d].fragment->length;", varNumber))
	}

	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else // stretching")

	f.varStretching(depth, varNumber, ctx)
}

func (f *Data) varStretching(depth, varNumber int, ctx *emitterContext) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, "stretching = 0;")
	f.PrintLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->stretchVarsNumber[%d] = %d;", patternNumber, ctx.patternCtx.entryPoint))

	f.PrintLabel(depth+1, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", ctx.bracketsIndex))

	f.PrintLabel(depth+1, "//Restore last offset at this point")
	f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset = CURR_FUNC_CALL->env->locals[%d].fragment->offset + "+
		" CURR_FUNC_CALL->env->locals[%d].fragment->length;", varNumber, varNumber))

	if ctx.isLeftMatching {
		f.printOffsetCheck(depth+1, prevStretchVarNumber, "")
	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("if (CURR_FUNC_CALL->env->locals[%d].fragment->length <= 0)", varNumber))
		f.printFailBlock(depth+1, prevStretchVarNumber, true)
	}

	f.PrintLabel(depth+1, "else")

	f.PrintLabel(depth+1, "{")

	if ctx.isLeftMatching {
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+2, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length++;", varNumber))
	} else {
		f.PrintLabel(depth+2, "fragmentOffset--;")
		f.PrintLabel(depth+2, fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d].fragment->length--;", varNumber))
	}

	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeVExprVar(depth int, ctx *emitterContext, varNumber int) {
	f.PrintLabel(depth, "if (!stretching) // Just init values")
	f.PrintLabel(depth, "{")
	f.matchingFreeTermVar(depth+1, ctx, varNumber)
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else // stretching")

	f.varStretching(depth, varNumber, ctx)
}
