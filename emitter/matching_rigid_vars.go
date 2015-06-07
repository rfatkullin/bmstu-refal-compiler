package emitter

import (
	"fmt"
)

import (
	_ "bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

const (
	LEFT_DIR int = iota
	RIGHT_DIR
)

func (emt *EmitterData) matchingRigidVars(depth int, value *tokens.Value, dir int) {

	varInfo, isLocalVar := emt.ctx.sentenceInfo.scope.VarMap[value.Name]
	isFixedVar := true

	if !isLocalVar {
		varInfo = emt.ctx.funcInfo.Env[value.Name]
	} else {
		_, isFixedVar = emt.ctx.fixedVars[value.Name]
	}

	varNumber := varInfo.Number

	varName := emt.getVarName(varNumber, isLocalVar)

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			emt.matchingRigidFixedTermVar(depth, varName, dir)
		} else {
			emt.matchingRigidFreeTermVar(depth, varName, dir)
			emt.ctx.fixedVars[varName] = emt.ctx.sentenceInfo.patternIndex
		}

		break

	case tokens.VT_S:
		if isFixedVar {
			emt.matchingRigidFixedSymbolVar(depth, varName, dir)
		} else {
			emt.matchingRigidFreeSymbolVar(depth, varName, dir)
			emt.ctx.fixedVars[varName] = emt.ctx.sentenceInfo.patternIndex
		}
		break
	}
}

func (emt *EmitterData) matchingRigidFreeSymbolVar(depth int, varName string, dir int) {
	emt.printLabel(depth, "//Matching rigid free symbol var")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printPatternFailBlock(depth)

	emt.setFragmentOffset(depth, dir)

	emt.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag == V_BRACKETS_TAG)"))
	emt.printPatternFailBlock(depth)

	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("%s->offset = fragmentOffset;", varName))
	emt.printLabel(depth+1, fmt.Sprintf("%s->length = 1;", varName))
	emt.printLabel(depth, "}")

	emt.updateBracketOffset(depth, dir)
}

func (emt *EmitterData) matchingRigidFixedSymbolVar(depth int, varName string, dir int) {
	emt.printLabel(depth, "//Matching rigid fixed symbol var")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printPatternFailBlock(depth)

	emt.setFragmentOffset(depth, dir)

	emt.printLabel(depth, fmt.Sprintf("if (!CHECK_SYMB_VAR(fragmentOffset, %s->offset))", varName))
	emt.printPatternFailBlock(depth)

	emt.updateBracketOffset(depth, dir)
}

func (emt *EmitterData) matchingRigidFreeTermVar(depth int, varName string, dir int) {
	emt.printLabel(depth, "//Matching rigid free term var")

	emt.setFragmentOffset(depth, dir)

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printPatternFailBlock(depth)

	emt.printLabel(depth, "else")
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("%s->offset = fragmentOffset;", varName))
	emt.printLabel(depth+1, fmt.Sprintf("%s->length = 1;", varName))
	emt.printLabel(depth, "}")

	emt.updateBracketOffset(depth, dir)
}

func (emt *EmitterData) matchingRigidFixedTermVar(depth int, varName string, dir int) {
	emt.printLabel(depth, "//Matching rigid fixed term var")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printPatternFailBlock(depth)

	emt.setFragmentOffset(depth, dir)

	emt.printLabel(depth, fmt.Sprintf("if (!CHECK_TERM_VAR(fragmentOffset, %s->offset))", varName))
	emt.printPatternFailBlock(depth)

	emt.updateBracketOffset(depth, dir)
}

func (emt *EmitterData) updateBracketOffset(depth, dir int) {

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) setFragmentOffset(depth, dir int) {

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - 1;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) getVarName(varNumber int, local bool) string {

	if local {
		return fmt.Sprintf("(CURR_FUNC_CALL->env->locals + %d)", varNumber)
	}

	return fmt.Sprintf("(CURR_FUNC_CALL->env->params + %d)", varNumber)
}
