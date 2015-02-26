package emitter

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
