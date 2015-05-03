package emitter

import (
	"fmt"
)

func (emitter *EmitterData) matchingIntLiteral(depth int, ctx *emitterContext, index int) {

	emitter.printLabel(depth-1, "//Matching int literal")

	emitter.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	emitter.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_INT_NUM_TAG || "+
		"intCmp(_memMngr.vterms[fragmentOffset].intNum, _memMngr.vterms[UINT64_C(%d)].intNum))", index))
	emitter.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	emitter.printLabel(depth, "fragmentOffset++;")
}

func (emitter *EmitterData) mathcingDoubleLiteral(depth int, ctx *emitterContext, index int) {

	emitter.printLabel(depth-1, "//Matching double literal")

	emitter.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	emitter.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_DOUBLE_NUM_TAG || "+
		"doubleCmp(_memMngr.vterms[fragmentOffset].doubleNum, _memMngr.vterms[UINT64_C(%d)].doubleNum))", index))
	emitter.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	emitter.printLabel(depth, "fragmentOffset++;")
}

func (emitter *EmitterData) matchingCompLiteral(depth int, ctx *emitterContext, index int) {

	emitter.printLabel(depth-1, "//Matching indetificator literal")

	emitter.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	emitter.printLabel(depth, fmt.Sprintf("if (!((_memMngr.vterms[fragmentOffset].tag == V_IDENT_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].str, _memMngr.vterms[UINT64_C(%d)].str)) || "+
		"(_memMngr.vterms[fragmentOffset].tag == V_CLOSURE_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].closure->ident, _memMngr.vterms[UINT64_C(%d)].str))))", index, index))
	emitter.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	emitter.printLabel(depth, "fragmentOffset++;")
}

func (emitter *EmitterData) matchingStrLiteral(depth int, ctx *emitterContext, strLen, index int) {

	emitter.printLabel(depth, "//Matching string literal")

	emitter.printLabel(depth, fmt.Sprintf("if (fragmentOffset + UINT64_C(%d) - 1 >= rightBound)", strLen))
	emitter.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	emitter.printLabel(depth, fmt.Sprintf("for (i = 0; i < UINT64_C(%d); i++)", strLen))
	emitter.printLabel(depth, "{")

	emitter.printLabel(depth+1, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset + i].tag != V_CHAR_TAG || "+
		"_memMngr.vterms[fragmentOffset + i].ch != _memMngr.vterms[UINT64_C(%d) + i].ch)", index))
	emitter.printFailBlock(depth+1, ctx.patternCtx.prevEntryPoint, true)

	emitter.printLabel(depth, "}")

	emitter.printLabel(depth, fmt.Sprintf("if (i < %d) // If check is failed", strLen))
	emitter.printLabel(depth+1, "break;")

	emitter.printLabel(depth, fmt.Sprintf("fragmentOffset += UINT64_C(%d);", strLen))
}
