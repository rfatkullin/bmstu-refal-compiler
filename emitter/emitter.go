package emitter

import (
	"fmt"
	"io"
)

import (
	fk "BMSTU-Refal-Compiler/emitter/funcs_keeper"
	sk "BMSTU-Refal-Compiler/emitter/scope_keeper"
	"BMSTU-Refal-Compiler/syntax"
)

type funcData struct {
	*syntax.Function
	emittedName string
}

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	CurrTermNum int
}

type patternContext struct {
	entryPoint     int
	prevEntryPoint int
}

type sentenceInfo struct {
	patternsCount  int
	index          int
	isFirstPattern bool
	isLastPattern  bool
	isLast         bool
	scope          *syntax.Scope
}

type emitterContext struct {
	entryPoint             int
	prevEntryPoint         int
	patternNumber          int
	maxPatternNumber       int
	nextSentenceEntryPoint int
	isNextActMatching      bool
	isLastAction           bool
	isFuncCallInConstruct  bool
	sentenceInfo           sentenceInfo
	fixedVars              map[string]int
	patternCtx             patternContext
	scopeKeeper            *sk.ScopeKeeper
	funcsKeeper            *fk.FuncsKeeper
	nestedNamedFuncs       []*fk.FuncInfo
	currFuncInfo           *fk.FuncInfo
}

func (f *Data) mainFunc(depth int, entryFuncName string) {

	f.PrintLabel(depth, "int main()")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "__initLiteralData();")
	f.PrintLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	f.PrintLabel(depth+1, "return 0;")
	f.PrintLabel(depth, "}")
}

func (f *Data) printInitLocals(depth, maxPatternNumber, varsNumber int) {

	f.PrintLabel(depth, "struct func_result_t funcRes;")
	f.PrintLabel(depth, "struct lterm_t* funcCallChain = 0;")
	f.PrintLabel(depth, "struct fragment_t* currFrag = 0;")
	f.PrintLabel(depth, "struct lterm_t** helper = 0;")
	f.PrintLabel(depth, "struct lterm_t* currTerm = 0;")
	f.PrintLabel(depth, "struct lterm_t* funcTerm = 0;")
	f.PrintLabel(depth, "int fragmentOffset = 0;")
	f.PrintLabel(depth, "int stretchingVarNumber = 0;")
	f.PrintLabel(depth, "int stretching = 0;")
	f.PrintLabel(depth, "int i = 0;")
	f.PrintLabel(depth, "int j = 0;")
	f.PrintLabel(depth, "if (*entryPoint == 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("env->locals = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));", maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->assembledFOVs = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));", maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->_FOVs = (struct lterm_t**)malloc(%d * sizeof(struct lterm_t*));", maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("env->stretchVarsNumber = (int*)malloc(%d * sizeof(int));", maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("for (i = 0; i < %d; i++)", maxPatternNumber))
	f.PrintLabel(depth+1, "{")
	f.PrintLabel(depth+2, "env->_FOVs[i] = 0;")
	f.PrintLabel(depth+2, "env->assembledFOVs[i] = 0;")
	f.PrintLabel(depth+2, fmt.Sprintf("env->locals[i] = (struct lterm_t*)malloc(%d * sizeof(struct lterm_t));", varsNumber))
	f.PrintLabel(depth+2, fmt.Sprintf("for (j = 0; j < %d; j++)", varsNumber))
	f.PrintLabel(depth+2, "{")
	f.PrintLabel(depth+3, "env->locals[i][j].tag = L_TERM_FRAGMENT_TAG;")
	f.PrintLabel(depth+3, "env->locals[i][j].fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));")
	f.PrintLabel(depth+2, "}")
	f.PrintLabel(depth+1, "}")
	f.initSretchVarNumbers(depth+1, maxPatternNumber)
	f.PrintLabel(depth, "}")
}

func (f *Data) printFreeLocals(depth, matchingNumber, varsNumber int) {

	f.PrintLabel(depth, "if (funcRes.status != CALL_RESULT)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "free(env->locals);")
	f.PrintLabel(depth+1, "free(env->stretchVarsNumber);")
	f.PrintLabel(depth+1, "free(env->assembledFOVs);")
	f.PrintLabel(depth, "}")
}

func (f *Data) calcMaxPatternsAndVarsNumbers(currFunc *syntax.Function) (int, int) {
	maxPatternsNumber := 0
	maxVarNumber := 0

	for _, s := range currFunc.Sentences {
		maxPatternsNumber = max(maxPatternsNumber, f.calcPatternsNumber(s))
		maxVarNumber = max(maxVarNumber, s.Scope.VarsNumber)
	}

	return maxPatternsNumber, maxVarNumber
}

func (f *Data) calcPatternsNumber(s *syntax.Sentence) int {
	// +1 s.Pattern
	number := 1

	for _, a := range s.Actions {
		if a.ActionOp == syntax.COLON {
			number++
		}
	}

	return number
}

func (f *Data) processFuncSentences(depth int, funcInfo *fk.FuncInfo, ctx *emitterContext) {

	maxVarsNumber := 0
	ctx.entryPoint = 0
	ctx.patternNumber = 0
	ctx.maxPatternNumber, maxVarsNumber = f.calcMaxPatternsAndVarsNumbers(funcInfo.Function)
	ctx.scopeKeeper = funcInfo.ScopeKeeper
	ctx.scopeKeeper.AddFuncScope(funcInfo.FuncName)
	ctx.currFuncInfo = funcInfo

	f.printInitLocals(depth, ctx.maxPatternNumber, maxVarsNumber)
	f.PrintLabel(depth, "while(*entryPoint >= 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (*entryPoint)")
	f.PrintLabel(depth+1, "{")

	for sentenceNumber, s := range funcInfo.Function.Sentences {
		ctx.scopeKeeper.AddSentenceScope(sentenceNumber)

		ctx.fixedVars = make(map[string]int)

		ctx.sentenceInfo.index = sentenceNumber
		ctx.sentenceInfo.scope = &s.Scope
		ctx.sentenceInfo.patternsCount = f.calcPatternsNumber(s)
		ctx.sentenceInfo.isLast = sentenceNumber == len(funcInfo.Function.Sentences)-1
		ctx.sentenceInfo.isFirstPattern = true
		ctx.sentenceInfo.isLastPattern = ctx.sentenceInfo.patternsCount == 1

		ctx.nextSentenceEntryPoint = ctx.entryPoint + ctx.sentenceInfo.patternsCount
		ctx.patternNumber = 0
		ctx.isFuncCallInConstruct = false
		ctx.prevEntryPoint = -1

		f.matchingPattern(depth+1, ctx, s.Pattern.Terms)

		ctx.sentenceInfo.isFirstPattern = false

		for index, a := range s.Actions {

			ctx.isLastAction = index == len(s.Actions)-1
			ctx.isNextActMatching = false
			if index+1 < len(s.Actions) && (s.Actions[index+1].ActionOp == syntax.COLON || s.Actions[index+1].ActionOp == syntax.DCOLON) {
				ctx.isNextActMatching = true
			}

			switch a.ActionOp {

			case syntax.COMMA: // ','
				f.ConstructResult(depth+2, ctx, &s.Scope, a.Expr)
				break

			case syntax.COLON: // ':'
				f.PrintLabel(depth+1, "} // Pattern case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.REPLACE: // '='
				ctx.prevEntryPoint = -1
				f.ConstructResult(depth+2, ctx, &s.Scope, a.Expr)
				break

			case syntax.DCOLON: // '::'
				ctx.prevEntryPoint = -1
				f.PrintLabel(depth+1, "} // Pattern case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'

			}
		}

		f.PrintLabel(depth+2, "*entryPoint = -1;")
		f.PrintLabel(depth+2, "break; //Successful end of sentence")
		f.PrintLabel(depth+1, "} // Pattern case end")

		ctx.scopeKeeper.PopLastSentenceScope()
	}

	f.PrintLabel(depth+1, "} // Entry point switch end")
	f.PrintLabel(depth, "} // Main while end")

	f.printFreeLocals(depth, ctx.maxPatternNumber, maxVarsNumber)
	f.PrintLabel(depth, "return funcRes;")
}

func (f *Data) predeclareGlobFuncs(depth, startIndex, funcsNumber int) {

	for i := 0; i < funcsNumber; i++ {
		f.PrintLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int* entryPoint, struct env_t* env, struct lterm_t* fieldOfView);", startIndex+i))
	}

	f.PrintLabel(depth, "")
}

func (f *Data) processGlobFuncs(depth int, ctx *emitterContext, funcs map[string]*syntax.Function) (entryFuncName string) {

	globs := make([]*fk.FuncInfo, 0)

	for funcName, currFunc := range funcs {

		funcInfo := ctx.funcsKeeper.AddFunc(sk.NewScopeKeeper(), currFunc)

		globs = append(globs, funcInfo)

		if funcName == "Go" {
			entryFuncName = funcInfo.EmittedFuncName
		}
	}

	for _, funcInfo := range globs {
		f.printFuncHeader(depth, funcInfo.EmittedFuncName)
		f.processFuncSentences(depth+1, funcInfo, ctx)
		f.PrintLabel(depth, fmt.Sprintf("} // func %s:%s\n", funcInfo.FuncName, funcInfo.EmittedFuncName)) // func block end
	}

	return
}

func (f *Data) processNestedFuncs(depth int, ctx *emitterContext) {

	for _, funcInfo := range ctx.nestedNamedFuncs {
		f.printFuncHeader(depth, funcInfo.EmittedFuncName)
		f.processFuncSentences(depth+1, funcInfo, ctx)

		funcName := "anonym"

		if funcInfo.Function.HasName {
			funcName = funcInfo.FuncName
		}

		f.PrintLabel(depth, fmt.Sprintf("} // func %s:%s\n", funcName, funcInfo.EmittedFuncName)) // func block end
	}
}

func (f *Data) addBuiltins(ctx *emitterContext) int {
	ctx.funcsKeeper.AddBuiltinFunc("Prout")
	ctx.funcsKeeper.AddBuiltinFunc("Card")

	return 2
}

func processFile(f Data) {
	unit := f.Ast
	depth := 0

	var ctx emitterContext
	ctx.funcsKeeper = fk.NewFuncsKeeper()
	ctx.nestedNamedFuncs = make([]*fk.FuncInfo, 0)

	f.PrintLabel(depth, fmt.Sprintf("// file:%s\n", f.Name))
	f.PrintHeaders()

	builtinsNumber := f.addBuiltins(&ctx)
	f.predeclareGlobFuncs(depth, builtinsNumber, len(unit.GlobMap)+unit.AnonymousNumber)

	f.printLiteralsAndHeapsInit(depth, unit)

	entryFuncName := f.processGlobFuncs(depth, &ctx, unit.GlobMap)

	f.processNestedFuncs(depth, &ctx)

	f.mainFunc(depth, entryFuncName)

	//ctx.funcsKeeper.PrintAllFuncs()
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
