package funcs_keeper

import (
	"fmt"
)
import (
	"BMSTU-Refal-Compiler/syntax"
)

type FuncInfo struct {
	ScopeName        string                     //func full scope name where it was declaration
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

func (funcsKeeper *FuncsKeeper) AddFunc(scopeName string, funcSyntTree *syntax.Function) *FuncInfo {

	funcNum := len(funcsKeeper.funcs)
	emittedFuncName := fmt.Sprintf("func_%d", funcNum)

	if funcSyntTree.FuncName == "" {
		funcSyntTree.FuncName = fmt.Sprintf("anonym_func_%d", funcNum)
	}

	funcInfo := FuncInfo{scopeName, emittedFuncName, funcSyntTree, make(map[string]syntax.ScopeVar, 8)}
	funcsKeeper.funcs[scopeName+funcSyntTree.FuncName] = funcInfo

	return &funcInfo
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
