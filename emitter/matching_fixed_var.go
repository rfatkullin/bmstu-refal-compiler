package emitter

import (
	"fmt"
)

func (f *Data) matchingFixedLocalSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d]", varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedEnvSymbolVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint
	lterm := fmt.Sprintf("CURR_FUNC_CALL->env->params[%d]", varNumber)
	f.matchingFixedSymbolVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedLocalExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("CURR_FUNC_CALL->env->locals[%d]", varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedEnvExprVar(depth int, ctx *emitterContext, varNumber int) {
	prevStretchVarNumber := ctx.patternCtx.prevEntryPoint

	lterm := fmt.Sprintf("CURR_FUNC_CALL->env->params[%d]", varNumber)
	f.matchingFixedExprVar(depth, prevStretchVarNumber, ctx, lterm)
}

func (f *Data) matchingFixedExprVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset + %s.fragment->length > rightBound)", lterm))
	f.printFailBlock(depth, prevStretchVarNumber, true)

	f.PrintLabel(depth, fmt.Sprintf("if (!eqFragment(fragmentOffset, %s.fragment->offset, %s.fragment->length))", lterm, lterm))
	f.printFailBlock(depth, prevStretchVarNumber, true)

	f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += %s.fragment->length;", lterm))
}

func (f *Data) matchingFixedSymbolVar(depth, prevStretchVarNumber int, ctx *emitterContext, lterm string) {

	f.PrintLabel(depth, "if (fragmentOffset >= rightBound ")
	f.PrintLabel(depth, fmt.Sprintf("|| (memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG) "+
		"|| (!eqSymbol(fragmentOffset, %s.fragment->offset)))", lterm))
	f.printFailBlock(depth, prevStretchVarNumber, true)

	f.PrintLabel(depth, "fragmentOffset++;")
}
