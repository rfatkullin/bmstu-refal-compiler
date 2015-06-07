package emitter

import (
	"bmstu-refal-compiler/syntax"
)

type patternContext struct {
	entryPoint     int
	prevEntryPoint int
}

type context struct {
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
	brIndex                int
	bracketsNumerator      int
	entryPoints            []*entryPoint
	entryPointNumerator    int
	ast                    *syntax.Unit
}

type entryPoint struct {
	entryPoint  int
	actionIndex int
}

func (ctx *context) initForNewFunc(currFunc *syntax.Function) {
	ctx.entryPointNumerator = 0
	ctx.maxPatternNumber, ctx.maxVarsNumber = getMaxPatternsAndVarsCount(currFunc)
	ctx.maxBracketsNumber = getMaxBracketsCountInFunc(currFunc)
	ctx.funcInfo = currFunc
}

func (ctx *context) initForNewSentence(sentencesCount, sentenceIndex int, sentence *syntax.Sentence) {
	ctx.isLeftMatching = true
	ctx.fixedVars = make(map[string]int)
	ctx.sentenceInfo.init(sentencesCount, sentenceIndex, sentence)
	ctx.bracketsNumerator = 0
	ctx.brIndex = 0
	ctx.clearEntryPoints()

	acts := ctx.sentenceInfo.sentence.Actions

	endCallAct := 0
	if acts[len(acts)-1].ActionOp == syntax.ARROW {
		endCallAct = 1
	}

	ctx.nextSentenceEntryPoint = ctx.entryPointNumerator +
		ctx.sentenceInfo.patternsCount +
		ctx.sentenceInfo.assembliesCount +
		2*(ctx.sentenceInfo.callActionsCount-endCallAct) +
		endCallAct
}

func (ctx *context) addPrevEntryPoint(newEntryPoint, newActionIndex int) {
	ctx.entryPoints = append(ctx.entryPoints, &entryPoint{entryPoint: newEntryPoint, actionIndex: newActionIndex})
}

func (ctx *context) clearEntryPoints() {
	ctx.entryPoints = make([]*entryPoint, 0)
}

func (ctx *context) getPrevEntryPoint() int {
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

func (ctx *context) needToAssembly() bool {
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
