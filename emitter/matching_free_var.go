package emitter

import (
	"fmt"
)

func (emt *EmitterData) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	emt.printOffsetCheck(depth, prevStretchVarNumber, "")
	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emt.printLabel(depth+1, "fragmentOffset++;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	emt.printOffsetCheck(depth, prevStretchVarNumber, " || _memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG")

	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emt.printLabel(depth+1, "fragmentOffset++;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {

	emt.printLabel(depth, "if (!stretching) // Just init values")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))

	if ctx.isLeftMatching {
		emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 0;", varNumber))
	} else {
		emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = rightBound - fragmentOffset;", varNumber))
		emt.printLabel(depth+1, fmt.Sprintf("fragmentOffset += (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber))
	}

	emt.printLabel(depth, "}")
	emt.printLabel(depth, "else // stretching")

	emt.varStretching(depth, varNumber, ctx)
}

func (emt *EmitterData) varStretching(depth, varNumber int, ctx *emitterContext) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	emt.printLabel(depth, "{")

	emt.printLabel(depth+1, "stretching = 0;")
	emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->stretchVarsNumber[%d] = %d;", patternNumber, ctx.patternCtx.entryPoint))

	emt.printLabel(depth+1, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", ctx.bracketsCurrentIndex))

	emt.printLabel(depth+1, "//Restore last offset at this point")
	emt.printLabel(depth+1, fmt.Sprintf("fragmentOffset = (CURR_FUNC_CALL->env->locals + %d)->offset + "+
		" (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber, varNumber))

	if ctx.isLeftMatching {
		emt.printOffsetCheck(depth+1, prevStretchVarNumber, "")
	} else {
		emt.printLabel(depth+1, fmt.Sprintf("if ((CURR_FUNC_CALL->env->locals + %d)->length <= 0)", varNumber))
		emt.printFailBlock(depth+1, prevStretchVarNumber, true)
	}

	emt.printLabel(depth+1, "else")

	emt.printLabel(depth+1, "{")

	if ctx.isLeftMatching {
		emt.printLabel(depth+2, "fragmentOffset++;")
		emt.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length++;", varNumber))
	} else {
		emt.printLabel(depth+2, "fragmentOffset--;")
		emt.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length--;", varNumber))
	}

	emt.printLabel(depth+1, "}")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeVExprVar(depth int, ctx *emitterContext, varNumber int) {
	emt.printLabel(depth, "if (!stretching) // Just init values")
	emt.printLabel(depth, "{")
	emt.matchingFreeTermVar(depth+1, ctx, varNumber)
	emt.printLabel(depth, "}")
	emt.printLabel(depth, "else // stretching")

	emt.varStretching(depth, varNumber, ctx)
}
