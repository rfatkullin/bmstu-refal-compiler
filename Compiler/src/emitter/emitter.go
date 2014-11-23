package emitter

import (
	"fmt"
	"io"
)

import (
	"strings"
	"syntax"
	"tokens"
)

const (
	tab = " "
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	CurrTermNum int
}

func (f *Data) Comment(s string) { fmt.Fprintf(f, "\t/* %s */\n", s) }

func (f *Data) Header() {
	fmt.Fprintf(f, "// file:%s\n\n", f.Name)
}

func (f *Data) mainFuncHeader() {

	fmt.Fprintf(f, "int main()\n{\n")
}

func (f *Data) mainFuncEnd() {
	fmt.Fprintf(f, "%sreturn 0;\n}\n", tab)
}

func (f *Data) funcHeader(name string) {
	fmt.Fprintf(f, "l_term* %s(vec_header* vecData) \n{\n", name)
	fmt.Fprintf(f, "%sstruct v_term* data = vecData.data;\n", tab)
	fmt.Fprintf(f, "%suint32_t length = vecData.size;\n", tab)
	fmt.Fprintf(f, "%suint32_t ok = 0;\n", tab)
}

func (f *Data) releaseOkVar(tabs string) {
	fmt.Fprintf(f, "%sok = 0;\n", tabs)
}

func (f *Data) setOkVar(tabs string) {
	fmt.Fprintf(f, "%sok = 1;\n", tabs)
}

func (f *Data) checkOKVar(tabs string) {
	fmt.Fprintf(f, "%sif (ok == 1)\n%s{", tabs)
}

func (f *Data) endBlock(tabs string) {
	fmt.Fprintf(f, "%s}\n")
}

func (f *Data) FuncEnd(name string) {
	fmt.Fprintf(f, "} // %s\n\n", name)
}

func genTabs(depth int) string {
	return strings.Repeat(tab, depth)
}

func (f *Data) processSymbol(termNumber, depth int) {

	var startVar = fmt.Sprintf("start_%d", termNumber)
	var followVar = fmt.Sprintf("follow_%d", termNumber)
	var startValue = "0"
	var tabs = genTabs(depth)

	if termNumber != 0 {
		var prevFollow = fmt.Sprintf("follow_%d", termNumber-1)

		fmt.Fprintf(f, "%sif (%s >= length) /*Откат*/;\n", tabs, prevFollow)

		startValue = fmt.Sprintf("follow_%d", termNumber-1)
	}

	fmt.Fprintf(f, "%sint %s = %s = %s;\n", tabs, startVar, followVar, startValue)

	fmt.Fprintf(f, "%sif (data[%s]->tag == V_TERM_SYMBOL_TAG) %s++;\n", tabs, startVar, followVar)
	fmt.Fprintf(f, "%s\telse /*Откат*/;\n", tabs)

	fmt.Fprintf(f, "//--------------------------------------\n")
}

func (f *Data) processExpr(termNumber, nestedDepth int) {

}

func (f *Data) processPattern(p *syntax.Expr) {

	if len(p.Terms) == 0 {
		f.setOkVar(genTabs(1))
	} else {
		for termIndex, term := range p.Terms {

			switch term.TermTag {
			case syntax.VAR:
				switch term.Value.VarType {

				case tokens.VT_S:
					f.processSymbol(termIndex, termIndex+1)
					break

				case tokens.VT_E:
					f.processExpr(termIndex, termIndex+1)
					break
				}

				//fmt.Fprintf(f, "%s ", term.Value.VarType.String())
				break
			}
		}
	}

	fmt.Fprintf(f, "\n")
}

func (f *Data) processAction(act *syntax.Action) {

	f.checkOKVar(genTabs(1))

	f.endBlock(genTabs(1))
}

// Инициализация vterm_t строкового литерала
// Пока только ASCII символы
func (f *Data) initStrVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)
	str := string(term.Value.Str)
	strLen := len(str)

	for i := 0; i < strLen; i++ {
		fmt.Fprintf(f, "%smemoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = %c};\n", tabs, str[i])
	}

	term.Index = f.CurrTermNum
	f.CurrTermNum += strLen
}

// Инициализация vterm_t для литералов целого типа
// Пока только обычные
func (f *Data) initIntNumVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%smemoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_INT_NUM_TAG, .intNum = %d};\n", tabs, term.Value.Int)
	term.Index = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для литералов вещественного типа
func (f *Data) initFloatVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%smemoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_FLOAT_NUM_TAG, .floatNum = %f};\n", tabs, term.Value.Float)
	term.Index = f.CurrTermNum
	f.CurrTermNum++
}

// Инициализация vterm_t для идентификатора
// Пока только ASCII символы
func (f *Data) initIdentVTerm(depth int, term syntax.Term) {
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%smemoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = %s};\n", tabs, string(term.Value.Name))

	term.Index = f.CurrTermNum
	f.CurrTermNum++
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

func (f *Data) initData(depth int) {
	unit := f.Ast
	tabs := genTabs(depth)

	fmt.Fprintf(f, "%smemoryManager.constTermOffset = 0;\n", tabs)
	fmt.Fprintf(f, "%sint i = 0;\n\n", tabs)

	for _, fun := range unit.GlobMap {
		for _, s := range fun.Sentences {
			for _, a := range s.Actions {
				f.initActionData(depth, a.Expr)
			}
		}
	}

	fmt.Fprintf(f, "\n")
}

func processFile(currFileData Data) {
	unit := currFileData.Ast

	currFileData.Header()

	for _, fun := range unit.GlobMap {
		currFileData.funcHeader(fun.FuncName)

		for _, s := range fun.Sentences {
			currFileData.releaseOkVar(genTabs(1))
			currFileData.processPattern(&s.Pattern)
			for _, a := range s.Actions {
				switch a.ActionOp {
				case syntax.COMMA: // ','
				case syntax.REPLACE: // '='
					currFileData.processAction(a)
					break
				case syntax.TARROW: // '->'
				case syntax.ARROW: // '=>'
				case syntax.COLON: // ':'
				case syntax.DCOLON: // '::'
				}
			}
		}

		currFileData.FuncEnd(fun.FuncName)
	}

	currFileData.mainFuncHeader()
	currFileData.initData(1)
	currFileData.mainFuncEnd()
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
