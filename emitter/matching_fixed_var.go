package emitter

import (
	"fmt"
)

func (emitter *EmitterData) matchingFixedLocalSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)", varNumber)
	emitter.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emitter *EmitterData) matchingFixedEnvSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->params + %d)", varNumber)
	emitter.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emitter *EmitterData) matchingFixedLocalExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)", varNumber)
	emitter.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emitter *EmitterData) matchingFixedEnvExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->params + %d)", varNumber)
	emitter.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emitter *EmitterData) matchingFixedExprVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	emitter.printLabel(depth, fmt.Sprintf("if (fragmentOffset + %s->length > rightBound)", lterm))
	emitter.printFailBlock(depth, prevStretchVarNumber, true)

	emitter.printLabel(depth, fmt.Sprintf("if (!eqFragment(fragmentOffset, %s->offset, %s->length))", lterm, lterm))
	emitter.printFailBlock(depth, prevStretchVarNumber, true)

	emitter.printLabel(depth, fmt.Sprintf("fragmentOffset += %s->length;", lterm))
}

func (emitter *EmitterData) matchingFixedSymbolVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	emitter.printLabel(depth, "if (fragmentOffset >= rightBound ")
	emitter.printLabel(depth, fmt.Sprintf("|| (_memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG) "+
		"|| (!eqSymbol(fragmentOffset, %s->offset)))", lterm))
	emitter.printFailBlock(depth, prevStretchVarNumber, true)

	emitter.printLabel(depth, "fragmentOffset++;")
}
