package emitter

import (
	"fmt"
	"io"
)

import (
	"syntax"
)

const (
	PARAMS_PTR = "%rsi"
	LOCALS_PTR = "%rdi"
	ARG_LEFT   = "%r8"
	ARG_RIGHT  = "%r9"
	HEAP_PTR   = "%r10"
	CUR_LKTERM = "%r11"
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
}

func (f *Data) Blank() { fmt.Fprint(f, "\n") }

func (f *Data) Comment(s string) { fmt.Fprintf(f, "\t/* %s */\n", s) }

func (f *Data) Header() {
	fmt.Fprintf(f, "\t/*%s*/\n", f.Name)
}

func (f *Data) DotExtern(mangledName string) { fmt.Fprintf(f, "\t.extern\t%s\n", mangledName) }

func (f *Data) DotGlobal(mangledName string) { fmt.Fprintf(f, "\t.globl\t%s\n", mangledName) }

func (f *Data) DotAlign(val int) { fmt.Fprintf(f, "\t.balign\t%d\n", val) }

func (f *Data) Label(name string) { fmt.Fprintf(f, "%s:\n", name) }

func (f *Data) Text(i int) { fmt.Fprintf(f, "\t.text\t%d\n", i) }

type functionStack []*syntax.Function

func (s *functionStack) Empty() bool { return len(*s) == 0 }

func (s *functionStack) Push(fun *syntax.Function) {
	*s = append(*s, fun)
}

func (s *functionStack) Pop() (fun *syntax.Function) {
	ln := len(*s) - 1
	fun = (*s)[ln]
	*s = (*s)[:ln]
	return
}

func decorate(name string) string {
	//TODO: implement real mangling
	return name
}

func processPattern(f Data, p *syntax.Expr) {
}

func processFile(f Data) {
	unit := f.Ast

	f.Header()
	if len(unit.Builtins) > 0 {
		f.Comment("Used built-in functions")
		for name, _ := range unit.Builtins {
			f.DotExtern(decorate(name))
		}
		f.Blank()
	}

	if len(unit.ExtMap) > 0 {
		f.Comment("External functions")
		for name, _ := range unit.ExtMap {
			f.DotExtern(decorate(name))
		}
		f.Blank()
	}

	if len(unit.GlobMap) > 0 {
		f.Comment("Global functions")
		for name, fun := range unit.GlobMap {
			if fun.IsEntry {
				f.DotGlobal(decorate(name))
			}
		}
		f.Blank()
	}

	f.Text(0)
	stack := make(functionStack, 0, len(unit.GlobMap)*2)
	for _, fun := range f.Ast.GlobMap {
		stack.Push(fun)
	}

	for !stack.Empty() {
		fun := stack.Pop()
		f.Label(decorate(fun.FuncName))

		for _, s := range fun.Sentences {
			processPattern(f, &s.Pattern)
			for _, a := range s.Actions {
				switch a.ActionOp {
				case syntax.COMMA: // ','
				case syntax.REPLACE: // '='
				case syntax.TARROW: // '->'
				case syntax.ARROW: // '=>'
				case syntax.COLON: // ':'
				case syntax.DCOLON: // '::'
				}
			}
		}
	}
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
