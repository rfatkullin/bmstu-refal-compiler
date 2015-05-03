package emitter

import (
	"fmt"
)

func (f *Data) matchingIntLiteral(depth int, ctx *emitterContext, index int) {

	f.printLabel(depth-1, "//Matching int literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_INT_NUM_TAG || "+
		"intCmp(_memMngr.vterms[fragmentOffset].intNum, _memMngr.vterms[UINT64_C(%d)].intNum))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.printLabel(depth, "fragmentOffset++;")
}

func (f *Data) mathcingDoubleLiteral(depth int, ctx *emitterContext, index int) {

	f.printLabel(depth-1, "//Matching double literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.printLabel(depth, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset].tag != V_DOUBLE_NUM_TAG || "+
		"doubleCmp(_memMngr.vterms[fragmentOffset].doubleNum, _memMngr.vterms[UINT64_C(%d)].doubleNum))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.printLabel(depth, "fragmentOffset++;")
}

func (f *Data) matchingCompLiteral(depth int, ctx *emitterContext, index int) {

	f.printLabel(depth-1, "//Matching indetificator literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.printLabel(depth, fmt.Sprintf("if (!((_memMngr.vterms[fragmentOffset].tag == V_IDENT_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].str, _memMngr.vterms[UINT64_C(%d)].str)) || "+
		"(_memMngr.vterms[fragmentOffset].tag == V_CLOSURE_TAG && "+
		"ustrEq(_memMngr.vterms[fragmentOffset].closure->ident, _memMngr.vterms[UINT64_C(%d)].str))))", index, index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.printLabel(depth, "fragmentOffset++;")
}

func (f *Data) matchingStrLiteral(depth int, ctx *emitterContext, strLen, index int) {

	f.printLabel(depth, "//Matching string literal")

	f.printLabel(depth, fmt.Sprintf("if (fragmentOffset + UINT64_C(%d) - 1 >= rightBound)", strLen))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.printLabel(depth, fmt.Sprintf("for (i = 0; i < UINT64_C(%d); i++)", strLen))
	f.printLabel(depth, "{")

	f.printLabel(depth+1, fmt.Sprintf("if (_memMngr.vterms[fragmentOffset + i].tag != V_CHAR_TAG || "+
		"_memMngr.vterms[fragmentOffset + i].ch != _memMngr.vterms[UINT64_C(%d) + i].ch)", index))
	f.printFailBlock(depth+1, ctx.patternCtx.prevEntryPoint, true)

	f.printLabel(depth, "}")

	f.printLabel(depth, fmt.Sprintf("if (i < %d) // If check is failed", strLen))
	f.printLabel(depth+1, "break;")

	f.printLabel(depth, fmt.Sprintf("fragmentOffset += UINT64_C(%d);", strLen))
}
