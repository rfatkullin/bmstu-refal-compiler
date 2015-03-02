package emitter

import (
	"fmt"
	"math/big"
	"strings"
)

import (
	"bmstu-refal-compiler/syntax"
)

func ReverseTerms(slice []*syntax.Term) (rSlice []*syntax.Term) {
	size := len(slice)
	rSlice = make([]*syntax.Term, 0)

	for index, _ := range slice {
		rSlice = append(rSlice, slice[size-index-1])
	}

	return
}

/// Unicode string --> char string"
func GetStrOfRunes(str string) string {

	runes := make([]string, 0)

	for _, rune := range str {
		runes = append(runes, fmt.Sprintf("%d", rune))
	}

	return strings.Join(runes, ",")
}

/// Bytes --> string
func GetStrOfBytes(num *big.Int) (str string, sign int, count int) {
	strs := make([]string, 0)
	bytes := num.Bytes()

	for _, byteVal := range num.Bytes() {
		strs = append(strs, fmt.Sprintf("%d", byteVal))
	}

	sign = 0
	if num.Sign() < 0 {
		sign = 1
	}

	str = strings.Join(strs, ",")
	count = len(bytes)

	return
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (ctx *emitterContext) setEnv(currFunc *syntax.Function) {
	ctx.env = make(map[string]syntax.ScopeVar, 0)
	s := &currFunc.Params

	for ; s != nil; s = s.Parent {
		if s.VarMap != nil {
			for varName, varInfo := range s.VarMap {
				ctx.env[varName] = syntax.ScopeVar{Number: len(ctx.env), VarType: varInfo.VarType}
			}
		}
	}
}

func (ctx *emitterContext) addNestedFunc(currFunc *syntax.Function) {
	ctx.nestedNamedFuncs = append(ctx.nestedNamedFuncs, currFunc)
	//	ctx.allFuncsMap[currFunc.Index] = currFunc
}
