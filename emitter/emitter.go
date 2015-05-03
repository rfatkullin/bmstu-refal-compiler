package emitter

import (
	"fmt"
	"io"
)

import (
	"bmstu-refal-compiler/syntax"
)

type Data struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	CurrTermNum int
}

func Handle(done chan<- bool, fs <-chan Data) {
	for f := range fs {
		processFile(f)
		f.Close()
		done <- true
	}
}

func processFile(f Data) {
	unit := f.Ast
	depth := 0

	f.printHeadersAndDefs(depth, unit.FuncsTotalCount)
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

func (f *Data) printHeadersAndDefs(depth, funcsTotalCount int) {

	f.printLabel(depth, "#include <stdlib.h>")
	f.printLabel(depth, "#include <stdio.h>\n")

	f.printLabel(depth, "#include <vmachine.h>")
	f.printLabel(depth, "#include <memory_manager.h>")
	f.printLabel(depth, "#include <defines/gc_macros.h>")
	f.printLabel(depth, "#include <builtins/builtins.h>")
	f.printLabel(depth, "#include <allocators/data_alloc.h>")
	f.printLabel(depth, "#include <allocators/vterm_alloc.h>")
	f.printLabel(depth, "#include <defines/data_struct_sizes.h>")

	f.printLabel(depth, "")

	for i := 0; i < funcsTotalCount; i++ {
		f.printLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", i))
	}

	f.printLabel(depth, "")
}

func (f *Data) processFuncs(depth int, funcs map[string]*syntax.Function) {
	for _, currFunc := range funcs {
		f.printLabel(depth, fmt.Sprintf("// %s", currFunc.FuncName))
		f.printLabel(depth, fmt.Sprintf("struct func_result_t %s(int entryStatus) \n{", f.genFuncName(currFunc.Index)))
		f.processFuncSentences(depth+1, currFunc)
		f.printLabel(depth, fmt.Sprintf("} // %s\n", currFunc.FuncName))
	}

}

func (f *Data) processFuncSentences(depth int, currFunc *syntax.Function) {
	sentencesCount := len(currFunc.Sentences)
	ctx := &emitterContext{}

	ctx.entryPointNumerator = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.maxBracketsNumber = getMaxBracketsCountInFunc(currFunc)
	ctx.funcInfo = currFunc

	f.printInitLocals(depth, ctx)

	f.printLabel(depth, "while(CURR_FUNC_CALL->entryPoint >= 0)")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, "switch (CURR_FUNC_CALL->entryPoint)")
	f.printLabel(depth+1, "{")

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
				f.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.REPLACE: // '='
				ctx.clearEntryPoints()
				f.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.COLON: // ':'
				f.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.clearEntryPoints()
				f.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				f.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.clearEntryPoints()
				f.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				f.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		f.printLabel(depth+2, "CURR_FUNC_CALL->entryPoint = -1;")
		f.printLabel(depth+2, "break; //Successful end of sentence")
		f.printLabel(depth+1, "} // Pattern case end")
	}

	f.printLabel(depth+1, "} // Entry point switch end")
	f.printLabel(depth, "} // Main while end")

	f.printLabel(depth, "return funcRes;")
}

func (f *Data) printInitLocals(depth int, ctx *emitterContext) {
	maxPatternNumber := ctx.maxPatternNumber
	maxVarsNumber := ctx.maxVarsNumber
	maxBracketsNumber := ctx.maxBracketsNumber

	f.printLabel(depth, "struct func_result_t funcRes;")
	f.printLabel(depth, "struct fragment_t* currFrag = 0;")
	f.printLabel(depth, "uint64_t fragmentOffset = 0;")
	f.printLabel(depth, "uint64_t rightBound = 0;")
	f.printLabel(depth, "int stretchingVarNumber = 0;")
	f.printLabel(depth, "int stretching = 0;")
	f.printLabel(depth, "int status = GC_OK;")
	f.printLabel(depth, "int prevStatus = GC_OK;")
	f.printLabel(depth, "int i = 0;")
	f.printLabel(depth, "int j = 0;")
	f.printLabel(depth, "if (entryStatus == FIRST_CALL)")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, fmt.Sprintf("checkAndCleanHeaps(0, ENV_SIZE(%d, %d, %d));", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	f.printLabel(depth+1, fmt.Sprintf("initEnvData(CURR_FUNC_CALL->env, %d, %d, %d);", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	f.printLabel(depth, "}")
	f.printLabel(depth, "else if (entryStatus == ROLL_BACK)")
	f.printLabel(depth+1, "stretching = 1;")
}

func (f *Data) mainFunc(depth int, entryFuncName string) {

	f.printLabel(depth, "int main(int argc, char** argv)")
	f.printLabel(depth, "{")
	f.printLabel(depth+1, "initAllocator(getHeapSizeFromCmdArgs(argc, argv));")
	f.printLabel(depth+1, "initLiteralData();")
	f.printLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", f.CurrTermNum))
	f.printLabel(depth+1, "initHeaps(vtermOffset);")
	f.printLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	f.printLabel(depth+1, "return 0;")
	f.printLabel(depth, "}")
}
