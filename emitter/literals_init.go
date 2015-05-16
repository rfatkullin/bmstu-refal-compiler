package emitter

import (
	"fmt"
	"unicode/utf8"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emt *EmitterData) printLiteralsAndHeapsInit(depth int, units []*syntax.Unit) {

	emt.printLabel(depth, "void initLiteralData()\n{")

	// Dummy-vterm. Обращение к vterm'у с нулевым смещением - признак ошибки.
	emt.printLabel(depth+1, "_memMngr.vterms[0] = (struct vterm_t){.tag = V_CHAR_TAG, .ch = 0}; // dummy-vterm.")

	for _, unit := range units {
		for _, currFunc := range unit.GlobMap {
			emt.initFuncLiterals(depth+1, currFunc)
		}
	}

	emt.printLabel(depth, "} // initLiteralData()\n")
}

func (emt *EmitterData) initActionLiterals(depth int, expr syntax.Expr) {

	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	for len(terms) > 0 {

		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.STR:
			emt.initStrVTerm(depth, term)
			break

		case syntax.COMP:
			emt.initIdentVTerm(depth, term, term.Value.Name)
			break

		case syntax.INT:
			emt.initIntNumVTerm(depth, term)
			break

		case syntax.FLOAT:
			emt.initFloatVTerm(depth, term)
			break

		case syntax.EXPR, syntax.EVAL:
			tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
			tmpTerms = append(tmpTerms, terms...)
			terms = tmpTerms
			break

		case syntax.FUNC:
			if !term.Function.HasName {
				term.Function.HasName = true
				term.FuncName = fmt.Sprintf("AnonymFunc_%d", emt.currTermNum)
				emt.initIdentVTerm(depth, term, term.FuncName)
			}

			emt.initIdentVTerm(depth, term, term.FuncName)
			emt.initFuncLiterals(depth, term.Function)
			break
		}
	}
}

func (emt *EmitterData) initFuncLiterals(depth int, currFunc *syntax.Function) {
	for _, s := range currFunc.Sentences {
		emt.initActionLiterals(depth, s.Pattern)
		for _, a := range s.Actions {
			emt.initActionLiterals(depth, a.Expr)
		}
	}
}

func (emt *EmitterData) initStrVTerm(depth int, term *syntax.Term) {
	term.IndexInLiterals = emt.currTermNum

	for i := 0; i < len(term.Value.Str); i++ {
		emt.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_CHAR_TAG, .ch = %d};", emt.currTermNum, term.Value.Str[i]))
		emt.currTermNum++
	}
}

// Инициализация vterm_t для литералов целого типа
// Пока только обычные
func (emt *EmitterData) initIntNumVTerm(depth int, term *syntax.Term) {
	bytesStr, sign, bytesCount := getStrOfBytes(term.Value.Int)

	emt.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_INT_NUM_TAG,"+
		" .intNum = allocateIntNumberLiteral((uint8_t[]){%s}, %d, UINT64_C(%d))};",
		emt.currTermNum, bytesStr, sign, bytesCount))

	term.IndexInLiterals = emt.currTermNum
	emt.currTermNum++
}

// Инициализация vterm_t для литералов вещественного типа
func (emt *EmitterData) initFloatVTerm(depth int, term *syntax.Term) {

	emt.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_DOUBLE_NUM_TAG, .doubleNum = %f};", emt.currTermNum, term.Value.Float))
	term.IndexInLiterals = emt.currTermNum
	emt.currTermNum++
}

// Инициализация vterm_t для идентификатора
func (emt *EmitterData) initIdentVTerm(depth int, term *syntax.Term, ident string) {
	runesStr := getStrOfRunes(ident)

	emt.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_IDENT_TAG, .str = allocateVStringLiteral((uint32_t[]){%s}, UINT64_C(%d))};",
		emt.currTermNum, runesStr, utf8.RuneCountInString(ident)))

	term.IndexInLiterals = emt.currTermNum
	emt.currTermNum++
}
