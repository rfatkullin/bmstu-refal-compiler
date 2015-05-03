package emitter

import (
	"fmt"
)

func (emt *EmitterData) matchingFixedLocalSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)", varNumber)
	emt.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emt *EmitterData) matchingFixedEnvSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->params + %d)", varNumber)
	emt.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emt *EmitterData) matchingFixedLocalExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)", varNumber)
	emt.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emt *EmitterData) matchingFixedEnvExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("(CURR_FUNC_CALL->env->params + %d)", varNumber)
	emt.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (emt *EmitterData) matchingFixedExprVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	emt.printLabel(depth, fmt.Sprintf("if (fragmentOffset + %s->length > rightBound)", lterm))
	emt.printFailBlock(depth, prevStretchVarNumber, true)

	emt.printLabel(depth, fmt.Sprintf("if (!eqFragment(fragmentOffset, %s->offset, %s->length))", lterm, lterm))
	emt.printFailBlock(depth, prevStretchVarNumber, true)

	emt.printLabel(depth, fmt.Sprintf("fragmentOffset += %s->length;", lterm))
}

func (emt *EmitterData) matchingFixedSymbolVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	emt.printLabel(depth, "if (fragmentOffset >= rightBound ")
	emt.printLabel(depth, fmt.Sprintf("|| (_memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG) "+
		"|| (!eqSymbol(fragmentOffset, %s->offset)))", lterm))
	emt.printFailBlock(depth, prevStretchVarNumber, true)

	emt.printLabel(depth, "fragmentOffset++;")
}
