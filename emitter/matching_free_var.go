package emitter

import (
	"fmt"
)

func (emitter *EmitterData) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	emitter.printOffsetCheck(depth, prevStretchVarNumber, "")
	emitter.printLabel(depth, "else")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emitter.printLabel(depth+1, "fragmentOffset++;")
	emitter.printLabel(depth, "}")
}

func (emitter *EmitterData) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	emitter.printOffsetCheck(depth, prevStretchVarNumber, " || _memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG")

	emitter.printLabel(depth, "else")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emitter.printLabel(depth+1, "fragmentOffset++;")
	emitter.printLabel(depth, "}")
}

func (emitter *EmitterData) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {

	emitter.printLabel(depth, "if (!stretching) // Just init values")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))

	if ctx.isLeftMatching {
		emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 0;", varNumber))
	} else {
		emitter.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = rightBound - fragmentOffset;", varNumber))
		emitter.printLabel(depth+1, fmt.Sprintf("fragmentOffset += (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber))
	}

	emitter.printLabel(depth, "}")
	emitter.printLabel(depth, "else // stretching")

	emitter.varStretching(depth, varNumber, ctx)
}

func (emitter *EmitterData) varStretching(depth, varNumber int, ctx *emitterContext) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	emitter.printLabel(depth, "{")

	emitter.printLabel(depth+1, "stretching = 0;")
	emitter.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->stretchVarsNumber[%d] = %d;", patternNumber, ctx.patternCtx.entryPoint))

	emitter.printLabel(depth+1, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", ctx.bracketsCurrentIndex))

	emitter.printLabel(depth+1, "//Restore last offset at this point")
	emitter.printLabel(depth+1, fmt.Sprintf("fragmentOffset = (CURR_FUNC_CALL->env->locals + %d)->offset + "+
		" (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber, varNumber))

	if ctx.isLeftMatching {
		emitter.printOffsetCheck(depth+1, prevStretchVarNumber, "")
	} else {
		emitter.printLabel(depth+1, fmt.Sprintf("if ((CURR_FUNC_CALL->env->locals + %d)->length <= 0)", varNumber))
		emitter.printFailBlock(depth+1, prevStretchVarNumber, true)
	}

	emitter.printLabel(depth+1, "else")

	emitter.printLabel(depth+1, "{")

	if ctx.isLeftMatching {
		emitter.printLabel(depth+2, "fragmentOffset++;")
		emitter.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length++;", varNumber))
	} else {
		emitter.printLabel(depth+2, "fragmentOffset--;")
		emitter.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length--;", varNumber))
	}

	emitter.printLabel(depth+1, "}")
	emitter.printLabel(depth, "}")
}

func (emitter *EmitterData) matchingFreeVExprVar(depth int, ctx *emitterContext, varNumber int) {
	emitter.printLabel(depth, "if (!stretching) // Just init values")
	emitter.printLabel(depth, "{")
	emitter.matchingFreeTermVar(depth+1, ctx, varNumber)
	emitter.printLabel(depth, "}")
	emitter.printLabel(depth, "else // stretching")

	emitter.varStretching(depth, varNumber, ctx)
}
