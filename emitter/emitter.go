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
	maxPatternNumber       int
	maxVarsNumber          int
	maxBracketsNumber      int
	nextSentenceEntryPoint int
	isThereFuncCall        bool
	sentenceInfo           sentenceInfo
	fixedVars              map[string]int
	patternCtx             patternContext
	isLeftMatching         bool
	funcInfo               *syntax.Function
	bracketsCurrentIndex   int
	bracketsNumerator      int
	entryPoints            []*entryPoint
	entryPointNumerator    int
}

type entryPoint struct {
	entryPoint  int
	actionIndex int
}

func (ctx *emitterContext) addPrevEntryPoint(newEntryPoint, newActionIndex int) {
	ctx.entryPoints = append(ctx.entryPoints, &entryPoint{entryPoint: newEntryPoint, actionIndex: newActionIndex})
}

func (ctx *emitterContext) clearEntryPoints() {
	ctx.entryPoints = make([]*entryPoint, 0)
}

func (ctx *emitterContext) getPrevEntryPoint() int {
	actionIndex := ctx.sentenceInfo.actionIndex

	if actionIndex == 0 {
		return -1
	}

	actionIndex--

	assemblyPresents := false
	s := ctx.sentenceInfo.sentence

	for i := len(ctx.entryPoints) - 1; i >= 0; i-- {
		if (s.Actions[ctx.entryPoints[i].actionIndex].ActionOp == syntax.COMMA) ||
			(s.Actions[ctx.entryPoints[i].actionIndex].ActionOp != syntax.REPLACE) {
			assemblyPresents = true
		}

		if (s.Actions[ctx.entryPoints[i].actionIndex].ActionOp == syntax.COLON) ||
			(s.Actions[ctx.entryPoints[i].actionIndex].ActionOp != syntax.DCOLON) &&
				assemblyPresents {
			return ctx.entryPoints[i].entryPoint
		}
	}

	return -1
}

func (ctx *emitterContext) needToAssembly() bool {
	actionIndex := ctx.sentenceInfo.actionIndex

	// Sentence Pattern
	if actionIndex == 0 {
		return true
	}

	// Sentence actions

	actionIndex--

	if actionIndex == 0 {
		return false
	}

	prevActionOp := ctx.sentenceInfo.sentence.Actions[actionIndex-1].ActionOp

	if prevActionOp == syntax.COLON || prevActionOp == syntax.DCOLON {
		return false
	}

	return true
}

func (f *Data) mainFunc(depth int, entryFuncName string) {

	f.PrintLabel(depth, "int main(int argc, char** argv)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "initAllocator(getHeapSizeFromCmdArgs(argc, argv));")
	f.PrintLabel(depth+1, "initLiteralData();")
	f.PrintLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", f.CurrTermNum))
	f.PrintLabel(depth+1, "initHeaps(vtermOffset);")
	f.PrintLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	f.PrintLabel(depth+1, "return 0;")
	f.PrintLabel(depth, "}")
}

func (f *Data) printInitLocals(depth int, ctx *emitterContext) {
	maxPatternNumber := ctx.maxPatternNumber
	maxVarsNumber := ctx.maxVarsNumber
	maxBracketsNumber := ctx.maxBracketsNumber

	f.PrintLabel(depth, "struct func_result_t funcRes;")
	f.PrintLabel(depth, "struct fragment_t* currFrag = 0;")
	f.PrintLabel(depth, "uint64_t fragmentOffset = 0;")
	f.PrintLabel(depth, "uint64_t rightBound = 0;")
	f.PrintLabel(depth, "int stretchingVarNumber = 0;")
	f.PrintLabel(depth, "int stretching = 0;")
	f.PrintLabel(depth, "int status = GC_OK;")
	f.PrintLabel(depth, "int prevStatus = GC_OK;")
	f.PrintLabel(depth, "int i = 0;")
	f.PrintLabel(depth, "int j = 0;")
	f.PrintLabel(depth, "if (entryStatus == FIRST_CALL)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("checkAndCleanHeaps(0, ENV_SIZE(%d, %d, %d));", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("initEnvData(CURR_FUNC_CALL->env, %d, %d, %d);", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else if (entryStatus == ROLL_BACK)")
	f.PrintLabel(depth+1, "stretching = 1;")
}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	sentencesCount := len(currFunc.Sentences)
	ctx := &emitterContext{}

	ctx.entryPointNumerator = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.maxBracketsNumber = getMaxBracketsCountInFunc(currFunc)
	ctx.funcInfo = currFunc

	f.printInitLocals(depth, ctx)

	f.PrintLabel(depth, "while(CURR_FUNC_CALL->entryPoint >= 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (CURR_FUNC_CALL->entryPoint)")
	f.PrintLabel(depth+1, "{")

	for sentenceIndex, sentence := range currFunc.Sentences {

		ctx.isLeftMatching = true
		ctx.fixedVars = make(map[string]int)
		ctx.sentenceInfo.init(sentencesCount, sentenceIndex, sentence)

		ctx.nextSentenceEntryPoint = ctx.entryPointNumerator +
			ctx.sentenceInfo.patternsCount + 2*ctx.sentenceInfo.callActionsCount
		ctx.bracketsNumerator = 0
		ctx.bracketsCurrentIndex = 0
		ctx.sentenceInfo.actionIndex = 0
		ctx.clearEntryPoints()

		f.matchingPattern(depth+1, ctx, sentence.Pattern.Terms)

		for index, a := range sentence.Actions {

			ctx.sentenceInfo.actionIndex = index + 1

			switch a.ActionOp {

			case syntax.COMMA: // ','
				f.ConstructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.REPLACE: // '='
				ctx.clearEntryPoints()
				f.ConstructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.COLON: // ':'
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.clearEntryPoints()
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.clearEntryPoints()
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		f.PrintLabel(depth+2, "CURR_FUNC_CALL->entryPoint = -1;")
		f.PrintLabel(depth+2, "break; //Successful end of sentence")
		f.PrintLabel(depth+1, "} // Pattern case end")
	}

	f.PrintLabel(depth+1, "} // Entry point switch end")
	f.PrintLabel(depth, "} // Main while end")

	f.PrintLabel(depth, "return funcRes;")
}

func (f *Data) predeclareFuncs(depth, funcsNumber int) {

	for i := 0; i < funcsNumber; i++ {
		f.PrintLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", i))
	}

	f.PrintLabel(depth, "")
}

func (f *Data) processFuncs(depth int, funcs map[string]*syntax.Function) {
	for _, currFunc := range funcs {
		f.printFuncHeader(depth, f.genFuncName(currFunc.Index))
		f.processFuncSentences(depth+1, currFunc)
		f.PrintLabel(depth, fmt.Sprintf("} // func %s:func_%d\n", currFunc.FuncName, currFunc.Index)) // func block end
	}

}

func processFile(f Data) {
	unit := f.Ast
	depth := 0

	f.PrintLabel(depth, fmt.Sprintf("// file:%s\n", f.Name))
	f.PrintHeaders()

	f.predeclareFuncs(depth, unit.FuncsTotalCount)

	f.printLiteralsAndHeapsInit(depth, unit)

	f.processFuncs(depth, unit.GlobMap)
	f.processFuncs(depth, unit.NestedFuncs)

	var goFunc *syntax.Function = nil
	var ok bool = false
	if goFunc, ok = unit.GlobMap["Go"]; !ok {
		if goFunc, ok = unit.GlobMap["GO"]; !ok {
			panic("Can't find entry point func! There is must be GO or Go func.")
		}
	}

	f.mainFunc(depth, f.genFuncName(goFunc.Index))
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}
