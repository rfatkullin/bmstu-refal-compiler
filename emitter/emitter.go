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
	maxVarsNumber          int
	nextSentenceEntryPoint int
	isThereFuncCall        bool
	sentenceInfo           sentenceInfo
	fixedVars              map[string]int
	patternCtx             patternContext
	isLeftMatching         bool
	funcInfo               *syntax.Function
}

func (f *Data) mainFunc(depth int, entryFuncName string) {

	f.PrintLabel(depth, "int main(int argc, char** argv)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "initLiteralData();")
	f.PrintLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", f.CurrTermNum))
	f.PrintLabel(depth+1, "initHeaps(vtermOffset);")
	f.PrintLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	f.PrintLabel(depth+1, "return 0;")
	f.PrintLabel(depth, "}")
}

func (f *Data) printInitLocals(depth, maxPatternNumber, varsNumber int) {

	f.PrintLabel(depth, "struct func_result_t funcRes;")
	f.PrintLabel(depth, "struct fragment_t* currFrag = 0;")
	f.PrintLabel(depth, "struct lterm_t* workFieldOfView = 0;")
	f.PrintLabel(depth, "uint64_t fragmentOffset = 0;")
	f.PrintLabel(depth, "uint64_t rightCheckOffset = 0;")
	f.PrintLabel(depth, "int stretchingVarNumber = 0;")
	f.PrintLabel(depth, "int stretching = 0;")
	f.PrintLabel(depth, "int success = 1;")
	f.PrintLabel(depth, "int i = 0;")
	f.PrintLabel(depth, "int j = 0;")
	f.PrintLabel(depth, "if (entryStatus == FIRST_CALL)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, fmt.Sprintf("checkAndCleanHeaps(0, ENV_SIZE(%d, %d));", varsNumber, maxPatternNumber))
	f.PrintLabel(depth+1, fmt.Sprintf("allocateEnvData(_currFuncCall->env, %d, %d);", varsNumber, maxPatternNumber))
	f.PrintLabel(depth, "}")
	f.PrintLabel(depth, "else if (entryStatus == ROLL_BACK)")
	f.PrintLabel(depth+1, "stretching = 1;")
}

func (f *Data) processFuncSentences(depth int, ctx *emitterContext, currFunc *syntax.Function) {
	sentencesCount := len(currFunc.Sentences)
	ctx.entryPoint = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.funcInfo = currFunc

	f.printInitLocals(depth, ctx.maxPatternNumber, ctx.maxVarsNumber)

	f.PrintLabel(depth, "while(_currFuncCall->entryPoint >= 0)")
	f.PrintLabel(depth, "{")
	f.PrintLabel(depth+1, "switch (_currFuncCall->entryPoint)")
	f.PrintLabel(depth+1, "{")

	for sentenceIndex, sentence := range currFunc.Sentences {

		ctx.isLeftMatching = true
		ctx.fixedVars = make(map[string]int)
		ctx.sentenceInfo.init(sentencesCount, sentenceIndex, sentence)

		ctx.nextSentenceEntryPoint = ctx.entryPoint +
			ctx.sentenceInfo.patternsCount + 2*ctx.sentenceInfo.callActionsCount
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
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.prevEntryPoint = -1
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.prevEntryPoint = -1
				f.PrintLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.ConstructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		f.PrintLabel(depth+2, "_currFuncCall->entryPoint = -1;")
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
