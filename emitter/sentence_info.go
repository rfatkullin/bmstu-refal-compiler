package emitter

import (
	"bmstu-refal-compiler/syntax"
)

type sentenceInfo struct {
	index            int
	patternIndex     int
	actionIndex      int
	patternsCount    int
	callActionsCount int
	actionsCount     int
	isLastPattern    bool
	isLast           bool
	scope            *syntax.Scope
	sentence         *syntax.Sentence
}

func (sentenceInfo *sentenceInfo) init(sentencesCount, sentenceIndex int, s *syntax.Sentence) {
	sentenceInfo.index = sentenceIndex
	sentenceInfo.sentence = s
	sentenceInfo.scope = &s.Scope
	sentenceInfo.patternsCount = getPatternsCount(s)
	sentenceInfo.callActionsCount = getCallActionsCount(s)
	sentenceInfo.isLast = sentenceIndex == sentencesCount-1
	sentenceInfo.isLastPattern = sentenceInfo.patternsCount == 1
	sentenceInfo.actionIndex = 0
	sentenceInfo.patternIndex = 0
	sentenceInfo.actionsCount = len(s.Actions)
}

func (sentenceInfo *sentenceInfo) isLastAction() bool {

	return sentenceInfo.actionIndex >= sentenceInfo.actionsCount
}

func (sentenceInfo *sentenceInfo) needToEval() bool {

	if sentenceInfo.isLastAction() {
		return false
	}

	actions := sentenceInfo.sentence.Actions
	index := sentenceInfo.actionIndex

	// index == actual index + 1
	switch actions[index].ActionOp {
	case syntax.COLON, // ':'
		syntax.DCOLON, // '::'
		syntax.TARROW, // '->'
		syntax.ARROW:  // '=>'
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

func getCallActionsCount(s *syntax.Sentence) int {
	number := 0

	for _, a := range s.Actions {
		if a.ActionOp == syntax.ARROW || a.ActionOp == syntax.TARROW {
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

func getBrackesCountInExpr(terms []*syntax.Term) int {
	count := 0

	for _, term := range terms {
		if term.TermTag == syntax.EXPR {
			count += getBrackesCountInExpr(term.Exprs[0].Terms) + 1
		}
	}

	return count
}

func getBracketsCountInSentence(s *syntax.Sentence) int {
	count := 0

	count += getBrackesCountInExpr(s.Pattern.Terms)

	for _, a := range s.Actions {
		if a.ActionOp == syntax.COLON || a.ActionOp == syntax.DCOLON {
			count += getBrackesCountInExpr(a.Expr.Terms)
		}
	}

	return count
}

func getMaxBracketsCountInFunc(currFunc *syntax.Function) int {
	maxBracketsCount := 0

	for _, s := range currFunc.Sentences {
		maxBracketsCount = max(maxBracketsCount, getBracketsCountInSentence(s))
	}

	return maxBracketsCount + 1 // +1 Assuming all expr in brackets
}
