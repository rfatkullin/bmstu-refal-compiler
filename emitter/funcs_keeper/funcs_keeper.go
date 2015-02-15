package funcs_keeper

import (
	"fmt"
)
import (
	sk "BMSTU-Refal-Compiler/emitter/scope_name_keeper"
	"BMSTU-Refal-Compiler/syntax"
)

type FuncInfo struct {
	ScopeKeeper      *sk.ScopeKeeper            //func full scope name where it was declaration
	EmittedFuncName  string                     //func name in C program pattern "func_#"
	*syntax.Function                            //func syntax tree struct
	EnvVarMap        map[string]syntax.ScopeVar //func all environment vars
}

type FuncsKeeper struct {
	funcs map[string]FuncInfo
}

func NewFuncsKeeper() *FuncsKeeper {
	return &FuncsKeeper{make(map[string]FuncInfo)}
}

func (funcsKeeper *FuncsKeeper) IsThereFunc(funcFullName string) (*FuncInfo, bool) {

	if funcInfo, ok := funcsKeeper.funcs[funcFullName]; ok {
		return &funcInfo, ok
	} else {
		return nil, ok
	}
}

func (funcsKeeper *FuncsKeeper) AddFunc(scopeKeeper *sk.ScopeKeeper, funcSyntTree *syntax.Function) *FuncInfo {

	funcNum := len(funcsKeeper.funcs)
	emittedFuncName := fmt.Sprintf("func_%d", funcNum)

	if funcSyntTree.FuncName == "" {
		funcSyntTree.FuncName = fmt.Sprintf("anonym_func_%d", funcNum)
	}

	funcInfo := FuncInfo{scopeKeeper.Copy(), emittedFuncName, funcSyntTree, make(map[string]syntax.ScopeVar, 8)}
	funcInfo.setEnv()
	funcsKeeper.funcs[scopeKeeper.GetFuncName(funcSyntTree.FuncName)] = funcInfo

	return &funcInfo
}

func (funcsKeeper *FuncsKeeper) AddBuiltinFunc(funcName string) {

	emittedFuncName := funcName

	funcInfo := FuncInfo{sk.NewScopeKeeper(), emittedFuncName, nil, make(map[string]syntax.ScopeVar, 0)}
	funcsKeeper.funcs[funcName] = funcInfo
}

func (funcsKeeper *FuncsKeeper) PrintAllFuncs() {
	for name, funcInfo := range funcsKeeper.funcs {
		fmt.Printf("%s\t%s\n", name, funcInfo.EmittedFuncName)
	}
}

func (funcInfo *FuncInfo) setEnv() {

	if funcInfo.EnvVarMap == nil {
		funcInfo.EnvVarMap = make(map[string]syntax.ScopeVar, 8)
	}

	env := funcInfo.EnvVarMap
	s := &funcInfo.Params

	for ; s != nil; s = s.Parent {
		if s.VarMap != nil {
			for varName, varInfo := range s.VarMap {
				env[varName] = syntax.ScopeVar{Number: len(env), VarType: varInfo.VarType}
			}
		}
	}
}
