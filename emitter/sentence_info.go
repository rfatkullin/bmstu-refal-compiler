package emitter

import (
	"BMSTU-Refal-Compiler/syntax"
)

type sentenceInfo struct {
	patternsCount int
	index         int
	patternIndex  int
	isLastPattern bool
	isLast        bool
	scope         *syntax.Scope
	sentence      *syntax.Sentence
}

func (sentenceInfo *sentenceInfo) init(sentencesCount, sentenceIndex int, s *syntax.Sentence) {
	sentenceInfo.index = sentenceIndex
	sentenceInfo.sentence = s
	sentenceInfo.scope = &s.Scope
	sentenceInfo.patternsCount = getPatternsCount(s)
	sentenceInfo.isLast = sentenceIndex == sentencesCount-1
	sentenceInfo.isLastPattern = sentenceInfo.patternsCount == 1
	sentenceInfo.patternIndex = 0
}

func getPatternsCount(s *syntax.Sentence) int {
	// +1 s.Pattern
	number := 1

	for _, a := range s.Actions {
		if a.ActionOp == syntax.COLON {
			number++
		}
	}

	return number
}

func getMaxPatternsAndVarsCount(currFunc *syntax.Function) (maxPatternsCount, maxVarCount int) {
	maxPatternsCount = 0
	maxVarCount = 0

	for _, s := range currFunc.Sentences {
		maxPatternsCount = max(maxPatternsCount, getPatternsCount(s))
		maxVarCount = max(maxVarCount, s.Scope.VarsNumber)
	}

	return maxPatternsCount, maxVarCount
}
