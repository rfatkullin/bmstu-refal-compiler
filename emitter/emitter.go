package emitter

import (
	"fmt"
	"io"
)

import (
	"bmstu-refal-compiler/syntax"
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

type emitterContext struct {
	entryPoint             int
	prevEntryPoint         int
	maxPatternNumber       int
	nextSentenceEntryPoint int
	isFuncCallInConstruct  bool
	sentenceInfo           sentenceInfo
	fixedVars              map[string]int
	patternCtx             patternContext
	isLeftMatching         bool
	funcInfo               *syntax.Function
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
	f.PrintLabel(depth, "struct lterm_t* helper = 0;")
	f.PrintLabel(depth, "struct lterm_t* currTerm = 0;")
	f.PrintLabel(depth, "struct lterm_t* funcTerm = 0;")
	f.PrintLabel(depth, "uint64_t fragmentOffset = 0;")
	f.PrintLabel(depth, "uint64_t rightCheckOffset = 0;")
	f.PrintLabel(depth, "int stretchingVarNumber = 0;")
	f.PrintLabel(depth, "int stretching = 0;")
	f.PrintLabel(depth, "int i = 0;")
	f.PrintLabel(depth, "int j = 0;")
	f.PrintLabel(depth, "if (entryStatus == FIRST_CALL)")
	f.PrintLabel(depth+1, fmt.Sprintf("allocateEnvData(env, %d, %d);", varsNumber, maxPatternNumber))
	f.PrintLabel(depth, "else if (entryStatus == ROLL_BACK)")
	f.PrintLabel(depth+1, "stretching = 1;")
}

func (f *Data) processFuncSentences(depth int, ctx *emitterContext, currFunc *syntax.Function) {
	maxVarsNumber := 0
	sentencesCount := len(currFunc.Sentences)
	ctx.entryPoint = 0
	ctx.maxPatternNumber, maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.funcInfo = currFunc

	f.printInitLocals(depth, ctx.maxPatternNumber, maxVarsNumber)

	f.PrintLabel(depth, "while(*entryPoint >= 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (*entryPoint)")
	f.PrintLabel(depth+1, "{")

	for sentenceIndex, sentence := range currFunc.Sentences {

		ctx.isLeftMatching = true
		ctx.fixedVars = make(map[string]int)
		ctx.sentenceInfo.init(sentencesCount, sentenceIndex, sentence)

		ctx.nextSentenceEntryPoint = ctx.entryPoint + ctx.sentenceInfo.patternsCount
		ctx.prevEntryPoint = -1

		f.matchingPattern(depth+1, ctx, sentence.Pattern.Terms)

		for index, a := range sentence.Actions {

			ctx.sentenceInfo.actionIndex = index

			switch a.ActionOp {

			case syntax.COMMA: // ','
				f.ConstructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.REPLACE: // '='
				ctx.prevEntryPoint = -1
				f.ConstructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.COLON: // ':'
				f.PrintLabel(depth+1, "} // Pattern case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.prevEntryPoint = -1
				f.PrintLabel(depth+1, "} // Pattern case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.prevEntryPoint = -1
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		f.PrintLabel(depth+2, "*entryPoint = -1;")
		f.PrintLabel(depth+2, "break; //Successful end of sentence")
		f.PrintLabel(depth+1, "} // Pattern case end")
	}

	f.PrintLabel(depth+1, "} // Entry point switch end")
	f.PrintLabel(depth, "} // Main while end")

	f.PrintLabel(depth, "return funcRes;")
}

func (f *Data) predeclareFuncs(depth, funcsNumber int) {

	for i := 0; i < funcsNumber; i++ {
		f.PrintLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int* entryPoint, struct env_t* env, struct lterm_t* fieldOfView, int entryStatus);", i))
	}

	f.PrintLabel(depth, "")
}

func (f *Data) processFuncs(depth int, ctx *emitterContext, funcs map[string]*syntax.Function) {
	for _, currFunc := range funcs {
		f.printFuncHeader(depth, f.genFuncName(currFunc.Index))
		f.processFuncSentences(depth+1, ctx, currFunc)
		f.PrintLabel(depth, fmt.Sprintf("} // func %s:func_%d\n", currFunc.FuncName, currFunc.Index)) // func block end
	}

}

func processFile(f Data) {
	unit := f.Ast
	depth := 0

	var ctx emitterContext

	f.PrintLabel(depth, fmt.Sprintf("// file:%s\n", f.Name))
	f.PrintHeaders()

	f.predeclareFuncs(depth, unit.FuncsTotalCount)

	f.printLiteralsAndHeapsInit(depth, unit)

	f.processFuncs(depth, &ctx, unit.GlobMap)
	f.processFuncs(depth, &ctx, unit.NestedFuncs)

	f.mainFunc(depth, f.genFuncName(unit.GlobMap["Go"].Index))
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
