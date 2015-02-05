package emitter

import (
	"fmt"
)

func (f *Data) matchingFreeTermVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {
	f.printOffsetCheck(depth, prevStretchVarNumber, "")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth+1, "else")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->length = 1;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "fragmentOffset++;")
	f.PrintLabel(depth+1, "}")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeSymbolVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {
	f.printOffsetCheck(depth, prevStretchVarNumber, " || memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG")
	f.PrintLabel(depth, "else")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->length = 1;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "fragmentOffset++;")
	f.PrintLabel(depth, "}")
}

func (f *Data) matchingFreeExprVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {
	f.PrintLabel(depth, "if (!stretching) // Just init values")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->offset = fragmentOffset;", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals[%d][%d].fragment->length = 0;", patternNumber, varNumber))
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else // stretching")
	f.PrintLabel(depth, "{")

	f.printOffsetCheck(depth+1, prevStretchVarNumber, "")

	f.PrintLabel(depth+1, "stretching = 0;")
	f.PrintLabel(depth+1, fmt.Sprintf("env->stretchVarsNumber[%d] = %d;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "//Restore last offset at this point")
	f.PrintLabel(depth+1, fmt.Sprintf("fragmentOffset = env->locals[%d][%d].fragment->offset + env->locals[%d][%d].fragment->length;", patternNumber, varNumber, patternNumber, varNumber))

	f.PrintLabel(depth+1, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth+1, "else")

	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "fragmentOffset += 1;")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d][%d].fragment->length += 1;", patternNumber, varNumber))
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth, "}")
}
