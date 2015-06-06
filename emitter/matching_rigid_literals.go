package emitter

import (
	"fmt"
)

func (emt *EmitterData) matchingRigidBr(depth int, dir int) {
	emt.printLabel(depth-1, "//Matching ()")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - 1", emt.ctx.brIndex))
	}

	emt.printLabel(depth, "if (!CHECK_BR(fragmentOffset))")
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) matchingRigidIntLiteral(depth int, termInd, dir int) {

	emt.printLabel(depth-1, "//Matching int literal")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - 1", emt.ctx.brIndex))
	}

	emt.printLabel(depth, fmt.Sprintf("if (!CHECK_INT_LIT(fragmentOffset, UINT64_C(%d))", termInd))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) matchingRigidDoubleLiteral(depth int, termInd, dir int) {

	emt.printLabel(depth-1, "//Matching rigid double literal")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - 1", emt.ctx.brIndex))
	}

	emt.printLabel(depth, fmt.Sprintf("if (!CHECK_DOUBLE_LIT(fragmentOffset, UINT64_C(%d)))", termInd))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) matchingRigidCompLiteral(depth int, termInd, dir int) {

	emt.printLabel(depth-1, "//Matching rigid indetificator literal")

	emt.printLabel(depth, fmt.Sprintf("if (1 > CURR_FRAG_LENGTH(%d))", emt.ctx.brIndex))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - 1", emt.ctx.brIndex))
	}

	emt.printLabel(depth, fmt.Sprintf("if (!CHECK_COMP_LIT(fragmentOffset, UINT64_C(%d))", termInd))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)]++;", emt.ctx.brIndex))
	}
}

func (emt *EmitterData) matchingRigidStrLiteral(depth int, strLen, termInd, dir int) {

	emt.printLabel(depth, "//Matching rigid string literal")

	emt.printLabel(depth, fmt.Sprintf("if (UINT64_C(%d) > CURR_FRAG_LENGTH(%d))", strLen, emt.ctx.brIndex))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_LEFT(%d);", emt.ctx.brIndex))
	} else {
		emt.printLabel(depth, fmt.Sprintf("fragmentOffset = CURR_FRAG_RIGHT(%d) - UINT64_C(%d)", emt.ctx.brIndex, strLen))
	}

	emt.printLabel(depth, fmt.Sprintf("for (i = 0; i < UINT64_C(%d); i++)", strLen))
	emt.printLabel(depth, "{")
	emt.printLabel(depth+1, fmt.Sprintf("if (!CHECK_CHAR_LIT(fragmentOffset + i, UINT64_C(%d) + i))", termInd))
	emt.printLabel(depth+2, "break;")
	emt.printLabel(depth, "}")

	emt.printLabel(depth, fmt.Sprintf("if (i < UINT64_C(%d)) // If check is failed", strLen))
	emt.printFailBlock(depth)

	if dir == LEFT_DIR {
		emt.printLabel(depth, fmt.Sprintf("ENV->brLeftOffset[UINT64_C(%d)] += UINT64_C(%d);", emt.ctx.brIndex, strLen))
	} else {
		emt.printLabel(depth, fmt.Sprintf("ENV->brRightOffset[UINT64_C(%d)] += UINT64_C(%d);", emt.ctx.brIndex, strLen))
	}
}
