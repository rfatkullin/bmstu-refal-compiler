package syntax

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/coords"
	"bmstu-refal-compiler/tokens"
)

type ActionOp int

const (
	COMMA   ActionOp = iota // ','
	REPLACE                 // '='
	TARROW                  // '->'
	ARROW                   // '=>'
	COLON                   // ':'
	DCOLON                  // '::'
)

var actionNames = map[ActionOp]string{
	COMMA:   "','",
	REPLACE: "'='",
	TARROW:  "'->'",
	ARROW:   "'=>'",
	COLON:   "':'",
	DCOLON:  "'::'",
}

func (a ActionOp) String() string {
	return actionNames[a]
}

var token2action = map[tokens.DomainTag]ActionOp{
	tokens.COMMA:   COMMA,
	tokens.REPLACE: REPLACE,
	tokens.TARROW:  TARROW,
	tokens.ARROW:   ARROW,
	tokens.COLON:   COLON,
	tokens.DCOLON:  DCOLON,
}

type TermTag int

const (
	L              TermTag = iota // $L modifier.
	R                             // $R modifier.
	STR                           // Character string.
	COMP                          // Compound symbol.
	INT                           // Integer number.
	FLOAT                         // Floating-point number.
	VAR                           // Variable.
	EXPR                          // Subexpression in parentheses.
	BRACED_EXPR                   // Subexpression in quoted braces.
	BRACKETED_EXPR                // Subexpression in quoted square brackets.
	ANGLED_EXPR                   // Subexpression in quoted angle brackets.
	EVAL                          // Subexpression inside evaluation brackets.
	FUNC                          // Nested function.
)

var tagnames = map[TermTag]string{
	L:              "$L",
	R:              "$R",
	STR:            "STR",
	COMP:           "COMP",
	INT:            "INT",
	FLOAT:          "FLOAT",
	VAR:            "VAR",
	EXPR:           "EXPR",
	BRACED_EXPR:    "BRACED_EXPR",
	BRACKETED_EXPR: "BRACKETED_EXPR",
	ANGLED_EXPR:    "ANGLED_EXPR",
	EVAL:           "EVAL",
	FUNC:           "FUNC",
}

var paren2term = map[tokens.DomainTag]TermTag{
	tokens.LPAREN:    EXPR,
	tokens.QLBRACE:   BRACED_EXPR,
	tokens.QLBRACKET: BRACKETED_EXPR,
	tokens.QLEVAL:    ANGLED_EXPR,
}

func (tag TermTag) String() string {
	return tagnames[tag]
}

type FuncHeader struct {
	coords.Pos
	HasName  bool
	FuncName string
	IsIdent  bool
	IsSe     bool
}

type Function struct {
	coords.Fragment
	FuncHeader
	IsEntry   bool
	Rollback  bool
	Params    Scope
	Sentences []*Sentence
	Env       map[string]ScopeVar
	Index     int
}

type Sentence struct {
	coords.Fragment
	Scope
	Pattern Expr
	Actions []*Action
}

type Action struct {
	Comment string
	coords.Fragment
	ActionOp
	Expr
}

type Expr struct {
	coords.Fragment
	Terms []*Term
}

type Term struct {
	Comment string
	coords.Fragment
	IndexInLiterals int
	TermTag
	tokens.Value
	Exprs []*Expr
	*Function
}

type ScopeVar struct {
	tokens.VarType
	Number int
}

type Scope struct {
	VarMap              map[string]ScopeVar
	VarsNumber          int
	AnonymousVarsNumber int
	FuncMap             map[string]int

	Parent *Scope
}

type Unit struct {
	Builtins        map[string]bool
	ExtMap          map[string]*FuncHeader
	GlobMap         map[string]*Function
	NestedFuncs     map[string]*Function
	FuncByNumber    map[int]*Function
	FuncsTotalCount int
}

func (f *Function) Len() int { return len(f.Sentences) }

func (f *Function) Add(s *Sentence) {
	f.Sentences = append(f.Sentences, s)
}

func (s *Sentence) Len() int { return len(s.Actions) }

func (s *Sentence) Add(a *Action) {
	s.Actions = append(s.Actions, a)
}

func (e *Expr) Len() int { return len(e.Terms) }

func (e *Expr) Add(t *Term) {
	e.Terms = append(e.Terms, t)
}

func (t *Term) Add(e *Expr) {
	t.Exprs = append(t.Exprs, e)
}

func (s *Scope) AddVar(vt tokens.VarType, n string) {
	if s.VarMap == nil {
		s.VarMap = make(map[string]ScopeVar, 8)
	}

	s.VarMap[n] = ScopeVar{vt, s.VarsNumber}
	s.VarsNumber++
}

func (s *Scope) AddAnonymousVar(vt tokens.VarType) (n string) {
	n = fmt.Sprintf("$%d", s.AnonymousVarsNumber)
	s.AnonymousVarsNumber++
	s.AddVar(vt, n)
	return
}

func (s *Scope) AddFunc(n string, index int) {
	if s.FuncMap == nil {
		s.FuncMap = make(map[string]int, 8)
	}
	s.FuncMap[n] = index
}

func (s *Scope) PropagateVar(vt tokens.VarType, n string, level int) {
	for i := 0; i < level; i++ {
		if _, ok := s.VarMap[n]; ok {
			return
		}
		s.AddVar(vt, n)
		s = s.Parent
	}
}

func (s *Scope) PropagateFunc(n string, level, index int) {
	for i := 0; i < level; i++ {
		if s.FuncMap != nil {
			if _, ok := s.FuncMap[n]; ok {
				return
			}
		}
		s.AddFunc(n, index)
		s = s.Parent
	}
}

func (s *Scope) FindVar(vt tokens.VarType, n string) (level int) {
	for level = 0; s != nil; s, level = s.Parent, level+1 {
		if _, ok := s.VarMap[n]; ok {
			return
		}
	}

	return -1
}

func (s *Scope) FindFunc(n string) (index int, level int) {
	for level := 0; s != nil; s, level = s.Parent, level+1 {
		if s.FuncMap != nil {
			if index, ok := s.FuncMap[n]; ok {
				return index, level
			}
		}
	}

	return 0, -1
}

func (currFunc *Function) setEnv() {
	currFunc.Env = make(map[string]ScopeVar, 0)
	s := &currFunc.Params

	for ; s != nil; s = s.Parent {
		if s.VarMap != nil {
			for varName, varInfo := range s.VarMap {
				currFunc.Env[varName] = ScopeVar{Number: len(currFunc.Env), VarType: varInfo.VarType}
			}
		}
	}
}
