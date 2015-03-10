package emitter

import (
	"fmt"
)

func (f *Data) matchingIntLiteral(depth int, ctx *emitterContext, index int) {

	f.PrintLabel(depth-1, "//Matching int literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_INT_NUM_TAG || "+
		"!IntCmp(memMngr.vterms[fragmentOffset].intNum, memMngr.vterms[UINT64_C(%d)].intNum))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "fragmentOffset++;")
}

func (f *Data) mathcingDoubleLiteral(depth int, ctx *emitterContext, index int) {

	f.PrintLabel(depth-1, "//Matching double literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_DOUBLE_NUM_TAG || "+
		"doubleCmp(memMngr.vterms[fragmentOffset].doubleNum, memMngr.vterms[UINT64_C(%d)].doubleNum))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "fragmentOffset++;")
}

func (f *Data) matchingCompLiteral(depth int, ctx *emitterContext, index int) {

	f.PrintLabel(depth-1, "//Matching indetificator literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_IDENT_TAG || "+
		"!UStrCmp(memMngr.vterms[fragmentOffset].str, memMngr.vterms[UINT64_C(%d)].str))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "fragmentOffset++;")
}

func (f *Data) matchingStrLiteral(depth int, ctx *emitterContext, strLen, index int) {

	f.PrintLabel(depth-1, "//Matching string literal")

	f.PrintLabel(depth, fmt.Sprintf("if (fragmentOffset + UINT64_C(%d) - 1 >= rightCheckOffset)", strLen))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < UINT64_C(%d); i++)", strLen))
	f.PrintLabel(depth, "{")

	f.PrintLabel(depth+1, fmt.Sprintf("if (memMngr.vterms[fragmentOffset + i].tag != V_CHAR_TAG || "+
		"memMngr.vterms[fragmentOffset + i].ch != memMngr.vterms[UINT64_C(%d) + i].ch)", index))
	f.printFailBlock(depth+1, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "}")

	f.PrintLabel(depth, fmt.Sprintf("if (i < %d) // If check is failed", strLen))
	f.PrintLabel(depth+1, "break;")

	f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += UINT64_C(%d);", strLen))
}
