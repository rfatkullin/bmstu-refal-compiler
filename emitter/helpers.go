package emitter

import (
	"fmt"
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
