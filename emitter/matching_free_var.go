package emitter

import (
	"fmt"
)

func (f *Data) matchingFreeTermVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, "")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))

	f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", varNumber))
	f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth+1, "else")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))

	f.PrintLabel(depth+2, "fragmentOffset++;")

	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	f.printOffsetCheck(depth, prevStretchVarNumber, " || memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG || memMngr.vterms[fragmentOffset].tag == V_BRACKET_CLOSE_TAG")

	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
	f.PrintLabel(depth+1, "fragmentOffset++;")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	patternNumber := ctx.sentenceInfo.patternIndex

	f.PrintLabel(depth, "if (!stretching) // Just init values")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->length = 0;", varNumber))
	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d].fragment->length = rightCheckOffset - fragmentOffset;", varNumber))
		f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset += env->locals[%d].fragment->length;", varNumber))
	}

	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else // stretching")
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, "stretching = 0;")
	f.PrintLabel(depth+1, fmt.Sprintf("env->stretchVarsNumber[%d] = %d;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "//Restore last offset at this point")

	f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset = env->locals[%d].fragment->offset + env->locals[%d].fragment->length;", varNumber, varNumber))

	if ctx.isLeftMatching {
		f.printOffsetCheck(depth+1, prevStretchVarNumber, "")
		f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_OPEN_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+1, "}")

	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("if (env->locals[%d].fragment->length <= 0)", varNumber))
		f.printFailBlock(depth+1, prevStretchVarNumber, true)
	}

	f.PrintLabel(depth+1, "else")

	f.PrintLabel(depth+1, "{")

	if ctx.isLeftMatching {
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length++;", varNumber))
	} else {
		f.PrintLabel(depth+2, "fragmentOffset--;")

		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_CLOSE_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length -= memMngr.vterms[fragmentOffset].inBracketLength;", varNumber))
		f.PrintLabel(depth+3, "fragmentOffset -= (memMngr.vterms[fragmentOffset].inBracketLength - 1);")
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length--;", varNumber))
		f.PrintLabel(depth+2, "}")
	}

	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth, "}")
}
