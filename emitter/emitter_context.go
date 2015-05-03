package emitter

import (
	"bmstu-refal-compiler/syntax"
)

type patternContext struct {
	entryPoint     int
	prevEntryPoint int
}

type emitterContext struct {
	maxPatternNumber       int
	maxVarsNumber          int
	maxBracketsNumber      int
	nextSentenceEntryPoint int
	isThereFuncCall        bool
	sentenceInfo           sentenceInfo
	fixedVars              map[string]int
	patternCtx             patternContext
	isLeftMatching         bool
	funcInfo               *syntax.Function
	bracketsCurrentIndex   int
	bracketsNumerator      int
	entryPoints            []*entryPoint
	entryPointNumerator    int
}

type entryPoint struct {
	entryPoint  int
	actionIndex int
}

func (ctx *emitterContext) addPrevEntryPoint(newEntryPoint, newActionIndex int) {
	ctx.entryPoints = append(ctx.entryPoints, &entryPoint{entryPoint: newEntryPoint, actionIndex: newActionIndex})
}

func (ctx *emitterContext) clearEntryPoints() {
	ctx.entryPoints = make([]*entryPoint, 0)
}

func (ctx *emitterContext) getPrevEntryPoint() int {
	actionIndex := ctx.sentenceInfo.actionIndex

	if actionIndex == 0 {
		return -1
	}

	actionIndex--

	assemblyPresents := false
	s := ctx.sentenceInfo.sentence

	for i := len(ctx.entryPoints) - 1; i >= 0; i-- {
		if (s.Actions[ctx.entryPoints[i].actionIndex].ActionOp == syntax.COMMA) ||
			(s.Actions[ctx.entryPoints[i].actionIndex].ActionOp != syntax.REPLACE) {
			assemblyPresents = true
		}

		if (s.Actions[ctx.entryPoints[i].actionIndex].ActionOp == syntax.COLON) ||
			(s.Actions[ctx.entryPoints[i].actionIndex].ActionOp != syntax.DCOLON) &&
				assemblyPresents {
			return ctx.entryPoints[i].entryPoint
		}
	}

	return -1
}

func (ctx *emitterContext) needToAssembly() bool {
	actionIndex := ctx.sentenceInfo.actionIndex

	// Sentence Pattern
	if actionIndex == 0 {
		return true
	}

	// Sentence actions

	actionIndex--

	if actionIndex == 0 {
		return false
	}

	prevActionOp := ctx.sentenceInfo.sentence.Actions[actionIndex-1].ActionOp

	if prevActionOp == syntax.COLON || prevActionOp == syntax.DCOLON {
		return false
	}

	return true
}
