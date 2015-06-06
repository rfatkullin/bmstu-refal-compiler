package emitter

import (
	"fmt"
)

func (emt *EmitterData) matchingFreeTermVar(depth int, varNumber int) {
	prevStretchVarNumber := emt.ctx.patternCtx.prevEntryPoint

	emt.printOffsetCheck(depth, prevStretchVarNumber, "")
	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emt.printLabel(depth+1, "fragmentOffset++;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeSymbolVar(depth int, varNumber int) {
	prevStretchVarNumber := emt.ctx.patternCtx.prevEntryPoint

	emt.printOffsetCheck(depth, prevStretchVarNumber, " || _memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG")

	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 1;", varNumber))
	emt.printLabel(depth+1, "fragmentOffset++;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeExprVar(depth int, varNumber int) {

	emt.printLabel(depth, "if (!stretching) // Just init values")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = fragmentOffset;", varNumber))

	if emt.ctx.isLeftMatching {
		emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = 0;", varNumber))
	} else {
		emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = rightBound - fragmentOffset;", varNumber))
		emt.printLabel(depth+1, fmt.Sprintf("fragmentOffset += (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber))
	}

	emt.printLabel(depth, "}")
	emt.printLabel(depth, "else // stretching")

	emt.varStretching(depth, varNumber)
}

func (emt *EmitterData) freeExprVarGetRest(depth int, varNumber int) {
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = CURR_FRAG_LEFT(UINT64_C(%d));", emt.ctx.brIndex))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = CURR_FRAG_LENGTH(UINT64_C(%d));", emt.ctx.brIndex))
}

func (emt *EmitterData) freeVExprVarGetRest(depth int, varNumber int) {
	prevStretchVarNumber := emt.ctx.patternCtx.prevEntryPoint

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printRollBackBlock(depth, prevStretchVarNumber, true)

	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->offset = CURR_FRAG_LEFT(UINT64_C(%d));", emt.ctx.brIndex))
	emt.printLabel(depth+1, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length = CURR_FRAG_LENGTH(UINT64_C(%d));", emt.ctx.brIndex))
}

func (emt *EmitterData) varStretching(depth, varNumber int) {
	prevStretchVarNumber := emt.ctx.patternCtx.prevEntryPoint
	patternNumber := emt.ctx.sentenceInfo.patternIndex

	emt.printLabel(depth, "{")

	emt.printLabel(depth+1, "stretching = 0;")
	emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->stretchVarsNumber[%d] = %d;", patternNumber, emt.ctx.patternCtx.entryPoint))

	emt.printLabel(depth+1, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", emt.ctx.brIndex))

	emt.printLabel(depth+1, "//Restore last offset at this point")
	emt.printLabel(depth+1, fmt.Sprintf("fragmentOffset = (CURR_FUNC_CALL->env->locals + %d)->offset + "+
		" (CURR_FUNC_CALL->env->locals + %d)->length;", varNumber, varNumber))

	if emt.ctx.isLeftMatching {
		emt.printOffsetCheck(depth+1, prevStretchVarNumber, "")
	} else {
		emt.printLabel(depth+1, fmt.Sprintf("if ((CURR_FUNC_CALL->env->locals + %d)->length <= 0)", varNumber))
		emt.printRollBackBlock(depth+1, prevStretchVarNumber, true)
	}

	emt.printLabel(depth+1, "else")

	emt.printLabel(depth+1, "{")

	if emt.ctx.isLeftMatching {
		emt.printLabel(depth+2, "fragmentOffset++;")
		emt.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length++;", varNumber))
	} else {
		emt.printLabel(depth+2, "fragmentOffset--;")
		emt.printLabel(depth+2, fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)->length--;", varNumber))
	}

	emt.printLabel(depth+1, "}")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) matchingFreeVExprVar(depth int, varNumber int) {
	emt.printLabel(depth, "if (!stretching) // Just init values")
	emt.printLabel(depth, "{")
	emt.matchingFreeTermVar(depth+1, varNumber)
	emt.printLabel(depth, "}")
	emt.printLabel(depth, "else // stretching")

	emt.varStretching(depth, varNumber)
}
