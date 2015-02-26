package emitter

import (
	"fmt"
)

func (f *Data) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.printOffsetCheck(depth, prevStretchVarNumber, "")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
		f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+1, "}")
	} else {
		f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_CLOSE_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
		f.PrintLabel(depth+2, "fragmentOffset -= memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+1, "}")
	}

	f.PrintLabel(depth+1, "else")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length = 1;", patternNumber, varNumber))

	if ctx.isLeftMatching {
		f.PrintLabel(depth+2, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth+2, "fragmentOffset--;")
	}

	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.printOffsetCheck(depth, prevStretchVarNumber, " || memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG || memMngr.vterms[fragmentOffset].tag == V_BRACKET_CLOSE_TAG")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->length = 1;", patternNumber, varNumber))

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth+1, "fragmentOffset--;")
	}

	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.PrintLabel(depth, "if (!stretching) // Just init values")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->length = 0;", patternNumber, varNumber))
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else // stretching")
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, "stretching = 0;")
	f.PrintLabel(depth+1, fmt.Sprintf("env->stretchVarsNumber[%d] = %d;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "//Restore last offset at this point")

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset = env->locals[%d][%d].fragment->offset + env->locals[%d][%d].fragment->length;", patternNumber, varNumber, patternNumber, varNumber))
		f.PrintLabel(depth+1, "if (fragmentOffset >= currFrag->offset + currFrag->length)")
		f.printFailBlock(depth+1, prevStretchVarNumber, true)
		f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
		f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+1, "}")

	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset = env->locals[%d][%d].fragment->offset;", patternNumber, varNumber))
		f.PrintLabel(depth+1, "if (fragmentOffset < leftCheckOffset)")
		f.printFailBlock(depth+1, prevStretchVarNumber, true)
		f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_CLOSE_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
		f.PrintLabel(depth+2, "fragmentOffset -= memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+1, "}")
	}

	f.PrintLabel(depth+1, "else")

	f.PrintLabel(depth+1, "{")

	if ctx.isLeftMatching {
		f.PrintLabel(depth+2, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth+2, "fragmentOffset--;")
	}

	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length += 1;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth, "}")
}
