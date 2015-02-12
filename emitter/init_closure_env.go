package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
)

func (f *Data) setEnvNestedFuncs(funcs map[string]*syntax.Function) {

	for _, currFunc := range funcs {

		currFunc.EnvVarMap = make(map[string]syntax.ScopeVar, 8)

		for s := &currFunc.Params; s != nil; s = s.Parent {
			if s.VarMap != nil {
				for varName, varInfo := range s.VarMap {
					currFunc.EnvVarMap[varName] = syntax.ScopeVar{Number: len(currFunc.EnvVarMap), VarType: varInfo.VarType}
				}
			}
		}

		fmt.Printf("Function %s\n", currFunc.FuncName)
		for varName, _ := range currFunc.EnvVarMap {
			fmt.Printf("\t%s\n", varName)
		}
	}
}
