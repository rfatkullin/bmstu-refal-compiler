package syntax

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/coords"
	"bmstu-refal-compiler/messages"
)

func analyse(ast chan<- *Unit, ms chan<- messages.Data,
	exts <-chan *FuncHeader, globals <-chan *Function, nested <-chan *Function, dialect int) {
	err := func(pos coords.Pos, s string) {
		ms <- messages.Data{pos, messages.ERROR, s}
	}

	warn := func(pos coords.Pos, s string) {
		ms <- messages.Data{pos, messages.WARNING, s}
	}

	errDuplicate := func(pos coords.Pos, what string) {
		err(pos, what+" having the same name is already defined")
	}

	errDuplicateExt := func(pos coords.Pos) {
		errDuplicate(pos, "external function")
	}

	errDuplicateGlobal := func(pos coords.Pos) {
		errDuplicate(pos, "global function")
	}

	errFuncnameInStringForm := func(pos coords.Pos) {
		err(pos, "function name cannot be in string form")
	}

	errEndOfSentenceExpected := func(pos coords.Pos) {
		err(pos, "end of sentence expected")
	}

	errIllegalModifier := func(pos coords.Pos, tag TermTag) {
		err(pos, fmt.Sprintf("illegal use of %v modifier", tag))
	}

	builtins := Builtins[dialect]

	checkFuncName := func(e *FuncHeader) bool {
		if !e.IsIdent && dialect == 5 {
			errFuncnameInStringForm(e.Pos)
		}

		_, ok := builtins[e.FuncName]
		if ok {
			err(e.Pos, "function having the same name is built-in")
		}
		return !ok
	}

	unit := Unit{
		Builtins:        make(map[string]bool, 16),
		ExtMap:          make(map[string]*FuncHeader, 16),
		GlobMap:         make(map[string]*Function, 64),
		FuncsTotalCount: 0,
	}
	ready := make(chan bool)

	go func() {
		for e := range exts {
			if checkFuncName(e) {
				if _, ok := unit.ExtMap[e.FuncName]; ok {
					errDuplicateExt(e.Pos)
				} else {
					unit.ExtMap[e.FuncName] = e
				}
			}
		}

		ready <- true
	}()

	issues := make(map[*coords.Pos]string, 1024)

	var (
		checkPattern func(*Function, *Sentence, *Expr)
		fillScope    func(*Sentence, *Expr)
		checkExpr    func(*Function, *Sentence, *Expr)
		checkActions func(*Function, *Sentence)
		checkBlock   func(*Function, *Scope)
	)

	checkPattern = func(f *Function, s *Sentence, p *Expr) {
		scope := &s.Scope
		for i, t := range p.Terms {
			switch t.TermTag {
			case L, // $L modifier.
				R: // $R modifier.
				if i != 0 {
					errIllegalModifier(t.Start, t.TermTag)
				}
			case VAR: // Variable.
				vt, name := t.VarType, t.Name
				if name != "" {
					if level := scope.FindVar(vt, name); level == -1 {
						scope.AddVar(vt, name)
					} else {
						f.Params.PropagateVar(vt, name, level)
					}
				} else {
					t.Name = scope.AddAnonymousVar(vt)
				}
			case EXPR, // Subexpression in parentheses.
				BRACED_EXPR,    // Subexpression in quoted braces.
				BRACKETED_EXPR, // Subexpression in quoted square brackets.
				ANGLED_EXPR:    // Subexpression in quoted angle brackets.
				checkPattern(f, s, t.Exprs[0])
			case EVAL: // Subexpression inside evaluation brackets.
				err(t.Start, "evaluation brackets are not allowed in pattern")
			case FUNC: // Nested function.
				if dialect == 7 {
					err(t.Start, "nested functions are not allowed in pattern")
				} else {
					err(t.Start, "blocks are not allowed in pattern")
				}
			}
		}
	}

	fillScope = func(s *Sentence, e *Expr) {
		scope := &s.Scope
		for _, t := range e.Terms {
			switch t.TermTag {
			case EXPR, // Subexpression in parentheses.
				BRACED_EXPR,    // Subexpression in quoted braces.
				BRACKETED_EXPR, // Subexpression in quoted square brackets.
				ANGLED_EXPR,    // Subexpression in quoted angle brackets.
				EVAL:           // Subexpression inside evaluation brackets.
				for _, ev := range t.Exprs {
					fillScope(s, ev)
				}
			case FUNC: // Nested function.
				if t.HasName {
					if _, level := scope.FindFunc(t.FuncName); level == -1 {
						scope.AddFunc(t.FuncName, unit.FuncsTotalCount)
						unit.FuncsTotalCount++
					} else {
						errDuplicate(t.Pos, "nested function")
					}
				}
			}
		}
	}

	checkCallableExpr := func(scope *Scope, e *Expr) {
		if e.Len() == 0 {
			if dialect == 5 {
				err(e.Start, "empty evaluation brackets not allowed")
			}
		} else {
			if tn := e.Terms[0]; tn.TermTag == COMP {
				if dialect == 5 && !tn.IsIdent {
					errFuncnameInStringForm(tn.Start)
				} else if _, level := scope.FindFunc(tn.Name); level == -1 {
					issues[&tn.Start] = tn.Name
				}
			} else if dialect == 5 {
				err(tn.Start, "function name is expected")
			} else if tn.TermTag != FUNC && tn.TermTag != VAR &&
				tn.TermTag != EVAL && tn.TermTag != FUNC {
				err(tn.Start, "expression that can be evaluated to function is expected")
			}
		}
	}

	checkExpr = func(f *Function, s *Sentence, e *Expr) {
		scope := &s.Scope
		for _, t := range e.Terms {
			switch t.TermTag {
			case L, // $L modifier.
				R: // $R modifier.
				errIllegalModifier(t.Start, t.TermTag)
			case COMP: // Compound symbol.
				if index, level := scope.FindFunc(t.Name); level != -1 {
					f.Params.PropagateFunc(t.Name, level, index)
				}
			case VAR: // Variable.
				vt, name := t.VarType, t.Name
				if name == "" {
					err(t.Start, "anonymous variables are not allowed here")
				} else if level := scope.FindVar(vt, name); level == -1 {
					err(t.Start, "undefined variable")
				} else {
					f.Params.PropagateVar(vt, name, level)
				}
			case EXPR, // Subexpression in parentheses.
				BRACED_EXPR,    // Subexpression in quoted braces.
				BRACKETED_EXPR, // Subexpression in quoted square brackets.
				ANGLED_EXPR:    // Subexpression in quoted angle brackets.
				checkExpr(f, s, t.Exprs[0])
			case EVAL: // Subexpression inside evaluation brackets.
				checkCallableExpr(scope, t.Exprs[0])
				for _, ev := range t.Exprs {
					checkExpr(f, s, ev)
				}
			case FUNC: // Nested function.
				if dialect == 5 {
					err(t.Start, "blocks are not allowed in expressions")
				}
				t.Params.Parent = &f.Params
				checkBlock(t.Function, scope)
			}
		}
	}

	if dialect == 7 {
		checkActions = func(f *Function, s *Sentence) {
			scope := &s.Scope
			for _, a := range s.Actions {
				switch a.ActionOp {
				case COMMA, // ','
					REPLACE: // '='
					fillScope(s, &a.Expr)
					checkExpr(f, s, &a.Expr)
				case TARROW, // '->'
					ARROW: // '=>'
					fillScope(s, &a.Expr)
					checkCallableExpr(scope, &a.Expr)
					checkExpr(f, s, &a.Expr)
				case COLON, // ':'
					DCOLON: // '::'
					checkPattern(f, s, &a.Expr)
				}
			}
		}
	} else {
		checkActions = func(f *Function, s *Sentence) {
			scope := &s.Scope
			for i, a := range s.Actions {
				switch a.ActionOp {
				case COMMA: // ','
					if i%2 == 1 {
						err(a.Start, "unexpected ','")
					}
					checkExpr(f, s, &a.Expr)
				case REPLACE: // '='
					checkExpr(f, s, &a.Expr)
					if i != s.Len()-1 {
						errEndOfSentenceExpected(a.Follow)
					}
					return
				case COLON: // ':'
					if i%2 == 0 {
						err(a.Start, "unexpected ':'")
					}

					e := &a.Expr
					if e.Len() >= 1 && e.Terms[0].TermTag == FUNC {
						t := e.Terms[0]

						if e.Len() > 1 || i != s.Len()-1 {
							errEndOfSentenceExpected(t.Follow)
						}

						a.ActionOp = ARROW
						t.Params.Parent = &f.Params
						checkBlock(t.Function, scope)
						return
					} else {
						checkPattern(f, s, e)
					}
				}
			}

			err(s.Actions[s.Len()-1].Follow, "sentence without right side")
		}
	}

	checkBlock = func(f *Function, parent *Scope) {
		if f.Len() == 0 {
			err(f.Start, "function must contain at least one sentence")
		} else {
			for _, s := range f.Sentences {
				s.Parent = parent
				checkPattern(f, s, &s.Pattern)
				checkActions(f, s)
			}
		}
	}

	for g := range globals {
		if checkFuncName(&g.FuncHeader) {
			if _, ok := unit.GlobMap[g.FuncName]; ok {
				errDuplicateGlobal(g.Pos)
			} else {
				unit.GlobMap[g.FuncName] = g
				unit.GlobMap[g.FuncName].Index = unit.FuncsTotalCount
				unit.FuncsTotalCount++
			}
		}

		checkBlock(g, nil)
	}

	<-ready

	for pos, name := range issues {
		if se, ok := builtins[name]; ok {
			if _, ok = unit.Builtins[name]; !ok {
				unit.Builtins[name] = se
			}
		} else if _, ok = unit.GlobMap[name]; !ok {
			if _, ok = unit.ExtMap[name]; !ok {
				err(*pos, "undefined function")
			}
		}
	}

	for name, g := range unit.GlobMap {
		if e, ok := unit.ExtMap[name]; ok {
			warn(g.Pos, fmt.Sprintf("hiding external function, defined at %v", e.Pos))
		}
	}

	for currFunc := range nested {
		currFunc.Index = unit.FuncsTotalCount
		unit.FuncsTotalCount++
	}

	ast <- &unit
	close(ast)
}
