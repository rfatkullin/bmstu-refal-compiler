package emitter

import (
	"fmt"
	"unicode/utf8"
)

func (f *Data) matchingIntLiteral(depth int, ctx *emitterContext, intNumber int) {

	f.PrintLabel(depth-1, fmt.Sprintf("//Matching %d literal", intNumber))

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_INT_NUM_TAG || "+
		"memMngr.vterms[fragmentOffset].intNum != %d)", intNumber))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	if ctx.isLeftMatching {
		f.PrintLabel(depth, "fragmentOffset++;")
	} else {
		f.PrintLabel(depth, "fragmentOffset--;")
	}
}

func (f *Data) matchingCompLiteral(depth int, ctx *emitterContext, compSymbol string) {

	identLen := utf8.RuneCountInString(compSymbol)

	f.PrintLabel(depth-1, fmt.Sprintf("//Matching %q literal", compSymbol))

	f.printOffsetCheck(depth, ctx.patternCtx.prevEntryPoint, "")

	f.PrintLabel(depth, fmt.Sprintf("if (memMngr.vterms[fragmentOffset].tag != V_IDENT_TAG || "+
		"memMngr.vterms[fragmentOffset].str->length != UINT64_C(%d))", identLen))
	f.printFailBlock(depth, ctx.patternCtx.prevEntryPoint, true)

	f.PrintLabel(depth, "{")
	runesStr := GetStrOfRunes(compSymbol)
	f.PrintLabel(depth+1, fmt.Sprintf("struct v_string strTmp = (struct v_string){.head = (uint32_t[]){%s}, .length = UINT64_C(%d)};", runesStr, identLen))
	f.PrintLabel(depth+1, "if (!UStrCmp(memMngr.vterms[fragmentOffset].str, &strTmp))")
	f.printFailBlock(depth+1, ctx.patternCtx.prevEntryPoint, true)
	f.PrintLabel(depth, "}")

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
