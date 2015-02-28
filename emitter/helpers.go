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
