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
	for emitter := range fs {
		processFile(emitter)
		emitter.Close()
		done <- true
	}
}

func processFile(emitter EmitterData) {
	unit := emitter.Ast
	depth := 0

	emitter.printHeadersAndDefs(depth, unit.FuncsTotalCount)
	emitter.printLiteralsAndHeapsInit(depth, unit)

	emitter.processFuncs(depth, unit.GlobMap)
	emitter.processFuncs(depth, unit.NestedFuncs)

	var goFunc *syntax.Function = nil
	var ok bool = false
	if goFunc, ok = unit.GlobMap["Go"]; !ok {
		if goFunc, ok = unit.GlobMap["GO"]; !ok {
			panic("Can't find entry point func! There is must be GO or Go func.")
		}
	}

	emitter.mainFunc(depth, emitter.genFuncName(goFunc.Index))
}

func (emitter *EmitterData) printHeadersAndDefs(depth, funcsTotalCount int) {

	emitter.printLabel(depth, "#include <stdlib.h>")
	emitter.printLabel(depth, "#include <stdio.h>\n")

	emitter.printLabel(depth, "#include <vmachine.h>")
	emitter.printLabel(depth, "#include <memory_manager.h>")
	emitter.printLabel(depth, "#include <defines/gc_macros.h>")
	emitter.printLabel(depth, "#include <builtins/builtins.h>")
	emitter.printLabel(depth, "#include <allocators/data_alloc.h>")
	emitter.printLabel(depth, "#include <allocators/vterm_alloc.h>")
	emitter.printLabel(depth, "#include <defines/data_struct_sizes.h>")

	emitter.printLabel(depth, "")

	for i := 0; i < funcsTotalCount; i++ {
		emitter.printLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", i))
	}

	emitter.printLabel(depth, "")
}

func (emitter *EmitterData) processFuncs(depth int, funcs map[string]*syntax.Function) {
	for _, currFunc := range funcs {
		emitter.printLabel(depth, fmt.Sprintf("// %s", currFunc.FuncName))
		emitter.printLabel(depth, fmt.Sprintf("struct func_result_t %s(int entryStatus) \n{", emitter.genFuncName(currFunc.Index)))
		emitter.processFuncSentences(depth+1, currFunc)
		emitter.printLabel(depth, fmt.Sprintf("} // %s\n", currFunc.FuncName))
	}

}

func (emitter *EmitterData) processFuncSentences(depth int, currFunc *syntax.Function) {
	sentencesCount := len(currFunc.Sentences)
	ctx := &emitterContext{}

	ctx.entryPointNumerator = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.maxBracketsNumber = getMaxBracketsCountInFunc(currFunc)
	ctx.funcInfo = currFunc

	emitter.printInitLocals(depth, ctx)

	emitter.printLabel(depth, "while(CURR_FUNC_CALL->entryPoint >= 0)")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, "switch (CURR_FUNC_CALL->entryPoint)")
	emitter.printLabel(depth+1, "{")

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

		emitter.matchingPattern(depth+1, ctx, sentence.Pattern.Terms)

		for index, a := range sentence.Actions {

			ctx.sentenceInfo.actionIndex = index + 1

			switch a.ActionOp {

			case syntax.COMMA: // ','
				emitter.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.REPLACE: // '='
				ctx.clearEntryPoints()
				emitter.constructAssembly(depth+2, ctx, a.Expr)
				break

			case syntax.COLON: // ':'
				emitter.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emitter.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				ctx.clearEntryPoints()
				emitter.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emitter.matchingPattern(depth+1, ctx, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				emitter.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emitter.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				ctx.clearEntryPoints()
				emitter.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emitter.constructFuncCallAction(depth+2, ctx, a.Expr.Terms)
				break
			}
		}

		emitter.printLabel(depth+2, "CURR_FUNC_CALL->entryPoint = -1;")
		emitter.printLabel(depth+2, "break; //Successful end of sentence")
		emitter.printLabel(depth+1, "} // Pattern case end")
	}

	emitter.printLabel(depth+1, "} // Entry point switch end")
	emitter.printLabel(depth, "} // Main while end")

	emitter.printLabel(depth, "return funcRes;")
}

func (emitter *EmitterData) printInitLocals(depth int, ctx *emitterContext) {
	maxPatternNumber := ctx.maxPatternNumber
	maxVarsNumber := ctx.maxVarsNumber
	maxBracketsNumber := ctx.maxBracketsNumber

	emitter.printLabel(depth, "struct func_result_t funcRes;")
	emitter.printLabel(depth, "struct fragment_t* currFrag = 0;")
	emitter.printLabel(depth, "uint64_t fragmentOffset = 0;")
	emitter.printLabel(depth, "uint64_t rightBound = 0;")
	emitter.printLabel(depth, "int stretchingVarNumber = 0;")
	emitter.printLabel(depth, "int stretching = 0;")
	emitter.printLabel(depth, "int status = GC_OK;")
	emitter.printLabel(depth, "int prevStatus = GC_OK;")
	emitter.printLabel(depth, "int i = 0;")
	emitter.printLabel(depth, "int j = 0;")
	emitter.printLabel(depth, "if (entryStatus == FIRST_CALL)")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, fmt.Sprintf("checkAndCleanHeaps(0, ENV_SIZE(%d, %d, %d));", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	emitter.printLabel(depth+1, fmt.Sprintf("initEnvData(CURR_FUNC_CALL->env, %d, %d, %d);", maxVarsNumber, maxPatternNumber, maxBracketsNumber))
	emitter.printLabel(depth, "}")
	emitter.printLabel(depth, "else if (entryStatus == ROLL_BACK)")
	emitter.printLabel(depth+1, "stretching = 1;")
}

func (emitter *EmitterData) mainFunc(depth int, entryFuncName string) {

	emitter.printLabel(depth, "int main(int argc, char** argv)")
	emitter.printLabel(depth, "{")
	emitter.printLabel(depth+1, "initAllocator(getHeapSizeFromCmdArgs(argc, argv));")
	emitter.printLabel(depth+1, "initLiteralData();")
	emitter.printLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", emitter.currTermNum))
	emitter.printLabel(depth+1, "initHeaps(vtermOffset);")
	emitter.printLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", entryFuncName))
	emitter.printLabel(depth+1, "return 0;")
	emitter.printLabel(depth, "}")
}
