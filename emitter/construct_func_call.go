package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emt *EmitterData) constructFunctionalVTerm(depth int, ctx *emitterContext, term *syntax.Term, emittedName string, funcIndex int) {

	env := make(map[string]syntax.ScopeVar, 0)
	rollback := 0

	// == -1 --> Builtins(no rollbacks, no env), != -1 --> Globs or Nested
	if funcIndex != -1 {
		currFunc := emt.Ast.FuncByNumber[funcIndex]
		env = currFunc.Env
		rollback = boolToInt(currFunc.Rollback)
	}

	emt.printLabel(depth, "//Start construction func term.")
	emt.printCheckGCCondition(depth, "currTerm", "chAllocateFragmentLTerm(1, &status)")
	emt.printCheckGCCondition(depth, "currTerm->fragment->offset", "chAllocateClosureVTerm(&status)")

	target := "_memMngr.vterms[currTerm->fragment->offset]"
	varStr := fmt.Sprintf("%s.closure", target)
	funcCallStr := fmt.Sprintf("chAllocateClosureStruct(%s, %d, _memMngr.vterms[%d].str, %d, &status)",
		emittedName, len(env), term.IndexInLiterals, rollback)

	emt.printCheckGCCondition(depth, varStr, funcCallStr)

	emt.printLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {
		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->locals + %d)->length;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := ctx.funcInfo.Env[needVarName]
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->params + %d)->offset;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->params + %d)->length;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	emt.printLabel(depth, "//Finish construction func term.")
}
