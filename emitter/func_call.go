package emitter

import (
	"fmt"
)

import (
	fk "bmstu-refal-compiler/emitter/funcs_keeper"
	"bmstu-refal-compiler/syntax"
)

func (f *Data) isFuncName(ident string, ctx *emitterContext) (*fk.FuncInfo, bool) {

	var funcInfo *fk.FuncInfo = nil
	var ok bool = false

	if funcInfo, ok = ctx.funcsKeeper.IsThereFunc(ident); ok {
		//Global func
		return funcInfo, ok
	}

	for _, scope := range ctx.scopeKeeper.GetAllScopes() {
		if funcInfo, ok := ctx.funcsKeeper.IsThereFunc(scope + ident); ok {
			//Nested func
			return funcInfo, ok
		}
	}

	return nil, false
}

func (f *Data) constructFunctionalVTerm(depth int, ctx *emitterContext, term *syntax.Term, funcInfo *fk.FuncInfo) {

	f.PrintLabel(depth, "//Start construction func term.")

	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = allocateClosure(%s, %d);", funcInfo.EmittedFuncName, len(funcInfo.EnvVarMap)))
	f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->ident = memMngr.vterms[%d].str;", term.IndexInLiterals))
	f.PrintLabel(depth, "currTerm->fragment->length = 1;")

	for needVarName, needVarInfo := range funcInfo.EnvVarMap {

		if parentLocalVarNumber, ok := ctx.sentenceInfo.scope.VarMap[needVarName]; ok {
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->env[%d] = env->locals[%d][%d];", needVarInfo.Number, ctx.entryPoint-1, parentLocalVarNumber.Number))
		} else {
			//Get from env of parent func
			parentEnvVarInfo, _ := ctx.currFuncInfo.EnvVarMap[needVarName]
			f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->env[%d] = env->params[%d];", needVarInfo.Number, parentEnvVarInfo.Number))
		}
	}

	f.PrintLabel(depth, "//Finish construction func term.")
}
