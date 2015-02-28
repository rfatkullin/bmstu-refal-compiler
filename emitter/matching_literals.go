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

	if ctx.isLeftMatching {
		f.PrintLabel(depth, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth, "fragmentOffset--;")
	}
}

func (f *Data) matchingCompLiteral(depth int, ctx *emitterContext, index int) {

	f.PrintLabel(depth-1, "//Matching indetificator literal")

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_IDENT_TAG || "+
		"!UStrCmp(memMngr.vterms[fragmentOffset].str, memMngr.vterms[UINT64_C(%d)].str))", index))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	if ctx.isLeftMatching {
		f.PrintLabel(depth, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth, "fragmentOffset--;")
	}
}

func (f *Data) matchingStrLiteral(depth int, ctx *emitterContext, str string) {

	f.PrintLabel(depth-1, fmt.Sprintf("//Matching %q literal", str))

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("for (i = 0; i < %d; i++)", len(str)))
	f.PrintLabel(depth, "{")

	if ctx.isLeftMatching {
		f.PrintLabel(depth+1, fmt.Sprintf("if (memMngr.vterms[fragmentOffset + i].tag != V_CHAR_TAG || "+
			"memMngr.vterms[fragmentOffset + i].ch != %q[i])", str))
	} else {
		f.PrintLabel(depth+1, fmt.Sprintf("if (memMngr.vterms[fragmentOffset - i].tag != V_CHAR_TAG || "+
			"memMngr.vterms[fragmentOffset - i].ch != %q[%d - i - 1])", str, len(str)))
	}

	f.printFailBlock(depth+1, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "}")

	f.PrintLabel(depth, fmt.Sprintf("if (i < %d) // If check is failed", len(str)))
	f.PrintLabel(depth+1, "break;")

	if ctx.isLeftMatching {
		f.PrintLabel(depth, fmt.Sprintf("fragmentOffset += %d;", len(str)))
	} else {
		f.PrintLabel(depth, fmt.Sprintf("fragmentOffset -= %d;", len(str)))
	}
}
