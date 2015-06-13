package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emt *EmitterData) checkNeedDataSize(depth int, terms []*syntax.Term) {
	chainsCount := 0
	vterms, data := emt.constrDataSizeStr(terms, &chainsCount)

	if len(vterms) > 0 {
		if vterms[len(vterms)-1] == '+' {
			vterms = vterms[0 : len(vterms)-1]
		}
	} else {
		vterms = "0"
	}

	if len(data) > 0 {
		if data[len(data)-1] == '+' {
			data = data[0 : len(data)-1]
		}
	} else {
		data = "0"
	}

	emt.printLabel(depth, fmt.Sprintf("if (GC_VTERM_OV(%s) || GC_DATA_OV(%s))", vterms, data))
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "collectGarbage();")
	emt.printLabel(depth+1, fmt.Sprintf("if (GC_VTERM_OV(%s) || GC_DATA_OV(%s))", vterms, data))
	emt.printLabel(depth+2, "PRINT_AND_EXIT(GC_MEMORY_OVERFLOW_MSG);")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) constrDataSizeStr(terms []*syntax.Term, chainsCount *int) (vterms, data string) {

	vterms = ""
	data = ""

	*chainsCount++

	for 0 < len(terms) {

		term := terms[0]

		if emt.isLiteral(term) {
			terms = emt.constrLiteralsDataSize(terms)
			data += "FRAGMENT_LTERM_SIZE+"
		} else {
			terms = terms[1:]

			if term.TermTag == syntax.EVAL {
				emt.ctx.isThereFuncCall = true
				data += "FUNC_CALL_LTERM_SIZE+"
				addVterms, addData := emt.constrDataSizeStr(term.Exprs[0].Terms, chainsCount)
				vterms += addVterms
				data += addData
			}

			if term.TermTag == syntax.EXPR {
				addVterms, addData := emt.constrDataSizeStr(term.Exprs[0].Terms, chainsCount)
				vterms += addVterms
				data += addData
			}

			if term.TermTag == syntax.VAR {
				data += "FRAGMENT_LTERM_SIZE+"
			}

			if term.TermTag == syntax.COMP || term.TermTag == syntax.FUNC {
				envSize := 0

				if term.Function != nil {
					envSize = len(term.Function.Env)
				}

				vterms += "1+"
				data += fmt.Sprintf("FRAGMENT_LTERM_SIZE+VCLOSURE_SIZE(%d)+", envSize)
			}
		}
	}

	if emt.ctx.isThereFuncCall {
		data += "CHAIN_LTERM_SIZE+"
	}

	return vterms, data
}

func (emt *EmitterData) constrLiteralsDataSize(terms []*syntax.Term) []*syntax.Term {

	literalsNumber := 0

	for _, term := range terms {

		if !emt.isLiteral(term) {
			break
		}

		literalsNumber++
	}

	return terms[literalsNumber:]
}
