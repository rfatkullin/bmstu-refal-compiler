package emitter

import (
	"fmt"
)

func (f *Data) matchingFixedSymbolVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {
	f.PrintLabel(depth, "if (fragmentOffset >= currFrag->length || memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG")
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag != memMngr.vterms[env->locals[%d][%d].fragment->offset].tag)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag == V_CHAR_TAG"+
		"  && memMngr.vterms[fragmentOffset].ch != memMngr.vterms[env->locals[%d][%d].fragment->offset].ch)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag == V_IDENT_TAG"+
		" && strcmp(memMngr.vterms[fragmentOffset].str, memMngr.vterms[env->locals[%d][%d].fragment->offset].str))", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag == V_INT_NUM_TAG"+
		" && memMngr.vterms[fragmentOffset].intNum != memMngr.vterms[env->locals[%d][%d].fragment->offset].intNum)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag == V_FLOAT_NUM_TAG"+
		" && memMngr.vterms[fragmentOffset].floatNum != memMngr.vterms[env->locals[%d][%d].fragment->offset].floatNum))", patternNumber, varNumber))
	f.printFailBlock(depth, prevStretchVarNumber)

	f.PrintLabel(depth, "fragmentOffset++;")
}

func (f *Data) matchingFixedExprVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {

	f.PrintLabel(depth, "if (fragmentOffset >= currFrag->length)")
	f.printFailBlock(depth, prevStretchVarNumber)

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < env->locals[%d][%d].fragment->length; i++)", patternNumber, varNumber))
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, fmt.Sprintf("if((memMngr.vterms[fragmentOffset + i].tag != memMngr.vterms[env->locals[%d][%d].fragment->offset + i].tag)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset + i].tag == V_CHAR_TAG"+
		"  && memMngr.vterms[fragmentOffset + i].ch != memMngr.vterms[env->locals[%d][%d].fragment->offset + i].ch)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset + i].tag == V_IDENT_TAG"+
		" && strcmp(memMngr.vterms[fragmentOffset + i].str, memMngr.vterms[env->locals[%d][%d].fragment->offset + i].str))", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset + i].tag == V_INT_NUM_TAG"+
		" && memMngr.vterms[fragmentOffset + i].intNum != memMngr.vterms[env->locals[%d][%d].fragment->offset + i].intNum)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset + i].tag == V_FLOAT_NUM_TAG"+
		" && memMngr.vterms[fragmentOffset + i].floatNum != memMngr.vterms[env->locals[%d][%d].fragment->offset + i].floatNum)", patternNumber, varNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset + i].tag == V_BRACKET_TAG"+
		" && memMngr.vterms[fragmentOffset + i].inBracketLength != memMngr.vterms[env->locals[%d][%d].fragment->offset + i].inBracketLength))", patternNumber, varNumber))
	f.PrintLabel(depth+1, "break;")

	f.PrintLabel(depth, "}")

	f.PrintLabel(depth, fmt.Sprintf("if(i < env->locals[%d][%d].fragment->length)", patternNumber, varNumber))
	f.printFailBlock(depth, prevStretchVarNumber)

	f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += env->locals[%d][%d].fragment->length;", patternNumber, varNumber))
}
