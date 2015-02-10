package emitter

import (
	"fmt"
)

import (
	"BMSTU-Refal-Compiler/syntax"
)

func (f *Data) searchVarsInTerms(terms []*syntax.Term, envVarMap map[string]syntax.ScopeVar) {

	for _, term := range terms {

		if term.TermTag == syntax.VAR {
			if _, ok := envVarMap[term.Value.Name]; !ok {
				envVarMap[term.Value.Name] = syntax.ScopeVar{Number: len(envVarMap), VarType: term.VarType}
				fmt.Printf("Unknown var: %s\n", term.Value.Name)
			}
		}
	}
}

func (f *Data) setEnvNestedFuncs(funcs map[string]*syntax.Function) {

	for _, currFunc := range funcs {
		if currFunc.EnvVarMap == nil {
			currFunc.EnvVarMap = make(map[string]syntax.ScopeVar, 8)
		}

		for _, s := range currFunc.Sentences {
			f.searchVarsInTerms(s.Pattern.Terms, currFunc.EnvVarMap)

			for _, a := range s.Actions {
				f.searchVarsInTerms(a.Terms, currFunc.EnvVarMap)
			}
		}
	}
}
