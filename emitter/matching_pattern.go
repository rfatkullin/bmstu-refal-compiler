package emitter

import (
	"fmt"
)

import (
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func (emt *EmitterData) matchingPattern(depth int, terms []*syntax.Term) {

	emt.checkAndAssemblyChain(depth + 1)

	emt.checkFragmentLength(depth+1, -1, false, terms)

	if len(terms) > 0 {
		emt.printLabel(depth+1, "else")
		emt.processPattern(depth+2, terms)
	}

	emt.processPatternFail(depth + 1)

	emt.ctx.addPrevEntryPoint(emt.ctx.entryPointNumerator, emt.ctx.sentenceInfo.actionIndex)
	emt.ctx.sentenceInfo.patternIndex++
}

func (emt *EmitterData) processEmptyPattern(depth int) {
	emt.printLabel(depth+1, "if (currFrag->length > 0)")
	emt.printRollBackBlock(depth+1, -1, false)
	emt.printLabel(depth+1, "break;")
}

func (emt *EmitterData) processPattern(depth int, terms []*syntax.Term) {

	emt.printLabel(depth, "while (stretchingVarNumber >= 0)")
	emt.printLabel(depth, "{")

	emt.printLabel(depth+1, "//From what stretchable variable start?")
	emt.printLabel(depth+1, "switch (stretchingVarNumber)")
	emt.printLabel(depth+1, "{")

	emt.ctx.patternCtx.entryPoint = 0
	emt.ctx.patternCtx.prevEntryPoint = -1

	emt.printFirstCase(depth, terms[0])

	emt.matchingTerms(depth+2, false, terms)

	emt.printLabel(depth+1, "} //pattern switch\n")

	emt.printLabel(depth+1, "if (!stretching)")
	emt.printLabel(depth+1, "{")
	emt.printLabel(depth+2, "if (fragmentOffset - currFrag->offset < currFrag->length)")
	emt.printRollBackBlock(depth+2, emt.ctx.patternCtx.prevEntryPoint, false)
	emt.printLabel(depth+2, "else")
	emt.printLabel(depth+3, "break; // Success!")
	emt.printLabel(depth+1, "}")

	emt.printLabel(depth, "} // Pattern while\n")
}

func (emt *EmitterData) printFirstCase(depth int, term *syntax.Term) {

	if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
		if _, ok := emt.ctx.fixedVars[term.Name]; !ok {
			return
		}
	}

	emt.ctx.patternCtx.entryPoint = 1
	emt.printLabel(depth+1, "case 0:")
}

func reverse(s []*syntax.Term) []*syntax.Term {
	rs := make([]*syntax.Term, 0)

	for i := len(s) - 1; i >= 0; i-- {
		rs = append(rs, s[i])
	}

	return rs
}

func (emt *EmitterData) checkRigidTerms(depth int, terms []*syntax.Term) []*syntax.Term {

	terms = emt.checkDirRigidTerms(depth, terms, LEFT_DIR)
	terms = reverse(terms)
	terms = emt.checkDirRigidTerms(depth, terms, RIGHT_DIR)
	terms = reverse(terms)

	return terms
}

func (emt *EmitterData) checkDirRigidTerms(depth int, terms []*syntax.Term, dir int) []*syntax.Term {
	i := 0

	for _, term := range terms {

		switch term.TermTag {
		case syntax.STR:
			emt.matchingRigidStrLiteral(depth, len(term.Value.Str), term.IndexInLiterals, dir)
			break
		case syntax.COMP:
			emt.matchingRigidCompLiteral(depth, term.IndexInLiterals, dir)
			break
		case syntax.INT:
			emt.matchingRigidIntLiteral(depth, term.IndexInLiterals, dir)
			break
		case syntax.FLOAT:
			emt.matchingRigidDoubleLiteral(depth, term.IndexInLiterals, dir)
			break
		case syntax.VAR:
			if term.VarType == tokens.VT_E || term.VarType == tokens.VT_V {
				return terms[i:]
			}
			emt.matchingRigidVars(depth, &term.Value, dir)
			break
		case syntax.EXPR:
			emt.matchingRigidBr(depth, dir)
			term.Exprs[0].Terms = emt.checkRigidTerms(depth, term.Exprs[0].Terms)
			break
		}

		i += 1
	}

	return terms[i:]
}

func (emt *EmitterData) isAllRigidTerms(terms []*syntax.Term) bool {

	for _, term := range terms {
		if term.TermTag == syntax.VAR && (term.VarType == tokens.VT_E || term.VarType == tokens.VT_V) {
			return false
		}
	}

	return true
}

func (emt *EmitterData) matchingTerms(depth int, inBrackets bool, terms []*syntax.Term) {
	parentMatchingOrder := emt.ctx.isLeftMatching

	terms = emt.checkRigidTerms(depth, terms)

	fmt.Printf("Length: %d\n", len(terms))

	emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))

	termsCount := len(terms)
	if termsCount == 0 {
		return
	}

	emt.ctx.isLeftMatching = !(terms[0].TermTag == syntax.R)

	for index, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			emt.matchingVariable(depth, &term.Value, emt.isAllRigidTerms(terms[index:]))
			break
		case syntax.STR:
			emt.matchingStrLiteral(depth, len(term.Value.Str), term.IndexInLiterals)
			break
		case syntax.COMP:
			emt.matchingCompLiteral(depth, term.IndexInLiterals)
			break
		case syntax.INT:
			emt.matchingIntLiteral(depth, term.IndexInLiterals)
			break
		case syntax.EXPR:
			emt.matchingExpr(depth, term.Exprs[0].Terms)
			break
		case syntax.FLOAT:
			emt.mathcingDoubleLiteral(depth, term.IndexInLiterals)
			break
		}
	}

	emt.ctx.isLeftMatching = parentMatchingOrder
}

func (emt *EmitterData) getMinLengthForTerms(terms []*syntax.Term) int {
	length := 0

	for _, term := range terms {

		switch term.TermTag {
		case syntax.VAR:
			varType := term.Value.VarType
			if varType == tokens.VT_T || varType == tokens.VT_S ||
				varType == tokens.VT_V {
				length += 1
			}
			break
		case syntax.STR:
			length += len(term.Value.Str)
			break
		case syntax.COMP, syntax.INT, syntax.EXPR, syntax.FLOAT:
			length += 1
			break
		}
	}

	return length
}

func (emt *EmitterData) matchingExpr(depth int, terms []*syntax.Term) {

	emt.ctx.bracketsNumerator++

	tmpBracketsCurrIndex := emt.ctx.brIndex
	bracketsIndex := emt.ctx.bracketsNumerator
	emt.ctx.brIndex = bracketsIndex

	emt.printLabel(depth, "//Check ().")
	emt.printOffsetCheck(depth, emt.ctx.patternCtx.prevEntryPoint, " || _memMngr.vterms[fragmentOffset].tag != V_BRACKETS_TAG")

	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[%d] = fragmentOffset;", bracketsIndex))
	emt.printLabel(depth, "rightBound = RIGHT_BOUND(fragmentOffset);")
	emt.printLabel(depth, "fragmentOffset = VTERM_BRACKETS(fragmentOffset)->offset;")

	emt.checkFragmentLength(depth, emt.ctx.patternCtx.prevEntryPoint, true, terms)

	emt.printLabel(depth, "//Start check in () terms.")
	emt.matchingTerms(depth, true, terms)

	emt.checkConsumeAllFragment(depth, emt.ctx.patternCtx.prevEntryPoint)

	emt.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->bracketsOffset[%d]);", tmpBracketsCurrIndex))
	emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FUNC_CALL->env->bracketsOffset[%d] + 1;", bracketsIndex))

	emt.printLabel(depth, "//End check in () terms.")

	emt.ctx.brIndex = tmpBracketsCurrIndex
}

func (emt *EmitterData) processPatternFail(depth int) {

	emt.printLabel(depth, "if (stretchingVarNumber < 0)")
	emt.printLabel(depth, "{")

	prevEntryPoint := emt.ctx.getPrevEntryPoint()
	//First pattern in current sentence
	if emt.ctx.sentenceInfo.patternIndex == 0 || prevEntryPoint == -1 {
		emt.processFailOfFirstPattern(depth + 1)
	} else {
		emt.processFailOfCommonPattern(depth+1, prevEntryPoint)
	}

	emt.printLabel(depth+1, "break;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) checkFragmentLength(depth, prevStertchingVarNumber int, withBreakStatement bool, terms []*syntax.Term) {

	if len(terms) == 0 {
		emt.printLabel(depth, "if (rightBound != fragmentOffset)")
	} else {
		emt.printLabel(depth, fmt.Sprintf("if (rightBound - fragmentOffset < %d)", emt.getMinLengthForTerms(terms)))
	}

	emt.printRollBackBlock(depth, prevStertchingVarNumber, withBreakStatement)
}

func (emt *EmitterData) processFailOfFirstPattern(depth int) {
	if emt.ctx.sentenceInfo.isLast {
		emt.printLabel(depth, "//First pattern of last sentence -> nothing to stretch -> fail!")
		emt.printLabel(depth, "funcRes = (struct func_result_t){.status = FAIL_RESULT, .fieldChain = 0, .callChain = 0};")
		emt.printLabel(depth, "CURR_FUNC_CALL->entryPoint = -1;")

	} else {
		emt.printLabel(depth, "//First pattern of current sentence -> jump to first pattern of next sentence!")
		emt.printLabel(depth, "stretching = 0;")
		emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", emt.ctx.nextSentenceEntryPoint))
		emt.printLabel(depth, "clearCurrFuncEnvData();")
	}
}

func (emt *EmitterData) processFailOfCommonPattern(depth, prevEntryPoint int) {
	emt.printLabel(depth, "//Jump to previouse pattern of same sentence!")
	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->entryPoint = %d;", prevEntryPoint))
}

func (emt *EmitterData) checkAndAssemblyChain(depth int) {
	patternIndex := emt.ctx.sentenceInfo.patternIndex

	emt.printLabel(depth, "if (!stretching)")
	emt.printLabel(depth, "{")

	if emt.ctx.sentenceInfo.actionIndex == 0 {
		if emt.ctx.sentenceInfo.index == 0 {
			emt.printLabel(depth+1, "ASSEMBLY_FIELD(0, CURR_FUNC_CALL->fieldOfView);")
		} else {
			emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[0];",
				patternIndex))
		}
	} else {
		if emt.ctx.needToAssembly() {
			emt.printLabel(depth+1, fmt.Sprintf("ASSEMBLY_FIELD(%d, CURR_FUNC_CALL->env->workFieldOfView);", patternIndex))
		} else {
			emt.printLabel(depth+1, fmt.Sprintf("CURR_FUNC_CALL->env->assembled[%d] = CURR_FUNC_CALL->env->assembled[%d];",
				patternIndex, patternIndex-1))
		}
	}

	emt.printLabel(depth, "} // !stretching")

	emt.printLabel(depth, fmt.Sprintf("currFrag = VTERM_BRACKETS(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emt.printLabel(depth, fmt.Sprintf("rightBound = RIGHT_BOUND(CURR_FUNC_CALL->env->assembled[%d]);", patternIndex))
	emt.printLabel(depth+1, "fragmentOffset = currFrag->offset;")
	emt.printLabel(depth, fmt.Sprintf("CURR_FUNC_CALL->env->bracketsOffset[0] = CURR_FUNC_CALL->env->assembled[%d];", patternIndex))
	emt.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = CURR_FUNC_CALL->env->stretchVarsNumber[%d];", emt.ctx.sentenceInfo.patternIndex))
}

func (emt *EmitterData) matchingVariable(depth int, value *tokens.Value, allRigid bool) {

	varInfo, isLocalVar := emt.ctx.sentenceInfo.scope.VarMap[value.Name]
	isFixedVar := true

	if !isLocalVar {
		varInfo = emt.ctx.funcInfo.Env[value.Name]
	} else {
		_, isFixedVar = emt.ctx.fixedVars[value.Name]
	}

	varNumber := varInfo.Number
	emt.printLabel(depth-1, fmt.Sprintf("//Matching %s variable", value.Name))

	switch value.VarType {
	case tokens.VT_T:
		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalExprVar(depth, varNumber)
			} else {
				emt.matchingFixedEnvExprVar(depth, varNumber)
			}
		} else {
			emt.matchingFreeTermVar(depth, varNumber)
			emt.ctx.fixedVars[value.Name] = emt.ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_S:
		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalSymbolVar(depth, varNumber)
			} else {
				emt.matchingFixedEnvSymbolVar(depth, varNumber)
			}

		} else {
			emt.matchingFreeSymbolVar(depth, varNumber)
			emt.ctx.fixedVars[value.Name] = emt.ctx.sentenceInfo.patternIndex
		}
		break

	case tokens.VT_E, tokens.VT_V:

		if isFixedVar {
			if isLocalVar {
				emt.matchingFixedLocalExprVar(depth, varNumber)
			} else {
				emt.matchingFixedEnvExprVar(depth, varNumber)
			}
		} else {
			emt.printLabel(depth-1, fmt.Sprintf("case %d:", emt.ctx.patternCtx.entryPoint))

			if !allRigid {
				if value.VarType == tokens.VT_E {
					emt.matchingFreeExprVar(depth, varNumber)
				} else {
					emt.matchingFreeVExprVar(depth, varNumber)
				}

				emt.ctx.patternCtx.prevEntryPoint = emt.ctx.patternCtx.entryPoint
				emt.ctx.patternCtx.entryPoint++
			} else {

				if value.VarType == tokens.VT_E {
					emt.freeExprVarGetRest(depth, varNumber)
				} else {
					emt.freeVExprVarGetRest(depth, varNumber)
				}
			}

			emt.ctx.fixedVars[value.Name] = emt.ctx.sentenceInfo.patternIndex
		}
		break
	}
}

func (emt *EmitterData) printRollBackBlock(depth, prevStretchVarNumber int, withBreakStatement bool) {

	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "stretching = 1;")
	emt.printLabel(depth+1, fmt.Sprintf("stretchingVarNumber = %d;", prevStretchVarNumber))
	if withBreakStatement {
		emt.printLabel(depth+1, "break;")
	}
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) printFailBlock(depth int) {

	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, "stretchingVarNumber = -1;")
	emt.printLabel(depth+1, "break;")
	emt.printLabel(depth, "}")
}

func (emt *EmitterData) printOffsetCheck(depth, prevStretchVarNumber int, optionalCond string) {

	emt.printLabel(depth, fmt.Sprintf("if (fragmentOffset >= rightBound%s)", optionalCond))
	emt.printRollBackBlock(depth, prevStretchVarNumber, true)
}

func (emt *EmitterData) checkConsumeAllFragment(depth, prevStretchVarNumber int) {
	emt.printLabel(depth, "if (fragmentOffset != rightBound)")
	emt.printRollBackBlock(depth, prevStretchVarNumber, true)
}
