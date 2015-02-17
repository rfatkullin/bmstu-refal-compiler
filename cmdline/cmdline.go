// Bauman Refal Compiler command line parser
package cmdline

import (
	"flag"
	"fmt"
	"os"
)

import "BMSTU-Refal-Compiler/chars"

// Variables bound to command line flags.
var (
	refal7           bool   // Set Refal-7 mode.
	refal5           bool   // Set Refal-5 mode.
	RelaxSe          bool   // Relax restrictions on side effects (Refal-7).
	NestedComments   bool   // Allow nested comments (Refal-7 and Refal-5).
	TotalFuncs       string // Set list of total functions
	ListCp           bool   // Display list of supported code pages
	CpSource         int    // Set code page of non-Unicode source files
	Help             bool   // Display this information.
	Ver              bool   // Display version.
	Syntax           bool   // Check for syntax errors, then stop.
	RecoveryWarnings bool   // Report syntax error recovery.
	LexemList        bool   // Output list of lexems.
	Ptree            bool   // Output parse tree.
	InitialFlast     bool   // Output FLAST built from parse tree.
	Fclass           bool   // Output side effects mark-up.
	RefalLike        bool   // Choose Refal-like FLAST output.
	ParseTime        bool   // Measure parse time.
	FbTime           bool   // Measure FLAST building time.
	SeTime           bool   // Measure SE-analysis time.
	GenTime          bool   // Measure code generation time.
	RepeatCnt        int    // Set the repeat counter (default is 1).
)

// Evaluated variables
var (
	Dialect  int      // Refal dialect
	Minihelp bool     // Print mini-help
	Sources  []string // Source file names
)

func init() {
	flag.BoolVar(&refal7, "refal7", false, "")
	flag.BoolVar(&refal7, "r7", false, "")

	flag.BoolVar(&refal5, "refal5", false, "")
	flag.BoolVar(&refal5, "r5", false, "")

	flag.BoolVar(&RelaxSe, "relax-se", false, "")
	flag.BoolVar(&RelaxSe, "se", false, "")

	flag.BoolVar(&NestedComments, "nested-comments", false, "")
	flag.BoolVar(&NestedComments, "nc", false, "")

	flag.StringVar(&TotalFuncs, "total-funcs", "", "")
	flag.StringVar(&TotalFuncs, "tf", "", "")

	flag.BoolVar(&ListCp, "list-cp", false, "")

	flag.IntVar(&CpSource, "cp-source", 1202, "")
	flag.IntVar(&CpSource, "cs", 1202, "")

	flag.BoolVar(&Help, "help", false, "")
	flag.BoolVar(&Help, "h", false, "")

	flag.BoolVar(&Ver, "ver", false, "")
	flag.BoolVar(&Ver, "v", false, "")

	flag.BoolVar(&Syntax, "syntax", false, "")
	flag.BoolVar(&Syntax, "s", false, "")

	flag.BoolVar(&RecoveryWarnings, "recovery-warnings", false, "")
	flag.BoolVar(&LexemList, "lexem-list", false, "")
	flag.BoolVar(&Ptree, "ptree", false, "")
	flag.BoolVar(&InitialFlast, "initial-flast", false, "")
	flag.BoolVar(&Fclass, "fclass", false, "")
	flag.BoolVar(&RefalLike, "refal-like", false, "")
	flag.BoolVar(&ParseTime, "parse-time", false, "")
	flag.BoolVar(&FbTime, "fb-time", false, "")
	flag.BoolVar(&SeTime, "se-time", false, "")
	flag.BoolVar(&GenTime, "gen-time", false, "")
	flag.IntVar(&RepeatCnt, "repeat-cnt", 1, "")

	flag.Parse()

	if !chars.CheckCp(CpSource) {
		fmt.Printf("illegal code page: %d\n", CpSource)
		os.Exit(1)
	}

	switch {
	case refal5:
		Dialect = 5
	case refal7:
		Dialect = 7
	default:
		Dialect = 7
	}

	Minihelp = flag.NFlag() == 0 && flag.NArg() == 0
	Sources = flag.Args()
}
