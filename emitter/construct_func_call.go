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
		rollback = boolToInt(currFunc.Rollback)
	}

	f.printLabel(depth, "//Start construction func term.")
	f.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")
	f.printCheckGCCondition(depth, "currTerm->fragment->offset", "chAllocateClosureVTerm(&status)")

	target := "_memMngr.vterms[currTerm->fragment->offset]"
	varStr := fmt.Sprintf("%s.closure", target)
	funcCallStr := fmt.Sprintf("chAllocateClosureStruct(%s, %d, _memMngr.vterms[%d].str, %d, &status)",
		emittedName, len(env), term.IndexInLiterals, rollback)

	f.printCheckGCCondition(depth, varStr, funcCallStr)

	f.printLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {
		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			f.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
			f.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->locals + %d)->length;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := ctx.funcInfo.Env[needVarName]
			f.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->params + %d)->offset;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
			f.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->params + %d)->length;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	f.printLabel(depth, "//Finish construction func term.")
}
