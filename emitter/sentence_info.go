package emitter

import (
	"bmstu-refal-compiler/syntax"
)

type sentenceInfo struct {
	index         int
	patternIndex  int
	actionIndex   int
	patternsCount int
	actionsCount  int
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
	sentenceInfo.actionsCount = len(s.Actions)
}

func (sentenceInfo *sentenceInfo) isLastAction() bool {

	return sentenceInfo.actionIndex+1 >= sentenceInfo.actionsCount
}

func (sentenceInfo *sentenceInfo) isNextMatchingAction() bool {

	if sentenceInfo.isLastAction() {
		return false
	}

	actions := sentenceInfo.sentence.Actions
	index := sentenceInfo.actionIndex

	if actions[index+1].ActionOp == syntax.COLON || actions[index+1].ActionOp == syntax.DCOLON {
		return true
	}

	return false
}

func getPatternsCount(s *syntax.Sentence) int {
	// +1 s.Pattern
	number := 1

	for _, a := range s.Actions {
		if a.ActionOp == syntax.COLON || a.ActionOp == syntax.DCOLON {
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
