package emitter

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

import (
	"bmstu-refal-compiler/syntax"
)

func genTabs(depth int) string {
	return strings.Repeat(tab, depth)
}

func (f *Data) PrintLabel(depth int, label string) {
	tabs := genTabs(depth)
	fmt.Fprintf(f, "%s%s\n", tabs, label)
}

func (f *Data) printFuncHeader(depth int, name string) {
	f.PrintLabel(depth, fmt.Sprintf("struct func_result_t %s(int* entryPoint, struct env_t* env, struct lterm_t* fieldOfView, int entryStatus) \n{", name))
}

func (f *Data) initActionLiterals(depth int, expr syntax.Expr) {

	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	for len(terms) > 0 {

		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.STR:
			f.initStrVTerm(depth, term)
			break

		case syntax.COMP:
			f.initIdentVTerm(depth, term)
			break

		case syntax.INT:
			f.initIntNumVTerm(depth, term)
			break

		case syntax.FLOAT:
			f.initFloatVTerm(depth, term)
			break

		case syntax.EXPR, syntax.EVAL:
			tmpTerms := append(make([]*syntax.Term, 0, len(term.Exprs[0].Terms)+len(terms)), term.Exprs[0].Terms...)
			tmpTerms = append(tmpTerms, terms...)
			terms = tmpTerms
			break

		case syntax.FUNC:
			f.initFuncLiterals(depth, term.Function)
			break
		}
	}
}

func (f *Data) initFuncLiterals(depth int, currFunc *syntax.Function) {

	for _, s := range currFunc.Sentences {
		f.initActionLiterals(depth, s.Pattern)
		for _, a := range s.Actions {
			f.initActionLiterals(depth, a.Expr)
		}
	}
}

func (f *Data) initStrVTerm(depth int, term *syntax.Term) {
	term.IndexInLiterals = f.CurrTermNum

	for i := 0; i < len(term.Value.Str); i++ {
		f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[%d] = (struct v_term){.tag = V_CHAR_TAG, .ch = %d};", f.CurrTermNum, term.Value.Str[i]))
		f.CurrTermNum++
	}
}

// Инициализация vterm_t для литералов целого типа
// Пока только обычные
func (f *Data) initIntNumVTerm(depth int, term *syntax.Term) {
	bytesStr, sign, bytesCount := GetStrOfBytes(term.Value.Int)

	f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[%d] = (struct v_term){.tag = V_INT_NUM_TAG,"+
		" .intNum = allocateIntNumberLiteral((uint8_t[]){%s}, %d, UINT64_C(%d))};",
		f.CurrTermNum, bytesStr, sign, bytesCount))

	term.IndexInLiterals = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для литералов вещественного типа
func (f *Data) initFloatVTerm(depth int, term *syntax.Term) {

	f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[%d] = (struct v_term){.tag = V_DOUBLE_NUM_TAG, .doubleNum = %f};", f.CurrTermNum, term.Value.Float))
	term.IndexInLiterals = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для идентификатора
func (f *Data) initIdentVTerm(depth int, term *syntax.Term) {
	ident := term.Value.Name
	runesStr := GetStrOfRunes(ident)

	f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[%d] = (struct v_term){.tag = V_IDENT_TAG, .str = allocateVStringLiteral((uint32_t[]){%s}, UINT64_C(%d))};",
		f.CurrTermNum, runesStr, utf8.RuneCountInString(ident)))

	term.IndexInLiterals = f.CurrTermNum
	f.CurrTermNum++
}

func (f *Data) printLiteralsAndHeapsInit(depth int, unit *syntax.Unit) {

	f.PrintLabel(depth, "void __initLiteralData()\n{")
	f.PrintLabel(depth+1, "initAllocator(1024 * 1024 * 1024);")

	f.initLiterals(depth+1, unit.GlobMap)

	f.PrintLabel(depth+1, fmt.Sprintf("initHeaps(2, %d);", f.CurrTermNum))

	//fmt.Fprintf(f, "%sdebugLiteralsPrint();\n", tabs)
	f.PrintLabel(depth, "} // __initLiteralData()\n")
}

func (f *Data) initLiterals(depth int, funcs map[string]*syntax.Function) {

	for _, currFunc := range funcs {
		f.initFuncLiterals(depth, currFunc)
	}

	fmt.Fprintf(f, "\n")
}

func (f *Data) PrintHeaders() {

	f.PrintLabel(0, "#include <stdlib.h>")
	f.PrintLabel(0, "#include <stdio.h>\n")
	f.PrintLabel(0, "#include <memory_manager.h>")
	f.PrintLabel(0, "#include <allocators.h>")
	f.PrintLabel(0, "#include <vmachine.h>")
	f.PrintLabel(0, "#include <builtins.h>")
	f.PrintLabel(0, "")
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
