package emitter

import (
	"fmt"
	"unicode/utf8"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emitter *EmitterData) printLiteralsAndHeapsInit(depth int, unit *syntax.Unit) {

	emitter.printLabel(depth, "void initLiteralData()\n{")
	emitter.initLiterals(depth+1, unit.GlobMap)
	emitter.printLabel(depth, "} // initLiteralData()\n")
}

func (emitter *EmitterData) initLiterals(depth int, funcs map[string]*syntax.Function) {

	// Dummy-vterm. Обращение к vterm'у с нулевым смещением - признак ошибки.
	emitter.printLabel(depth, "_memMngr.vterms[0] = (struct vterm_t){.tag = V_CHAR_TAG, .ch = 0}; // dummy-vterm.")

	emitter.currTermNum = 1
	for _, currFunc := range funcs {
		emitter.initFuncLiterals(depth, currFunc)
	}

	fmt.Fprintf(emitter, "\n")
}

func (emitter *EmitterData) initActionLiterals(depth int, expr syntax.Expr) {

	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	for len(terms) > 0 {

		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.STR:
			emitter.initStrVTerm(depth, term)
			break

		case syntax.COMP:
			emitter.initIdentVTerm(depth, term, term.Value.Name)
			break

		case syntax.INT:
			emitter.initIntNumVTerm(depth, term)
			break

		case syntax.FLOAT:
			emitter.initFloatVTerm(depth, term)
			break

		case syntax.EXPR, syntax.EVAL:
			tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
			tmpTerms = append(tmpTerms, terms...)
			terms = tmpTerms
			break

		case syntax.FUNC:
			if !term.Function.HasName {
				term.Function.HasName = true
				term.FuncName = fmt.Sprintf("AnonymFunc_%d", emitter.currTermNum)
				emitter.initIdentVTerm(depth, term, term.FuncName)
			}

			emitter.initIdentVTerm(depth, term, term.FuncName)
			emitter.initFuncLiterals(depth, term.Function)
			break
		}
	}
}

func (emitter *EmitterData) initFuncLiterals(depth int, currFunc *syntax.Function) {
	for _, s := range currFunc.Sentences {
		emitter.initActionLiterals(depth, s.Pattern)
		for _, a := range s.Actions {
			emitter.initActionLiterals(depth, a.Expr)
		}
	}
}

func (emitter *EmitterData) initStrVTerm(depth int, term *syntax.Term) {
	term.IndexInLiterals = emitter.currTermNum

	for i := 0; i < len(term.Value.Str); i++ {
		emitter.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_CHAR_TAG, .ch = %d};", emitter.currTermNum, term.Value.Str[i]))
		emitter.currTermNum++
	}
}

// Инициализация vterm_t для литералов целого типа
// Пока только обычные
func (emitter *EmitterData) initIntNumVTerm(depth int, term *syntax.Term) {
	bytesStr, sign, bytesCount := getStrOfBytes(term.Value.Int)

	emitter.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_INT_NUM_TAG,"+
		" .intNum = allocateIntNumberLiteral((uint8_t[]){%s}, %d, UINT64_C(%d))};",
		emitter.currTermNum, bytesStr, sign, bytesCount))

	term.IndexInLiterals = emitter.currTermNum
	emitter.currTermNum++
}

// Инициализация vterm_t для литералов вещественного типа
func (emitter *EmitterData) initFloatVTerm(depth int, term *syntax.Term) {

	emitter.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_DOUBLE_NUM_TAG, .doubleNum = %f};", emitter.currTermNum, term.Value.Float))
	term.IndexInLiterals = emitter.currTermNum
	emitter.currTermNum++
}

// Инициализация vterm_t для идентификатора
func (emitter *EmitterData) initIdentVTerm(depth int, term *syntax.Term, ident string) {
	runesStr := getStrOfRunes(ident)

	emitter.printLabel(depth, fmt.Sprintf("_memMngr.vterms[%d] = (struct vterm_t){.tag = V_IDENT_TAG, .str = allocateVStringLiteral((uint32_t[]){%s}, UINT64_C(%d))};",
		emitter.currTermNum, runesStr, utf8.RuneCountInString(ident)))

	term.IndexInLiterals = emitter.currTermNum
	emitter.currTermNum++
}
