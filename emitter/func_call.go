package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (f *Data) constructFunctionalVTerm(depth int, ctx *emitterContext, term *syntax.Term, emittedName string, funcIndex int) {

	env := make(map[string]syntax.ScopeVar, 0)
	rollback := 0

	// == -1 --> Builtins(no rollbacks, no env), != -1 --> Globs or Nested
	if funcIndex != -1 {
		currFunc := f.Ast.FuncByNumber[funcIndex]
		env = currFunc.Env
		rollback = BoolToInt(currFunc.Rollback)
	}

	f.PrintLabel(depth, "//Start construction func term.")

	f.PrintLabel(depth, "currTerm = allocateFragmentLTerm(1);")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = gcAllocateClosure(%s, %d, memMngr.vterms[%d].str, %d);",
		emittedName, len(env), term.IndexInLiterals, rollback))
	f.PrintLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {

		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->params[%d] = env->locals[%d];",
				needVarInfo.Number, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := env[needVarName]
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->params[%d] = env->params[%d];",
				needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	f.PrintLabel(depth, "//Finish construction func term.")
}
