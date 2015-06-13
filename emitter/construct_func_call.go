package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
)

func (emt *EmitterData) constructFunctionalVTerm(depth int, term *syntax.Term, emittedName string, currFunc *syntax.Function) {

	env := make(map[string]syntax.ScopeVar, 0)
	rollback := 0

	if currFunc != nil {
		env = currFunc.Env
		rollback = boolToInt(currFunc.Rollback)
	}

	emt.printLabel(depth, "//Start construction func term.")
	emt.printLabel(depth, "ALLC_FRAG_LTERM(currTerm)")
	emt.printLabel(depth, "currTerm->fragment->offset = _memMngr.vtermsOffset;")
	emt.printLabel(depth, "ALLC_CLOSURE_VTERM")

	target := "_memMngr.vterms[currTerm->fragment->offset]"
	varStr := fmt.Sprintf("%s.closure", target)

	emt.printLabel(depth, fmt.Sprintf("ALLC_CLOSURE_STRUCT(%s, %s, %d, _memMngr.vterms[%d].str, %d)",
		varStr, emittedName, len(env), term.IndexInLiterals, rollback))

	emt.printLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range env {
		if parentLocalVarNumber, ok := emt.ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->locals + %d)->offset;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->locals + %d)->length;",
				target, needVarInfo.Number, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := emt.ctx.funcInfo.Env[needVarName]
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->offset = (CURR_FUNC_CALL->env->params + %d)->offset;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
			emt.printLabel(depth, fmt.Sprintf("(%s.closure->params + %d)->length = (CURR_FUNC_CALL->env->params + %d)->length;",
				target, needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	emt.printLabel(depth, "//Finish construction func term.")
}
