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
	named := true

	// == -1 --> Builtins(no rollbacks, no env), != -1 --> Globs or Nested
	if funcIndex != -1 {
		currFunc := f.Ast.FuncByNumber[funcIndex]
		env = currFunc.Env
		rollback = BoolToInt(currFunc.Rollback)
		named = currFunc.HasName
		//fmt.Printf("Func name: %s\n", currFunc.FuncName)
	}

	f.PrintLabel(depth, "//Start construction func term.")

	f.PrintLabel(depth, "currTerm = chAllocateFragmentLTerm(1, &status);")
	f.printCheckGCCondition(depth)

	f.PrintLabel(depth, "currTerm->fragment->offset = chAllocateClosureVTerm(&status);")
	f.printCheckGCCondition(depth)

	target := "memMngr.vterms[currTerm->fragment->offset]"

	if named {
		f.PrintLabel(depth, fmt.Sprintf("%s.closure = chAllocateClosureStruct(%s, %d, memMngr.vterms[%d].str, %d, &status);",
			target, emittedName, len(env), term.IndexInLiterals, rollback))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("%s.closure = chAllocateClosureStruct(%s, %d, 0, %d, &status);",
			target, emittedName, len(env), rollback))
	}

	f.printCheckGCCondition(depth)

	f.PrintLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {
		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			f.PrintLabel(depth, fmt.Sprintf("%s.closure->params[%d].fragment->offset = _currFuncCall->env->locals[%d].fragment->offset;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
			f.PrintLabel(depth, fmt.Sprintf("%s.closure->params[%d].fragment->length = _currFuncCall->env->locals[%d].fragment->length;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := ctx.funcInfo.Env[needVarName]
			f.PrintLabel(depth, fmt.Sprintf("%s.closure->params[%d].fragment->offset = _currFuncCall->env->params[%d].fragment->offset;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
			f.PrintLabel(depth, fmt.Sprintf("%s.closure->params[%d].fragment->length = _currFuncCall->env->params[%d].fragment->length;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	f.PrintLabel(depth, "//Finish construction func term.")
}
