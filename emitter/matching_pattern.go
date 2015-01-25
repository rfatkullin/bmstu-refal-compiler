package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
	"BMSTU-Refal-Compiler/tokens"
)

func (f *Data) matchingPattern(depth int, p *syntax.Expr, scope *syntax.Scope) {

	if len(p.Terms) == 0 {
		return
	}

	f.PrintLabel(depth, "fieldOfView->current = getAssembliedChain(fieldOfView->current);")
	f.PrintLabel(depth, "struct lterm_t* fragmentTerm = fieldOfView->current->begin;")
	f.PrintLabel(depth, "int fragmentBegin = fragmentTerm->fragment->offset;")
	f.PrintLabel(depth, "int fragmentOffset = 0;")
	f.PrintLabel(depth, "int firstIteration = 1;")

	f.PrintLabel(depth, "while (1)")
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, "switch (s)")
	for _, term := range p.Terms {

		switch term.TermTag {
		case syntax.VAR:
			f.matchingVariable(depth+2, &term.Value, scope)
			break
		}

	}
	f.PrintLabel(depth+1, "} //patterns switch")

	f.PrintLabel(depth, "}")
}

func (f *Data) matchingVariable(depth int, value *tokens.Value, scope *syntax.Scope) {

	varNumber := scope.VarMap[value.Name]
	prevStretchVarNumber := -1

	switch value.VarType {
	case tokens.VT_T:
		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("result = %d;", prevStretchVarNumber))
		f.PrintLabel(depth+2, "break;")
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length = inBracketLength;", varNumber))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_S:
		f.PrintLabel(depth+1, "if (fragmentOffset >= fragmentTerm->fragment->length || memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("result = %d;", prevStretchVarNumber))
		f.PrintLabel(depth+2, "break;")
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 1;", varNumber))
		f.PrintLabel(depth+2, "fragmentOffset++;")
		f.PrintLabel(depth+1, "}")
		break

	case tokens.VT_E:

		f.PrintLabel(depth, fmt.Sprintf("case %d:", varNumber))
		f.PrintLabel(depth+1, "if (!stretching)")
		f.PrintLabel(depth+1, "{")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->offset = fragmentOffset;", varNumber))
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length = 0;", varNumber))
		f.PrintLabel(depth+1, "}")
		f.PrintLabel(depth+1, "else")
		f.PrintLabel(depth+1, "{")

		f.PrintLabel(depth+2, "if (fragmentOffset >= fragmentTerm->fragment->length)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, fmt.Sprintf("result = %d;", prevStretchVarNumber))
		f.PrintLabel(depth+3, "break;")
		f.PrintLabel(depth+2, "}")

		f.PrintLabel(depth+2, "if (memMngr.vterms[fragmentOffset].tag == V_BRACKET_TAG)")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += memMngr.vterms[fragmentOffset].inBracketLength;")
		f.PrintLabel(depth+3, fmt.Sprintf("env->locals[%d].fragment->length += inBracketLength;", varNumber))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+2, "else")
		f.PrintLabel(depth+2, "{")
		f.PrintLabel(depth+3, "fragmentOffset += 1;")
		f.PrintLabel(depth+2, fmt.Sprintf("env->locals[%d].fragment->length += 1;", varNumber))
		f.PrintLabel(depth+2, "}")
		f.PrintLabel(depth+1, "}")
		break
	}
}

func (f *Data) processSymbol(termNumber, depth int) {

}
