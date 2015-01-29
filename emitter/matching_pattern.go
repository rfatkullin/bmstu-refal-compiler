package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
	"BMSTU-Refal-Compiler/tokens"
)

func (f *Data) matchingPattern(depth, patternNumber, allPatternNumber int, p *syntax.Expr, scope *syntax.Scope, currEntryPoint *int) {

	if len(p.Terms) == 0 {
		return
	}

	f.PrintLabel(depth, fmt.Sprintf("case %d:", *currEntryPoint))
	f.PrintLabel(depth, fmt.Sprintf("{"))

	f.checkAndAssemblyChain(depth+1, patternNumber)

	f.PrintLabel(depth+1, "int fragmentOffset = currFrag->offset;")
	f.PrintLabel(depth+1, fmt.Sprintf("int stretchingVarNumber = stretchVarsNumber[%d];", patternNumber))
	f.PrintLabel(depth+1, "int stretching = 0;\n")

	f.PrintLabel(depth+1, "while (stretchingVarNumber >= 0)")
	f.PrintLabel(depth+1, "{")

	f.PrintLabel(depth+2, "//From what stretchable variable start?")
	f.PrintLabel(depth+2, "switch (stretchingVarNumber)")
	f.PrintLabel(depth+2, "{")

	prevStretchVarNumber := -1
	for _, term := range p.Terms {

		switch term.TermTag {
		case syntax.VAR:
			f.matchingVariable(depth+2, patternNumber, &term.Value, scope, &prevStretchVarNumber)
			break
		}
	}

	f.PrintLabel(depth+1, "} //pattern switch\n")

	f.PrintLabel(depth+1, "if (stretchingVarNumber >= 0)")
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "if (fragmentOffset - fragmentTerm->fragment->offset < fragmentTerm->fragment->length)")
	f.PrintLabel(depth+3, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	f.PrintLabel(depth+2, "else")
	f.PrintLabel(depth+3, "break; // Success!")
	f.PrintLabel(depth+1, "}")

	f.PrintLabel(depth, "}")
}

func (f *Data) initStretchVarNumbersArray(depth, matchingNumber int) {

	f.PrintLabel(depth, "//TO FIX: Set to zero after every sentence!")
	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; ++i )", matchingNumber))
	f.PrintLabel(depth+1, "stretchVarsNumber[i] = 0;")
}

func (f *Data) checkAndAssemblyChain(depth, indexInSentence int) {
	f.PrintLabel(depth, fmt.Sprintf("if (assembledFOVs[%d] == 0)", indexInSentence))
	f.PrintLabel(depth+1, fmt.Sprintf("assembledFOVs[%d] = getAssembliedChain(fieldOfView->current);", indexInSentence))
	f.PrintLabel(depth, fmt.Sprintf("currFrag = assembledFOVs[%d]->frag;", indexInSentence))
}

func (f *Data) matchingVariable(depth, patternNumber int, value *tokens.Value, scope *syntax.Scope, prevStretchVarNumber *int) {

	varNumber := scope.VarMap[value.Name].Number

	switch value.VarType {
	case tokens.VT_T:
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))
		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.printTermCheckFailBlock(depth+1, *prevStretchVarNumber)
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length = memMngr.vterms[fragmentOffset].inBracketLength;", varNumber))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_S:
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))

		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length || memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.printTermCheckFailBlock(depth+1, *prevStretchVarNumber)

		f.PrintLabel(depth+1, "else")

		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_E:

		f.PrintLabel(depth, fmt.Sprintf("case %d:", varNumber))
		f.PrintLabel(depth+1, fmt.Sprintf("//Matching %s variable", value.Name))
		f.PrintLabel(depth+1, "if (!stretching) // Just init values")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 0;", varNumber))
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else // stretching")
		f.PrintLabel(depth+1, "{")

		f.PrintLabel(depth+2, "stretching = 0;")

		f.PrintLabel(depth+2, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.printTermCheckFailBlock(depth+2, *prevStretchVarNumber)

		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length += memMngr.vterms[fragmentOffset].inBracketLength;;", varNumber))
		f.PrintLabel(depth+2, "}")

		f.PrintLabel(depth+2, "else")

		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += 1;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length += 1;", varNumber))
		f.PrintLabel(depth+2, "}")

		f.PrintLabel(depth+2, fmt.Sprintf("stretchVarsNumber[%d] = %d;", patternNumber, varNumber))
		f.PrintLabel(depth+1, "}")

		*prevStretchVarNumber = varNumber
		break
	}
}

func (f *Data) printTermCheckFailBlock(depth, prevStretchVarNumber int) {
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	f.PrintLabel(depth+1, "break;")
	f.PrintLabel(depth, "}")
}

func (f *Data) processSymbol(termNumber, depth int) {

}
