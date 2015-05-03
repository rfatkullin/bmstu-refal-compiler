package emitter

import (
	"fmt"
)

func (f *Data) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, "")
	f.printLabel(depth, "else")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	f.printLabel(depth+1, "fragmentOffset++;")
	f.printLabel(depth, "}")
}

func (f *Data) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, " || _memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG")

	f.printLabel(depth, "else")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	f.printLabel(depth+1, "fragmentOffset++;")
	f.printLabel(depth, "}")
}

func (f *Data) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {

	f.printLabel(depth, "if (!stretching) // Just init values")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))

	if ctx.isLeftMatching {
		f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 0;", varNumber))
	} else {
		f.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = rightBound - fragmentOffset;", varNumber))
		f.printLabel(depth+1, fmt.Sprintf("fragmentOffset += (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber))
	}

	f.printLabel(depth, "}")
	f.printLabel(depth, "else // stretching")

	f.varStretching(depth, varNumber, ctx)
}

func (f *Data) varStretching(depth, varNumber int, ctx *emitterContext) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.printLabel(depth, "{")

	f.printLabel(depth+1, "stretching = 0;")
	f.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->stretchVarsNumber[%d] = %d;", patternNumber, ctx.patternCtx.entryPoint))

	f.printLabel(depth+1, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", ctx.bracketsCurrentIndex))

	f.printLabel(depth+1, "//Restore last offset at this point")
	f.printLabel(depth+1, fmt.Sprintf("fragmentOffset = (CURR_FUNC_CALL->env->locals + %d)->offset + "+
		" (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber, varNumber))

	if ctx.isLeftMatching {
		f.printOffsetCheck(depth+1, prevStretchVarNumber, "")
	} else {
		f.printLabel(depth+1, fmt.Sprintf("if ((CURR_FUNC_CALL->env->locals + %d)->length <= 0)", varNumber))
		f.printFailBlock(depth+1, prevStretchVarNumber, true)
	}

	f.printLabel(depth+1, "else")

	f.printLabel(depth+1, "{")

	if ctx.isLeftMatching {
		f.printLabel(depth+2, "fragmentOffset++;")
		f.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length++;", varNumber))
	} else {
		f.printLabel(depth+2, "fragmentOffset--;")
		f.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length--;", varNumber))
	}

	f.printLabel(depth+1, "}")
	f.printLabel(depth, "}")
}

func (f *Data) matchingFreeVExprVar(depth int, ctx *emitterContext, varNumber int) {
	f.printLabel(depth, "if (!stretching) // Just init values")
	f.printLabel(depth, "{")
	f.matchingFreeTermVar(depth+1, ctx, varNumber)
	f.printLabel(depth, "}")
	f.printLabel(depth, "else // stretching")

	f.varStretching(depth, varNumber, ctx)
}
