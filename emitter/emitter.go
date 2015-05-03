package emitter

import (
	"fmt"
	"io"
)

import (
	"bmstu-refal-compiler/syntax"
)

type EmitterData struct {
	Name string
	Ast  *syntax.Unit
	io.WriteCloser
	currTermNum int
	context     *emitterContext
}

func ConstructEmitterData(name string, ast *syntax.Unit, io io.WriteCloser) EmitterData {
	return EmitterData{name, ast, io, 0, nil}
}

func Handle(done chan<- bool, fs <-chan EmitterData) {
	for emt := range fs {
		processFile(emt)
		emt.Close()
		done <- true
	}
}

func processFile(emt EmitterData) {
	unit := emt.Ast
	depth := 0

	emt.printHeadersAndDefs(depth, unit.FuncsTotalCount)
	emt.printLiteralsAndHeapsInit(depth, unit)

	emt.processFuncs(depth, unit.GlobMap)
	emt.processFuncs(depth, unit.NestedFuncs)

	var goFunc *syntax.Function = nil
	var ok bool = false
	if goFunc, ok = unit.GlobMap["Go"]; !ok {
		if goFunc, ok = unit.GlobMap["GO"]; !ok {
			panic("Can't find entry point func! There is must be GO or Go func.")
		}
	}

	emt.mainFunc(depth, emt.genFuncName(goFunc.Index))
}

func (emt *EmitterData) printHeadersAndDefs(depth, funcsTotalCount int) {

	emt.printLabel(depth, "#include <stdlib.h>")
	emt.printLabel(depth, "#include <stdio.h>\n")

	emt.printLabel(depth, "#include <vmachine.h>")
	emt.printLabel(depth, "#include <memory_manager.h>")
	emt.printLabel(depth, "#include <defines/gc_macros.h>")
	emt.printLabel(depth, "#include <builtins/builtins.h>")
	emt.printLabel(depth, "#include <allocators/data_alloc.h>")
	emt.printLabel(depth, "#include <allocators/vterm_alloc.h>")
	emt.printLabel(depth, "#include <defines/data_struct_sizes.h>")

	emt.printLabel(depth, "")

	for i := 0; i < funcsTotalCount; i++ {
		emt.printLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", i))
	}

	emt.printLabel(depth, "")
}

func (emt *EmitterData) processFuncs(depth int, funcs map[string]*syntax.Function) {
	for _, currFunc := range funcs {
		emt.printLabel(depth, fmt.Sprintf("// %s", currFunc.FuncName))
		emt.printLabel(depth, fmt.Sprintf("struct func_result_t %s(int entryStatus) \n{", emt.genFuncName(currFunc.Index)))
		emt.processFuncSentences(depth+1, currFunc)
		emt.printLabel(depth, fmt.Sprintf("} // %s\n", currFunc.FuncName))
	}

}

func (emt *EmitterData) processFuncSentences(depth int, currFunc *syntax.Function) {
	sentencesCount := len(currFunc.Sentences)
	ctx := &emitterContext{}

	ctx.entryPointNumerator = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.maxBracketsNumber = getMaxBracketsCountInFunc(currFunc)
	ctx.funcInfo = currFunc

	emt.printInitLocals(depth, ctx)

	emt.printLabel(depth, "while(CURR_FUNC_CALL->entryPoint >= 0)")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "switch (CURR_FUNC_CALL->entryPoint)")
	emt.printLabel(depth+1, "{")

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

		emt.matchingPattern(depth+1, ctx, sentence.Pattern.Terms)

		for index, a := range sentence.Actions {

			ctx.sentenceInfo.actionIndex = index + 1

			switch a.ActionOp {

			case syntax.COMMA: // ','
				emt.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.REPLACE: // '='
				ctx.clearEntryPoints()
				emt.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.COLON: // ':'
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.clearEntryPoints()
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.clearEntryPoints()
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		emt.printLabel(depth+2, "CURR_FUNC_CALL->entryPoint = -1;")
		emt.printLabel(depth+2, "break; //Successful end of sentence")
		emt.printLabel(depth+1, "} // Pattern case end")
	}

	emt.printLabel(depth+1, "} // Entry point switch end")
	emt.printLabel(depth, "} // Main while end")

	emt.printLabel(depth, "return funcRes;")
}

func (emt *EmitterData) printInitLocals(depth int, ctx *emitterContext) {
	maxPatternNumber := ctx.maxPatternNumber
	maxVarsNumber := ctx.maxVarsNumber
	maxBracketsNumber := ctx.maxBracketsNumber

	emt.printLabel(depth, "struct func_result_t funcRes;")
	emt.printLabel(depth, "struct fragment_t* currFrag = 0;")
	emt.printLabel(depth, "uint64_t fragmentOffset = 0;")
	emt.printLabel(depth, "uint64_t rightBound = 0;")
	emt.printLabel(depth, "int stretchingVarNumber = 0;")
	emt.printLabel(depth, "int stretching = 0;")
	emt.printLabel(depth, "int status = GC_OK;")
	emt.printLabel(depth, "int prevStatus = GC_OK;")
	emt.printLabel(depth, "int i = 0;")
	emt.printLabel(depth, "int j = 0;")
	emt.printLabel(depth, "if (entryStatus == FIRST_CALL)")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("checkAndCleanHeaps(0, ENV_SIZE(%d, %d, %d));", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	emt.printLabel(depth+1, fmt.Sprintf("initEnvData(CURR_FUNC_CALL->env, %d, %d, %d);", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	emt.printLabel(depth, "}")
	emt.printLabel(depth, "else if (entryStatus == ROLL_BACK)")
	emt.printLabel(depth+1, "stretching = 1;")
}

func (emt *EmitterData) mainFunc(depth int, entryFuncName string) {

	emt.printLabel(depth, "int main(int argc, char** argv)")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "initAllocator(getHeapSizeFromCmdArgs(argc, argv));")
	emt.printLabel(depth+1, "initLiteralData();")
	emt.printLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", emt.currTermNum))
	emt.printLabel(depth+1, "initHeaps(vtermOffset);")
	emt.printLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	emt.printLabel(depth+1, "return 0;")
	emt.printLabel(depth, "}")
}
