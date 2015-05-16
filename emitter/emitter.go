package emitter

import (
	"fmt"
	"io"
	"os"
)

import (
	"bmstu-refal-compiler/syntax"
)

type EmitterData struct {
	io.WriteCloser
	currTermNum  int
	ctx          *context
	dialect      int
	FuncByNumber map[int]*syntax.Function
	AllGlobals   map[string]*syntax.Function
}

func Handle(done chan<- bool, unitsChan <-chan *syntax.Unit, targetSourceName string, dialect int) {
	var units []*syntax.Unit = make([]*syntax.Unit, 0)

	for unit := range unitsChan {
		units = append(units, unit)
	}

	emt := constructEmitter(targetSourceName, dialect, units)
	emt.startEmit(units)

	emt.Close()

	for i := 0; i < len(units); i++ {
		done <- true
	}
}

func constructEmitter(targetSourceName string, dialect int, units []*syntax.Unit) EmitterData {

	var err error = nil

	emt := EmitterData{nil, 1, &context{}, dialect,
		make(map[int]*syntax.Function),
		make(map[string]*syntax.Function)}

	if emt.WriteCloser, err = os.Create(targetSourceName); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for _, unit := range units {

		for _, gFunc := range unit.GlobMap {
			emt.FuncByNumber[gFunc.Index] = gFunc
			emt.AllGlobals[gFunc.FuncName] = gFunc
		}

		for _, nFunc := range unit.NestedFuncs {
			emt.FuncByNumber[nFunc.Index] = nFunc
		}
	}

	return emt
}

func (emt *EmitterData) startEmit(units []*syntax.Unit) {
	var (
		depth     int              = 0
		entryFunc *syntax.Function = nil
		ok        bool             = false
	)

	if entryFunc, ok = emt.AllGlobals["GO"]; !ok {
		if entryFunc, ok = emt.AllGlobals["Go"]; !ok {
			panic("Can't find entry point func! There is must be GO or Go func.")
		}
	}

	emt.printHeadersAndDefs(depth, units)
	emt.printLiteralsAndHeapsInit(depth, units)

	for _, unit := range units {
		emt.ctx.ast = unit
		emt.processFile(depth)
	}

	emt.processEntryPoint(depth, entryFunc)
}

func (emt *EmitterData) processFile(depth int) {

	emt.processFuncs(depth, emt.ctx.ast.GlobMap)
	emt.processFuncs(depth, emt.ctx.ast.NestedFuncs)
}

func (emt *EmitterData) printHeadersAndDefs(depth int, units []*syntax.Unit) {

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

	for _, unit := range units {

		for _, gFunc := range unit.GlobMap {
			emt.printLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", gFunc.Index))
		}

		for _, nFunc := range unit.NestedFuncs {
			emt.printLabel(depth, fmt.Sprintf("struct func_result_t func_%d(int entryStatus);", nFunc.Index))
		}
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

	emt.ctx.initForNewFunc(currFunc)

	emt.printInitLocals(depth)

	emt.printLabel(depth, "while(CURR_FUNC_CALL->entryPoint >= 0)")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "switch (CURR_FUNC_CALL->entryPoint)")
	emt.printLabel(depth+1, "{")

	for sentenceIndex, sentence := range currFunc.Sentences {

		emt.ctx.initForNewSentence(sentencesCount, sentenceIndex, sentence)

		emt.matchingPattern(depth+1, sentence.Pattern.Terms)

		for index, a := range sentence.Actions {

			emt.ctx.sentenceInfo.actionIndex = index + 1

			switch a.ActionOp {

			case syntax.COMMA: // ','
				emt.constructAssembly(depth+2, a.Expr)
				break

			case syntax.REPLACE: // '='
				emt.ctx.clearEntryPoints()
				emt.constructAssembly(depth+2, a.Expr)
				break

			case syntax.COLON: // ':'
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.matchingPattern(depth+1, a.Expr.Terms)
				break

			case syntax.DCOLON: // '::'
				emt.ctx.clearEntryPoints()
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.matchingPattern(depth+1, a.Expr.Terms)
				break

			case syntax.TARROW: // '->'
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.constructFuncCallAction(depth+2, a.Expr.Terms)
				break

			case syntax.ARROW: // '=>'
				emt.ctx.clearEntryPoints()
				emt.printLabel(depth+1, "} // Pattern or Call Action case end\n")
				emt.constructFuncCallAction(depth+2, a.Expr.Terms)
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

func (emt *EmitterData) printInitLocals(depth int) {
	maxPatternNumber := emt.ctx.maxPatternNumber
	maxVarsNumber := emt.ctx.maxVarsNumber
	maxBracketsNumber := emt.ctx.maxBracketsNumber

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

func (emt *EmitterData) processEntryPoint(depth int, entryFunc *syntax.Function) {

	emt.printLabel(depth, "int main(int argc, char** argv)")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "initAllocator(getHeapSizeFromCmdArgs(argc, argv));")
	emt.printLabel(depth+1, "initLiteralData();")
	emt.printLabel(depth+1, fmt.Sprintf("uint64_t vtermOffset = initArgsData(UINT64_C(%d), argc, argv);", emt.currTermNum))
	emt.printLabel(depth+1, "initHeaps(vtermOffset);")
	emt.printLabel(depth+1, fmt.Sprintf("mainLoop(\"Go\", %s);", emt.genFuncName(entryFunc.Index)))
	emt.printLabel(depth+1, "return 0;")
	emt.printLabel(depth, "}")
}
