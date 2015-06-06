package emitter

import (
	"fmt"
)

func (emt *EmitterData) matchingIntLiteral(depth int, index int) {

	emt.printLabel(depth-1, "//Matching int literal")

	emt.printOffsetCheck(depth, emt.ctx.patternCtx.prevEntryPoint, "")

	emt.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_INT_NUM_TAG || "+
		"intCmp(_memMngr.vterms[fragmentOffset].intNum, _memMngr.vterms[UINT64_C(%d)].intNum))", index))
	emt.printRollBackBlock(depth, emt.ctx.patternCtx.prevEntryPoint, true)

	emt.printLabel(depth, "fragmentOffset++;")
}

func (emt *EmitterData) mathcingDoubleLiteral(depth int, index int) {

	emt.printLabel(depth-1, "//Matching double literal")

	emt.printOffsetCheck(depth, emt.ctx.patternCtx.prevEntryPoint, "")

	emt.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_DOUBLE_NUM_TAG || "+
		"doubleCmp(_memMngr.vterms[fragmentOffset].doubleNum, _memMngr.vterms[UINT64_C(%d)].doubleNum))", index))
	emt.printRollBackBlock(depth, emt.ctx.patternCtx.prevEntryPoint, true)

	emt.printLabel(depth, "fragmentOffset++;")
}

func (emt *EmitterData) matchingCompLiteral(depth int, index int) {

	emt.printLabel(depth-1, "//Matching indetificator literal")

	emt.printOffsetCheck(depth, emt.ctx.patternCtx.prevEntryPoint, "")

	emt.printLabel(depth, fmt.Sprintf("if (!((_memMngr.vterms[fragmentOffset].tag == V_IDENT_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].str, _memMngr.vterms[UINT64_C(%d)].str)) || "+
		"(_memMngr.vterms[fragmentOffset].tag == V_CLOSURE_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].closure->ident, _memMngr.vterms[UINT64_C(%d)].str))))", index, index))
	emt.printRollBackBlock(depth, emt.ctx.patternCtx.prevEntryPoint, true)

	emt.printLabel(depth, "fragmentOffset++;")
}

func (emt *EmitterData) matchingStrLiteral(depth int, strLen, index int) {

	emt.printLabel(depth, "//Matching string literal")

	emt.printLabel(depth, fmt.Sprintf("if (fragmentOffset + UINT64_C(%d) - 1 >= rightBound)", strLen))
	emt.printRollBackBlock(depth, emt.ctx.patternCtx.prevEntryPoint, true)

	emt.printLabel(depth, fmt.Sprintf("for (i = 0; i < UINT64_C(%d); i++)", strLen))
	emt.printLabel(depth, "{")

	emt.printLabel(depth+1, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset + i].tag != V_CHAR_TAG || "+
		"_memMngr.vterms[fragmentOffset + i].ch != _memMngr.vterms[UINT64_C(%d) + i].ch)", index))
	emt.printRollBackBlock(depth+1, emt.ctx.patternCtx.prevEntryPoint, true)

	emt.printLabel(depth, "}")

	emt.printLabel(depth, fmt.Sprintf("if (i < %d) // If check is failed", strLen))
	emt.printLabel(depth+1, "break;")

	emt.printLabel(depth, fmt.Sprintf("fragmentOffset += UINT64_C(%d);", strLen))
}
