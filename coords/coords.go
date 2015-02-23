package coords

import "fmt"

type (
	// Source file position.
	Pos struct {
		Offs      int // Offset from the beginning of file (in bytes).
		Line, Col int // Line and column (first line is 1, and also first column is 1).
	}

	// Coordinates of source file fragment.
	Fragment struct {
		Start  Pos // Position of fragment's first character.
		Follow Pos // Position of first character following the fragment.
	}
)

func (x Pos) String() string {
	if x.Line == 0 || x.Col == 0 {
		return fmt.Sprintf("[offs %d]", x.Offs)
	}

	return fmt.Sprintf("%d:%02d", x.Line, x.Col)
}

func (x Fragment) String() string {
	return fmt.Sprintf("(%v)-(%v)", x.Start, x.Follow)
}
