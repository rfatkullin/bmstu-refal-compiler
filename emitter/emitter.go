package emitter

import (
	"fmt"
	"io"
)

import (
	"BMSTU-Refal-Compiler/syntax"
)

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

type emitterContext struct {
	entryPoint               int
	patternNumber            int
	inSentencePatternNumber  int
	maxPatternNumber         int
	nextSentenceEntryPoint   int
	sentenceNumber           int
	isFirstPatternInSentence bool
	isLastPatternInSentence  bool
	isLastSentence           bool
	isNextActMatching        bool
	isLastAction             bool
	isFuncCallInConstruct    bool
	sentenceScope            *syntax.Scope
	fixedVars                map[string]bool
	patternCtx               patternContext
}

func (f *Data) mainFunc(depth int) {

	f.PrintLabel(depth, "int main()")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "__initLiteralData();")
	f.PrintLabel(depth+1, "mainLoop(Go);")
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
	f.PrintLabel(depth+1, fmt.Sprintf("env->stretchVarsNumber = (int*)malloc(%d * sizeof(int));", maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("for (i = 0; i < %d; i++)", maxPatternNumber))
	f.PrintLabel(depth+1, "{")
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

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	//isTheresPatternsExists := f.isTherePatternsExists(currFunc.Sentences)
	var ctx emitterContext
	maxVarsNumber := 0
	ctx.entryPoint = 0
	ctx.patternNumber = 0
	ctx.maxPatternNumber, maxVarsNumber = f.calcMaxPatternsAndVarsNumbers(currFunc)

	f.printInitLocals(depth, ctx.maxPatternNumber, maxVarsNumber)

	//if isTheresPatternsExists {
	f.PrintLabel(depth, "while(*entryPoint >= 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (*entryPoint)")
	f.PrintLabel(depth+1, "{")
	//}

	for sentenceNumber, s := range currFunc.Sentences {

		ctx.fixedVars = make(map[string]bool)
		ctx.sentenceNumber = sentenceNumber
		ctx.sentenceScope = &s.Scope
		ctx.inSentencePatternNumber = f.calcPatternsNumber(s)
		ctx.isLastSentence = sentenceNumber == len(currFunc.Sentences)-1
		ctx.isFirstPatternInSentence = true
		ctx.isLastPatternInSentence = ctx.inSentencePatternNumber == 1
		ctx.nextSentenceEntryPoint = ctx.entryPoint + ctx.inSentencePatternNumber
		ctx.patternNumber = 0
		ctx.isFuncCallInConstruct = false

		f.matchingPattern(depth+1, &ctx, s.Pattern.Terms)

		ctx.isFirstPatternInSentence = false

		for index, a := range s.Actions {

			ctx.isLastAction = index == len(s.Actions)-1
			ctx.isNextActMatching = false
			if index+1 < len(s.Actions) && (s.Actions[index+1].ActionOp == syntax.COLON || s.Actions[index+1].ActionOp == syntax.DCOLON) {
				ctx.isNextActMatching = true
			}

			switch a.ActionOp {

			case syntax.REPLACE, syntax.COMMA: // '=' ','
				f.ConstructResult(depth+2, &ctx, &s.Scope, a.Expr)
				break

			case syntax.COLON: // ':'
				f.PrintLabel(depth+1, "} // Pattern case end\n")
				f.matchingPattern(depth+1, &ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
			case syntax.ARROW: // '=>'
			case syntax.DCOLON: // '::'
			}
		}

		f.PrintLabel(depth+2, "*entryPoint = -1;")
		f.PrintLabel(depth+2, "break; //Successful end of sentence")
		f.PrintLabel(depth+1, "} // Pattern case end")
	}

	f.PrintLabel(depth+1, "} // Entry point switch end")
	f.PrintLabel(depth, "} // Main while end")

	f.printFreeLocals(depth, ctx.maxPatternNumber, maxVarsNumber)
	f.PrintLabel(depth, "return funcRes;")

}

func processFile(f Data) {
	unit := f.Ast
	depth := 0

	f.PrintLabel(depth, fmt.Sprintf("// file:%s\n", f.Name))
	f.PrintHeaders()

	f.initLiteralDataFunc(depth)

	for _, currFunc := range unit.GlobMap {
		f.printFuncHeader(depth, currFunc.FuncName)
		f.processFuncSentences(depth+1, currFunc)
		f.PrintLabel(depth, fmt.Sprintf("} // func %s\n", currFunc.FuncName)) // func block end
	}

	f.mainFunc(depth)
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
