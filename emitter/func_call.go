package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (f *Data) constructFunctionalVTerm(depth int, ctx *emitterContext, term *syntax.Term, emittedName string, funcIndex int) {

	env := make(map[string]syntax.ScopeVar, 0)

	if funcIndex != -1 {
		env = f.Ast.FuncByNumber[funcIndex].Env
	}

	f.PrintLabel(depth, "//Start construction func term.")

	f.PrintLabel(depth, "currTerm = allocateFragmentLTerm();")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = allocateClosure(%s, %d, %d);", emittedName, len(env), term.IndexInLiterals))
	f.PrintLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {

		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->params[%d] = env->locals[%d][%d];", needVarInfo.Number, ctx.entryPoint-1, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := env[needVarName]
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->params[%d] = env->params[%d];", needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	f.PrintLabel(depth, "//Finish construction func term.")
}
