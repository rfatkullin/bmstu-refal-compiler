// Bauman Refal Compiler parser and semantic analyzer package
package syntax

import (
	"BMSTU-Refal-Compiler/coords"
	//"encoding/json"
	"BMSTU-Refal-Compiler/messages"
	"BMSTU-Refal-Compiler/tokens"
)

import (
	"fmt"
	"io"
)

func Handle(ast chan<- *Unit, ms chan<- messages.Data, ts <-chan tokens.Data, dialect int) {
	exts := make(chan *FuncHeader, 128)
	globals := make(chan *Function, 128)
	nested := make(chan *Function, 128)
	go analyse(ast, ms, exts, globals, nested, dialect)

	var tok tokens.Data
	var prevPos coords.Pos

	err := func(s string, toPrev bool) {
		switch {
		case tok.Match(tokens.END_OF_INPUT):
			ms <- messages.Data{tok.Start, messages.ERROR, "unexpected end of file"}
		case toPrev:
			ms <- messages.Data{prevPos, messages.ERROR, s}
		default:
			ms <- messages.Data{tok.Start, messages.ERROR, s}
		}
	}

	outsideExpr := true

	next := func() {
		prevPos = tok.Follow

		if tok.Match(tokens.END_OF_INPUT) {
			_, _ = <-ts
			/*if !closed(ts) {
				<-ts
			}*/
		} else {
			if tok = <-ts; tok.Match(tokens.QUOTE) {
				if outsideExpr {
					errPos := tok.Start
					err("quoting is not allowed outside expressions", false)

					tok = <-ts
					for !tok.Match(tokens.QUOTE | tokens.END_OF_INPUT) {
						tok = <-ts
					}

					ms <- messages.Data{errPos, messages.RECOVERY_WARNING,
						fmt.Sprintf("missing input till %v", tok.Follow)}
				} else {
					tok = <-ts
				}
			}
			//fmt.Printf("Reached %v (%s)\n",tok.Follow,tok.Comment)
		}
	}

	errExpected := func(mask tokens.DomainTag, toPrev bool) {
		err(mask.ExpectedString(), toPrev)
	}

	errFuncnameExpected := func(toPrev bool) {
		err("function name expected", toPrev)
	}

	missTill := func(mask tokens.DomainTag) {
		start, flag := tok.Start, false
		for !tok.Match(mask | tokens.END_OF_INPUT) {
			flag = true
			next()
		}

		if flag {
			ms <- messages.Data{start, messages.RECOVERY_WARNING,
				fmt.Sprintf("missing input till %v", prevPos)}
		}
	}

	expect := func(mask tokens.DomainTag, miss tokens.DomainTag) {
		if !tok.Match(mask) {
			errExpected(mask, false)
			missTill(mask | miss)
		}
	}

	var (
		declaration func()
		declTail    func(*Function)
		nameList    func(bool)
		semicolons  func()
		block       func(*Function)
		repSentence func(*Function)
		sentence    func(tokens.DomainTag) *Sentence
		repAction   func(*Sentence, bool)
		action      func() *Action
		expr        func(*Expr, tokens.DomainTag)
		term        func() *Term
	)

	// program ::= declaration program
	//           | ε.
	// FIRST(program)  = { ε, EXTERN, SE, ENTRY, COMPOUND }
	// FOLLOW(program) = { END_OF_INPUT }
	// Always advances input: ?
	program := func() {
		firsts := tokens.EXTERN | tokens.ENTRY | tokens.COMPOUND
		if dialect == 7 {
			firsts |= tokens.SE
		}

		for {
			if !tok.Match(tokens.END_OF_INPUT) {
				expect(firsts, 0)
			}

			if tok.Match(tokens.END_OF_INPUT) {
				break
			}

			declaration() // declaration's FIRST set is ensured here
		}

		next()
	}

	// declaration ::= mod decl-tail
	//               | EXTERN ext-tail.
	// mod         ::= SE mod
	//               | ENTRY mod
	//               | ε.
	// ext-tail    ::= ext-tail2
	//               | SE ext-tail2.
	// ext-tail2   ::= COMPOUND name-list semicolons.
	// FIRST(declaration)  = { EXTERN, SE, ENTRY, COMPOUND }
	// FOLLOW(declaration) = { END_OF_INPUT, EXTERN, SE, ENTRY, COMPOUND }
	// Always advances input: yes, if first set ensured.
	declaration = func() {
		start := tok.Start

		if tok.Match(tokens.EXTERN) {
			next()

			isSe := false
			if tok.Match(tokens.SE) {
				next()
				isSe = true
			}

			if tok.Match(tokens.COMPOUND) {
				exts <- &FuncHeader{tok.Start, true, tok.Name, tok.IsIdent, isSe}
				next()
				nameList(isSe)
			} else {
				errFuncnameExpected(false)
				if tok.Match(tokens.COMMA) {
					nameList(isSe)
				}
			}

			if tok.Match(tokens.SEMICOLON) {
				semicolons()
			}
		} else {
			f := new(Function)
			f.Start = start

			for {
				if tok.Match(tokens.SE) {
					if f.IsSe {
						err("extra $SE keyword", false)
					} else {
						f.IsSe = true
					}
				} else if tok.Match(tokens.ENTRY) {
					if f.IsEntry {
						err("extra $ENTRY keyword", false)
					} else {
						f.IsEntry = true
					}

				} else {
					break
				}

				next()
			}

			declTail(f)
			f.Follow = prevPos
			globals <- f
		}
	}

	// decl-tail      ::= COMPOUND body.
	// body           ::= block opt-semicolons
	//                  | expr-f action rep-action semicolons.
	// opt-semicolons ::= semicolons
	//                  | ε.
	// FIRST(decl-tail)  = { COMPOUND }
	// FOLLOW(decl-tail) = { END_OF_INPUT, EXTERN, SE, ENTRY, COMPOUND }
	// Always advances input: ?
	declTail = func(f *Function) {
		outsideExpr = false

		if tok.Match(tokens.COMPOUND) {
			f.HasName = true
			f.FuncName = tok.Name
			f.IsIdent = tok.IsIdent
			f.Pos = tok.Start
			next()
		} else {
			errFuncnameExpected(false)
		}

		if tok.Match(tokens.LBRACE | tokens.QMARK) {
			block(f) // block's FIRST set is ensured here
			if tok.Match(tokens.SEMICOLON) {
				semicolons()
			}
		} else if dialect == 7 {
			f.Add(sentence(0))

			if tok.Match(tokens.SEMICOLON) {
				semicolons()
			} else {
				errExpected(tokens.SEMICOLON, true)
			}
		} else {
			errExpected(tokens.LBRACE, true)
		}

		outsideExpr = true
	}

	// name-list ::= COMMA COMPOUND name-list
	//             | ε.
	// FIRST(name-list)  = { ε, COMMA }
	// FOLLOW(name-list) = { SEMICOLON }
	// Always advances input: no
	nameList = func(isSe bool) {
		for !tok.Match(tokens.SEMICOLON) {
			if tok.Match(tokens.COMMA) {
				next()

				for tok.Match(tokens.COMMA) {
					err("extra ','", false)
					next()
				}

				if tok.Match(tokens.SE) && tok.Start.Line == prevPos.Line {
					if isSe {
						err("extra $SE keyword", false)
					} else {
						err("$SE must follow $EXTERN", false)
					}
					next()
				}
			} else {
				if !(tok.Match(tokens.COMPOUND) && tok.Start.Line == prevPos.Line) {
					// Probably missed ';'. Exit.
					errExpected(tokens.SEMICOLON, true)
					return
				}

				// Maybe we suddenly missed ','?
				errExpected(tokens.COMMA, true)
			}

			if tok.Match(tokens.COMPOUND) {
				exts <- &FuncHeader{tok.Start, true, tok.Name, tok.IsIdent, isSe}
				next()
			} else if tok.Match(tokens.SEMICOLON) {
				// Unexpected end of current declaration. Exit.
				errFuncnameExpected(true)
				return
			} else if tok.Match(tokens.EXTERN | tokens.ENTRY | tokens.SE) {
				// Maybe beginning of the next declaration? Exit.
				errExpected(tokens.SEMICOLON, true)
				return
			} else {
				// Absolute garbage detected. Exit.
				err("function name, or ';' expected", true)
				return
			}
		}
	}

	// semicolons ::= SEMICOLON.
	// FIRST(semicolons)  = { SEMICOLON }
	// FOLLOW(semicolons) = { END_OF_INPUT, EXTERN, SE, ENTRY, COMPOUND, RBRACE, LBRACE, QMARK,
	//                        COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW, L, R, STRING,
	//                        INTEGER, FLOAT, VAR, LPAREN, QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC }
	// Always advances input: yes
	semicolons = func() {
		next()
		for tok.Match(tokens.SEMICOLON) {
			err("extra semicolon", false)
			next()
		}
	}

	// block ::= LBRACE rep-sentence
	//         | QMARK LBRACE rep-sentence.
	// FIRST(block)  = { LBRACE, QMARK }
	// FOLLOW(block) = { END_OF_INPUT, EXTERN, SE, ENTRY, COMPOUND, SEMICOLON, COMMA, REPLACE,
	//                   COLON, DCOLON, TARROW, ARROW, RBRACE, LBRACE, QMARK, L, R, STRING,
	//                   INTEGER, FLOAT, VAR, LPAREN, QLBRACE, QLBRACKET, QLEVAL,
	//                   LEVAL, FUNC, RPAREN, TILDE, REVAL }
	// Always advances input: yes, if FIRST set ensured
	block = func(f *Function) {
		if tok.Match(tokens.QMARK) {
			f.Rollback = true
			next()
		}

		if tok.Match(tokens.LBRACE) {
			next()
		} else {
			errExpected(tokens.LBRACE, true)
		}

		repSentence(f)
	}

	// rep-sentence ::= sentence tl-rep-sentence
	//                | RBRACE.
	// tl-rep-sentence ::= semicolons rep-sentence
	//                   | RBRACE.
	// FIRST(rep-sentence)  = { RBRACE, LBRACE, QMARK, COMMA, REPLACE, COLON, DCOLON, TARROW,
	//                          ARROW, L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR, LPAREN,
	//                          QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC }
	// FOLLOW(rep-sentence) = { END_OF_INPUT, EXTERN, SE, ENTRY, COMPOUND, SEMICOLON, COMMA,
	//                          REPLACE, COLON, DCOLON, TARROW, ARROW, RBRACE, LBRACE,
	//                          QMARK, L, R, STRING, INTEGER, FLOAT, VAR, LPAREN,
	//                          QLBRACE, QLBRACKET, QLEVAL, LEVAL,
	//                          FUNC, RPAREN, TILDE, REVAL }
	// Always advances input: no
	repSentence = func(f *Function) {
		beginning := tokens.L | tokens.R | tokens.STRING | tokens.COMPOUND | tokens.INTEGER |
			tokens.FLOAT | tokens.VAR | tokens.LPAREN | tokens.QLBRACE | tokens.QLBRACKET |
			tokens.QLEVAL | tokens.LEVAL | tokens.FUNC |
			tokens.COMMA | tokens.REPLACE | tokens.COLON |
			tokens.DCOLON | tokens.TARROW | tokens.ARROW

		for !tok.Match(tokens.RBRACE) {
			if !tok.Match(beginning) {
				if !tok.Match(tokens.EXTERN | tokens.SE | tokens.ENTRY) {
					err("pattern expected", false)
				}
				missTill(beginning | tokens.EXTERN | tokens.SE | tokens.ENTRY |
					tokens.SEMICOLON | tokens.RBRACE)
			}

			if tok.Match(beginning) {
				f.Add(sentence(tokens.RBRACE)) // sentence's FIRST set is ensured here
			} else {
				break
			}

			if tok.Match(tokens.SEMICOLON) {
				semicolons()
			} else {
				break
			}
		}

		if !tok.Match(tokens.RBRACE) {
			errExpected(tokens.RBRACE, true)
		} else {
			next()
		}
	}

	// sentence ::= expr-f action rep-action.
	// FIRST(sentence)  = { L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR, LPAREN,
	//                      QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC,
	//                      COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW }
	// FOLLOW(sentence) = { SEMICOLON, RBRACE }
	// Always advances input: yes, if FIRST set ensured.
	sentence = func(terminator tokens.DomainTag) *Sentence {
		s := new(Sentence)
		s.Start = tok.Start
		expr(&s.Pattern, 0)

		operations := tokens.COMMA | tokens.REPLACE | tokens.COLON |
			tokens.DCOLON | tokens.TARROW | tokens.ARROW
		if !tok.Match(operations) {
			if dialect == 5 {
				err("right side or condition expected", true)
			} else {
				err("action expected", true)
			}
			missTill(operations | tokens.EXTERN | tokens.SE | tokens.ENTRY |
				tokens.SEMICOLON | terminator)
		}

		if tok.Match(operations) {
			s.Add(action()) // action's FIRST set is ensured here
			repAction(s, true)
		}

		s.Follow = prevPos
		return s
	}

	// rep-action ::= action rep-action
	//              | ε.
	// FIRST(rep-action)  = { ε, COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW }
	// FOLLOW(rep-action) = { SEMICOLON, RBRACE }
	// Always advances input: no.
	repAction = func(s *Sentence, rbrace bool) {
		operations := tokens.COMMA | tokens.REPLACE | tokens.COLON |
			tokens.DCOLON | tokens.TARROW | tokens.ARROW
		follow := tokens.SEMICOLON
		if rbrace {
			follow |= tokens.RBRACE
		}

		for {
			if !tok.Match(operations | follow) {
				if dialect == 5 {
					err("right side or condition expected", false)
				} else {
					err("action expected", false)
				}
				missTill(operations | follow | tokens.EXTERN | tokens.SE | tokens.ENTRY)
			}

			if !tok.Match(operations) {
				break
			}

			s.Add(action()) // action's FIRST set is ensured here
		}
	}

	// action ::= operation expr.
	// operation ::= COMMA | REPLACE | COLON | DCOLON | TARROW | ARROW.
	// FIRST(action)  = { COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW }
	// FOLLOW(action) = { COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW, SEMICOLON, RBRACE }
	// Always advances input: yes, if FIRST set ensured.
	action = func() *Action {
		a := new(Action)
		a.Start = tok.Start

		a.ActionOp = token2action[tok.DomainTag]
		a.Comment = a.ActionOp.String()
		next()
		expr(&a.Expr, 0)

		a.Follow = prevPos
		return a
	}

	// expr-f ::= term-f expr
	//          | ε.
	// expr ::= term expr
	//        | ε.
	// FIRST(expr-f)  = { ε, L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR,
	//                    LPAREN, QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC }
	// FOLLOW(expr-f) = { COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW }
	// FIRST(expr)  = { ε, LBRACE, QMARK, L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR,
	//                  LPAREN, QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC }
	// FOLLOW(expr) = { COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW, SEMICOLON, RBRACE,
	//                  RPAREN, TILDE, REVAL }
	// Always advances input: no
	expr = func(e *Expr, terminator tokens.DomainTag) {
		e.Start = tok.Start

		count := 0
	L:
		for tok.Match(tokens.LBRACE | tokens.QMARK | tokens.L | tokens.R | tokens.STRING |
			tokens.COMPOUND | tokens.INTEGER | tokens.FLOAT | tokens.VAR | tokens.LPAREN |
			tokens.QLBRACE | tokens.QLBRACKET | tokens.QLEVAL | tokens.LEVAL | tokens.FUNC |
			/* added for error recovery: */
			tokens.RPAREN | tokens.REVAL | tokens.QRBRACE |
			tokens.QRBRACKET | tokens.QREVAL) {
			switch {
			case tok.Match(terminator):
				break L
			case tok.Match(tokens.RPAREN):
				err("no matching '('", false)
				next()
			case tok.Match(tokens.REVAL):
				err("no matching '<'", false)
				next()
			case tok.Match(tokens.QRBRACE):
				err("no matching quoted '{'", false)
				next()
			case tok.Match(tokens.QRBRACKET):
				err("no matching quoted '['", false)
				next()
			case tok.Match(tokens.QREVAL):
				err("no matching quoted '<'", false)
				next()
			default:
				e.Add(term()) // term's FIRST set is ensured here
				count++
			}
		}

		if count == 0 {
			e.Follow = e.Start
		} else {
			e.Follow = prevPos
		}
	}

	// term     ::= term-f
	//            | block.
	// term-f   ::= L
	//            | R
	//            | STRING
	//            | COMPOUND
	//            | INTEGER
	//            | FLOAT
	//            | VAR
	//            | LPAREN expr RPAREN
	//            | QLBRACE expr QRBRACE
	//            | QLBRACKET expr QRBRACKET
	//            | QLEVAL expr QREVAL
	//            | LEVAL expr rep-expr REVAL
	//            | FUNC COMPOUND block.
	// rep-expr ::= TILDE expr rep-expr
	//            | ε.
	// FIRST(term)  = { LBRACE, QMARK, L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR, LPAREN,
	//                  QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC }
	// FOLLOW(term) = { COMMA, REPLACE, COLON, DCOLON, TARROW, ARROW, SEMICOLON, RBRACE,
	//                  LBRACE, QMARK, L, R, STRING, COMPOUND, INTEGER, FLOAT, VAR, LPAREN,
	//                  QLBRACE, QLBRACKET, QLEVAL, LEVAL, FUNC, RPAREN, TILDE, REVAL }
	// Always advances input: yes, if FIRST set ensured.
	term = func() *Term {
		t := new(Term)
		t.Start = tok.Start

		switch tok.DomainTag {
		case tokens.L:
			t.TermTag = L
			t.Comment = L.String()
			next()
		case tokens.R:
			t.TermTag = R
			t.Comment = R.String()
			next()
		case tokens.STRING:
			t.TermTag = STR
			t.Comment = fmt.Sprintf("%v '%s'", STR, string(tok.Str))
			t.Value = tok.Value
			next()
		case tokens.COMPOUND:
			t.TermTag = COMP
			t.Comment = fmt.Sprintf("%v '%s'", COMP, tok.Name)
			t.Value = tok.Value
			next()
		case tokens.INTEGER:
			t.TermTag = INT
			t.Comment = fmt.Sprintf("%v %v", INT, tok.Int)
			t.Value = tok.Value
			next()
		case tokens.FLOAT:
			t.TermTag = FLOAT
			t.Comment = fmt.Sprintf("%v %f", FLOAT, tok.Float)
			t.Value = tok.Value
			next()
		case tokens.VAR:
			t.TermTag = VAR
			t.Comment = fmt.Sprintf("%v %v.%s", VAR, tok.VarType, tok.Name)
			t.Value = tok.Value
			next()
		case tokens.LPAREN,
			tokens.QLBRACE,
			tokens.QLBRACKET,
			tokens.QLEVAL:
			t.TermTag = paren2term[tok.DomainTag]
			t.Comment = t.TermTag.String()
			right := tok.DomainTag.Pair()
			next()

			e := new(Expr)
			expr(e, right)
			t.Add(e)

			if tok.Match(right) {
				next()
			} else {
				err(fmt.Sprintf("%v expected", right), true)
			}
		case tokens.LEVAL:
			t.TermTag = EVAL
			t.Comment = EVAL.String()
			next()

			e := new(Expr)
			expr(e, tokens.REVAL)
			t.Add(e)

			for tok.Match(tokens.TILDE) {
				next()
				e := new(Expr)
				expr(e, tokens.REVAL)
				t.Add(e)
			}

			if tok.Match(tokens.REVAL) {
				next()
			} else {
				errExpected(tokens.REVAL, true)
			}
		case tokens.FUNC:
			t.TermTag = FUNC
			t.Comment = FUNC.String()
			next()

			f := new(Function)
			t.Function = f
			f.Start = tok.Start

			if tok.Match(tokens.COMPOUND) {
				f.HasName = true
				f.FuncName = tok.Name
				f.IsIdent = tok.IsIdent
				f.Pos = tok.Start
				next()
			} else {
				errFuncnameExpected(false)
			}

			if tok.Match(tokens.QMARK | tokens.LBRACE) {
				block(f) // block's FIRST set is ensured here
				f.Follow = prevPos
			} else {
				errExpected(tokens.LBRACE, true)
				f.Follow = tok.Start
			}
			nested <- f
		case tokens.QMARK,
			tokens.LBRACE:
			t.TermTag = FUNC
			t.Comment = FUNC.String()

			f := new(Function)
			t.Function = f
			f.Start = tok.Start
			block(f) // block's FIRST set is ensured here
			f.Follow = prevPos
			nested <- f
		}

		t.Follow = prevPos
		return t
	}

	next()
	program()
	close(exts)
	close(globals)
	close(nested)
}

func Intercept(out chan<- *Unit, w io.WriteCloser, in <-chan *Unit) {
	tree := <-in
	out <- tree
	close(out)
	//bs, _ := json.MarshalIndent(tree, "", "\t")
	w.Write([]byte(PrintUnit(tree)))
	//w.Write(bs)
	w.Close()
}
