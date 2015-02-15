package emitter

import (
	"fmt"
)

import (
	fk "BMSTU-Refal-Compiler/emitter/funcs_keeper"
	_ "BMSTU-Refal-Compiler/syntax"
)

func (f *Data) isFuncName(ident string, ctx *emitterContext) (*fk.FuncInfo, bool) {

	var funcInfo *fk.FuncInfo = nil
	var ok bool = false

	//fmt.Printf("Search global func: %s\n", ident)
	if funcInfo, ok = ctx.funcsKeeper.IsThereFunc(ident); ok {
		//Global func
		return funcInfo, ok
	}

	//fmt.Printf("Search nested func: %s\n", ident)

	for _, scope := range ctx.scopeKeeper.GetAllScopes() {
		//fmt.Printf("\tSearch in scope: %d %s\n", ind, scope)
		if funcInfo, ok := ctx.funcsKeeper.IsThereFunc(scope + ident); ok {
			return funcInfo, ok
		}
	}

	return nil, false
}

func (f *Data) constructFunctionalVTerm(depth int, ctx *emitterContext, ident string, funcInfo *fk.FuncInfo) {

	f.PrintLabel(depth, "//Start construction func term.")

	f.PrintLabel(depth, "currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));")
	f.PrintLabel(depth, "currTerm->tag = L_TERM_FRAGMENT_TAG;")
	f.PrintLabel(depth, "currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));")
	f.PrintLabel(depth, fmt.Sprintf("currTerm->fragment->offset = allocateClosure(%s, %d);", funcInfo.EmittedFuncName, len(funcInfo.EnvVarMap)))
	f.PrintLabel(depth, fmt.Sprintf("memMngr.vterms[currTerm->fragment->offset].closure->ident = %q;", ident))
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
