package emitter

import (
	"fmt"
)

func (f *Data) matchingFixedLocalSymbolVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {
	lterm := fmt.Sprintf("env->locals[%d][%d]", patternNumber, varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, lterm)
}

func (f *Data) matchingFixedEnvSymbolVar(depth, prevStretchVarNumber, varNumber int) {
	lterm := fmt.Sprintf("env->params[%d]", varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, lterm)
}

func (f *Data) matchingFixedLocalExprVar(depth, prevStretchVarNumber, patternNumber, varNumber int) {

	lterm := fmt.Sprintf("env->locals[%d][%d]", patternNumber, varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, lterm)
}

func (f *Data) matchingFixedEnvExprVar(depth, prevStretchVarNumber, varNumber int) {
	lterm := fmt.Sprintf("env->params[%d]", varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, lterm)
}

func (f *Data) matchingFixedExprVar(depth, prevStretchVarNumber int, lterm string) {
	f.printOffsetCheck(depth, prevStretchVarNumber, "")

	checkTerm := "memMngr.vterms[fragmentOffset + i]"
	patternTerm := fmt.Sprintf("memMngr.vterms[%s.fragment->offset + i]", lterm)

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %s.fragment->length; i++)", lterm))
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, fmt.Sprintf("if((%s.tag != %s.tag)", checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_CHAR_TAG && %s.ch != %s.ch)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_IDENT_TAG && strcmp(%s.str, %s.str))", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_INT_NUM_TAG && %s.intNum != %s.intNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_FLOAT_NUM_TAG && %s.floatNum != %s.floatNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_BRACKET_TAG && %s.inBracketLength != %s.inBracketLength))", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, "break;")

	f.PrintLabel(depth, "}")

	f.PrintLabel(depth, fmt.Sprintf("if(i < %s.fragment->length)", lterm))
	f.printFailBlock(depth, prevStretchVarNumber, true)

	f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += %s.fragment->length;", lterm))
}

func (f *Data) matchingFixedSymbolVar(depth, prevStretchVarNumber int, lterm string) {
	checkTerm := "memMngr.vterms[fragmentOffset]"
	patternTerm := fmt.Sprintf("memMngr.vterms[%s.fragment->offset]", lterm)

	f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset >= currFrag->offset + currFrag->length || %s.tag == V_BRACKET_TAG", checkTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag != %s.tag)", checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_CHAR_TAG && %s.ch != %s.ch)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_IDENT_TAG && strcmp(%s.str, %s.str))", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_INT_NUM_TAG && %s.intNum != %s.intNum)", checkTerm, checkTerm, patternTerm))
	f.PrintLabel(depth+1, fmt.Sprintf("|| (%s.tag == V_FLOAT_NUM_TAG && %s.floatNum != %s.floatNum))", checkTerm, checkTerm, patternTerm))

	f.printFailBlock(depth, prevStretchVarNumber, true)

	f.PrintLabel(depth, "fragmentOffset++;")
}
