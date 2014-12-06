package emitter

import (
	"fmt"
	//"io"
)

import (
	"strings"
	"syntax"
	//"tokens"
)

const (
	tab = " "
)

func genTabs(depth int) string {
	return strings.Repeat(tab, depth)
}

func (f *Data) Comment(s string) { fmt.Fprintf(f, "\t/* %s */\n", s) }

func (f *Data) PrintLabel(depth int, label string) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%s%s", tabs, label)
}

func (f *Data) funcHeader(name string) {
	fmt.Fprintf(f, "struct func_result_t %s(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) \n{\n", name)
}

func (f *Data) initData(depth int) {
	unit := f.Ast

	for _, fun := range unit.GlobMap {
		for _, s := range fun.Sentences {
			for _, a := range s.Actions {
				f.initActionData(depth, a.Expr)
			}
		}
	}

	fmt.Fprintf(f, "\n")
}

func (f *Data) initActionData(depth int, expr syntax.Expr) {

	terms := make([]*syntax.Term, len(expr.Terms))
	copy(terms, expr.Terms)

	for len(terms) > 0 {

		term := terms[0]
		terms = terms[1:]

		switch term.TermTag {

		case syntax.STR:
			f.initStrVTerm(depth, *term)
			break

		case syntax.COMP:
			f.initIdentVTerm(depth, *term)
			break

		case syntax.INT:
			f.initIntNumVTerm(depth, *term)
			break

		case syntax.FLOAT:
			f.initFloatVTerm(depth, *term)
			break

		case syntax.EXPR:
			terms = append(terms, term.Exprs[0].Terms...)
			break

		case syntax.EVAL:
			//TO DO
			break

		case syntax.FUNC:
			//TO DO
			break

		case syntax.BRACED_EXPR:
		case syntax.BRACKETED_EXPR:
		case syntax.ANGLED_EXPR:
			//Пока считаем, что тут не может быть литералов
			break

		case syntax.VAR:
		case syntax.L:
		case syntax.R:
			//Не литералы
			break
		}
	}
}

// Инициализация vterm_t строкового литерала
// Пока только ASCII символы
func (f *Data) initStrVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)
	str := string(term.Value.Str)
	strLen := len(str)

	for i := 0; i < strLen; i++ {
		fmt.Fprintf(f, "%s*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = '%c'};\n", tabs, str[i])
	}

	term.Index = f.CurrTermNum
	f.CurrTermNum += strLen
}

// Инициализация vterm_t для литералов целого типа
// Пока только обычные
func (f *Data) initIntNumVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%s*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = %d};\n", tabs, term.Value.Int)
	term.Index = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для литералов вещественного типа
func (f *Data) initFloatVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%s*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_FLOAT_NUM_TAG, .floatNum = %f};\n", tabs, term.Value.Float)
	term.Index = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для идентификатора
// Пока только ASCII символы
func (f *Data) initIdentVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%s*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = \"%s\"};\n", tabs, string(term.Value.Name))

	term.Index = f.CurrTermNum
	f.CurrTermNum++
}

func (f *Data) initLiteralDataFunc(depth int) {
	tabs := genTabs(depth + 1)

	fmt.Fprintf(f, "void __initLiteralData()\n{\n")
	fmt.Fprintf(f, "%sinitAllocator(1024 * 1024 * 1024);\n", tabs)
	f.initData(depth + 1)
	fmt.Fprintf(f, "%sinitHeaps(2);\n", tabs)
	fmt.Fprintf(f, "} // __initLiteralData()\n\n")
}

func (f *Data) PrintHeaders() {

	f.PrintLabel(0, "#include <memory_manager.h>\n")
	f.PrintLabel(0, "#include <v_machine.h>\n")
	f.PrintLabel(0, "\n")
}
