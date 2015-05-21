// Bauman Refal Compiler main package.
package main

// Standard packages
import (
	"fmt"
	"os"
	"path"
	"time"
)

// Bauman Refal Compiler packages
import (
	"bmstu-refal-compiler/chars"
	"bmstu-refal-compiler/cmdline"
	_ "bmstu-refal-compiler/coords"
	"bmstu-refal-compiler/emitter"
	"bmstu-refal-compiler/messages"
	"bmstu-refal-compiler/syntax"
	"bmstu-refal-compiler/tokens"
)

func changeExt(file, newExt string) string {
	ext := path.Ext(file)
	file = file[0 : len(file)-len(ext)]
	return fmt.Sprintf("%s.%s", file, newExt)
}

func main() {

	fmt.Println()
	fmt.Println("Bauman Refal Compiler")
	fmt.Println("(c) 2003-2010 Sergei Yu. Skorobogatov")
	fmt.Println()

	if cmdline.Minihelp {
		fmt.Println("Usage: refalc [ option... ] filename...")
		fmt.Println("(Use '--help' or '-h' to display list of options)")
		fmt.Println()
	}

	if cmdline.Ver {
		fmt.Println("Version 0.2.0")
		fmt.Println("Build information:")
		fmt.Println()
	}

	if cmdline.Help {
		fmt.Println("Usage: refalc [ option... ] filename...")
		fmt.Println()
		fmt.Println("Language options:")
		fmt.Println("  --refal7           or  -r7  Set Refal-7 mode")
		fmt.Println("  --refal5           or  -r5  Set Refal-5 mode")
		fmt.Println("  --relax-se         or  -se  Relax restrictions on side effects (Refal-7)")
		fmt.Println("  --nested-comments  or  -nc  Allow nested comments (Refal-7 and Refal-5)")
		fmt.Println("  --total-funcs <filename>")
		fmt.Println("    or    -tf <filename>      Set list of total functions")
		fmt.Println()
		fmt.Println("I18n options:")
		fmt.Println("  --list-cp                   Display list of supported code pages")
		fmt.Println("  --cp-source <value>")
		fmt.Println("    or    -cs <value>         Set code page of non-Unicode source files")
		fmt.Println()
		fmt.Println("Miscellaneous options:")
		fmt.Println("  --help             or  -h   Display this information")
		fmt.Println("  --ver              or  -v   Display version")
		fmt.Println("  --syntax           or  -s   Check for syntax errors, then stop")
		fmt.Println()
		fmt.Println("Internal structures output options (self-debugging):")
		fmt.Println("  --recovery-warnings         Report syntax error recovery")
		fmt.Println("  --lexem-list                Output list of lexems")
		fmt.Println("  --ptree                     Output parse tree")
		fmt.Println("  --initial-flast             Output FLAST built from parse tree")
		fmt.Println("  --fclass                    Output side effects mark-up")
		fmt.Println("  --refal-like                Choose Refal-like FLAST output")
		fmt.Println()
		fmt.Println("Time measurement options (self-profiling):")
		fmt.Println("  --parse-time                Measure parse time")
		fmt.Println("  --fb-time                   Measure FLAST building time")
		fmt.Println("  --se-time                   Measure SE-analysis time")
		fmt.Println("  --gen-time                  Measure code generation time")
		fmt.Println("  --repeat-cnt <value>        Set the repeat counter (default is 1)")
		fmt.Println()
	}

	if cmdline.ListCp {
		fmt.Println("List of supported code pages (to use with '--cp-source' option):")
		fmt.Println("  37    - IBM EBCDIC (US-Canada)")
		fmt.Println("  437   - OEM United States")
		fmt.Println("  500   - IBM EBCDIC (International)")
		fmt.Println("  708   - Arabic (ASMO 708)")
		fmt.Println("  720   - Arabic (DOS)")
		fmt.Println("  737   - Greek (DOS)")
		fmt.Println("  775   - Baltic (DOS)")
		fmt.Println("  850   - Western European (DOS)")
		fmt.Println("  852   - Central European (DOS)")
		fmt.Println("  855   - OEM Cyrillic")
		fmt.Println("  857   - Turkish (DOS)")
		fmt.Println("  858   - OEM Multilingual Latin I")
		fmt.Println("  860   - Portuguese (DOS)")
		fmt.Println("  861   - Icelandic (DOS)")
		fmt.Println("  862   - Hebrew (DOS)")
		fmt.Println("  863   - French Canadian (DOS)")
		fmt.Println("  864   - Arabic (864)")
		fmt.Println("  865   - Nordic (DOS)")
		fmt.Println("  866   - Cyrillic (DOS)")
		fmt.Println("  869   - Greek, Modern (DOS)")
		fmt.Println("  870   - IBM EBCDIC (Multilingual Latin-2)")
		fmt.Println("  874   - Thai (Windows)")
		fmt.Println("  875   - IBM EBCDIC (Greek Modern)")
		fmt.Println("  102   - IBM EBCDIC (Turkish Latin-5)")
		fmt.Println("  1047  - IBM Latin-1")
		fmt.Println("  1140  - IBM EBCDIC (US-Canada-Euro)")
		fmt.Println("  1141  - IBM EBCDIC (Germany-Euro)")
		fmt.Println("  1142  - IBM EBCDIC (Denmark-Norway-Euro)")
		fmt.Println("  1143  - IBM EBCDIC (Finland-Sweden-Euro)")
		fmt.Println("  1144  - IBM EBCDIC (Italy-Euro)")
		fmt.Println("  1145  - IBM EBCDIC (Spain-Euro)")
		fmt.Println("  1146  - IBM EBCDIC (UK-Euro)")
		fmt.Println("  1147  - IBM EBCDIC (France-Euro)")
		fmt.Println("  1148  - IBM EBCDIC (International-Euro)")
		fmt.Println("  1149  - IBM EBCDIC (Icelandic-Euro)")
		fmt.Println("  1250  - Central European (Windows)")
		fmt.Println("  1251  - Cyrillic (Windows)")
		fmt.Println("  1252  - Western European (Windows)")
		fmt.Println("  1253  - Greek (Windows)")
		fmt.Println("  1254  - Turkish (Windows)")
		fmt.Println("  1255  - Hebrew (Windows)")
		fmt.Println("  1256  - Arabic (Windows)")
		fmt.Println("  1257  - Baltic (Windows)")
		fmt.Println("  1258  - Vietnamese (Windows)")
		fmt.Println("  10000 - Western European (Mac)")
		fmt.Println("  10004 - Arabic (Mac)")
		fmt.Println("  10005 - Hebrew (Mac)")
		fmt.Println("  10006 - Greek (Mac)")
		fmt.Println("  10007 - Cyrillic (Mac)")
		fmt.Println("  10010 - Romanian (Mac)")
		fmt.Println("  10017 - Ukrainian (Mac)")
		fmt.Println("  10021 - Thai (Mac)")
		fmt.Println("  10029 - Central European (Mac)")
		fmt.Println("  10079 - Icelandic (Mac)")
		fmt.Println("  10081 - Turkish (Mac)")
		fmt.Println("  10082 - Croatian (Mac)")
		fmt.Println("  20105 - Western European (IA5)")
		fmt.Println("  20106 - German (IA5)")
		fmt.Println("  20107 - Swedish (IA5)")
		fmt.Println("  20108 - Norwegian (IA5)")
		fmt.Println("  20127 - US-ASCII")
		fmt.Println("  20269 - ISO-6937")
		fmt.Println("  20273 - IBM EBCDIC (Germany)")
		fmt.Println("  20277 - IBM EBCDIC (Denmark-Norway)")
		fmt.Println("  20278 - IBM EBCDIC (Finland-Sweden)")
		fmt.Println("  20280 - IBM EBCDIC (Italy)")
		fmt.Println("  20284 - IBM EBCDIC (Spain)")
		fmt.Println("  20285 - IBM EBCDIC (UK)")
		fmt.Println("  20290 - IBM EBCDIC (Japanese katakana)")
		fmt.Println("  20297 - IBM EBCDIC (France)")
		fmt.Println("  20420 - IBM EBCDIC (Arabic)")
		fmt.Println("  20423 - IBM EBCDIC (Greek)")
		fmt.Println("  20424 - IBM EBCDIC (Hebrew)")
		fmt.Println("  20833 - IBM EBCDIC (Korean Extended)")
		fmt.Println("  20838 - IBM EBCDIC (Thai)")
		fmt.Println("  20866 - Cyrillic (KOI8-R)")
		fmt.Println("  20871 - IBM EBCDIC (Icelandic)")
		fmt.Println("  20880 - IBM EBCDIC (Cyrillic Russian)")
		fmt.Println("  20905 - IBM EBCDIC (Turkish)")
		fmt.Println("  20924 - IBM Latin-1")
		fmt.Println("  21025 - IBM EBCDIC (Cyrillic Serbian-Bulgarian)")
		fmt.Println("  21027 - Ext Alpha Lowercase")
		fmt.Println("  21866 - Cyrillic (KOI8-U)")
		fmt.Println("  28591 - Western European (ISO)")
		fmt.Println("  28592 - Central European (ISO)")
		fmt.Println("  28593 - Latin 3 (ISO)")
		fmt.Println("  28594 - Baltic (ISO)")
		fmt.Println("  28595 - Cyrillic (ISO)")
		fmt.Println("  28596 - Arabic (ISO)")
		fmt.Println("  28597 - Greek (ISO)")
		fmt.Println("  28598 - Hebrew (ISO-Visual)")
		fmt.Println("  28599 - Turkish (ISO)")
		fmt.Println("  28603 - Estonian (ISO)")
		fmt.Println("  28605 - Latin 9 (ISO)")
		fmt.Println("  38598 - Hebrew (ISO-Logical)")
		fmt.Println("  1200  - UTF-16 (little-endian)")
		fmt.Println("  1201  - UTF-16 (big-endian)")
		fmt.Println("  1202  - UTF-32 (little-endian)")
		fmt.Println("  1203  - UTF-32 (big-endian)")
		fmt.Println("  65001 - UTF-8")
		fmt.Println()
		fmt.Println("NOTE 1: Default encoding is UTF-8.")
		fmt.Println()
		fmt.Println("NOTE 2: Unicode preamble (Byte-Order Mark) in the beginning")
		fmt.Println("  of the source file overrides '--cp-source' option.")
		fmt.Println()
	}

	done := make(chan bool, 16)
	fs := make(chan *syntax.Unit, 16)
	fileCount := 0
	targetSourceFileName := ""

	if len(cmdline.Sources) > 0 {
		targetSourceFileName = changeExt(cmdline.Sources[0], "c")
	}

	go emitter.Handle(done, fs, targetSourceFileName, cmdline.Dialect)

	t := time.Now()
	for _, x := range cmdline.Sources {
		fmt.Printf("%s:\n", x)

		f, err := os.Open(x)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		ms := make(chan messages.Data, 16)
		summ := make(chan messages.Summary)
		go messages.Handle(ms, summ)

		rs := make(chan chars.Rune, 1024)
		go chars.Handle(rs, ms, f, cmdline.CpSource)

		ts := make(chan tokens.Data, 128*1024)
		go tokens.Handle(ts, ms, rs, cmdline.Dialect, cmdline.NestedComments)

		if cmdline.LexemList {
			if d, err := os.Create(changeExt(x, "tokens")); err == nil {
				ts2 := make(chan tokens.Data, 128)
				go tokens.Intercept(ts2, d, ts)
				ts = ts2
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}

		ast := make(chan *syntax.Unit)
		go syntax.Handle(ast, ms, ts, cmdline.Dialect)

		if cmdline.Ptree {
			if d, err := os.Create(changeExt(x, "ptree")); err == nil {
				ast2 := make(chan *syntax.Unit, 1)
				go syntax.Intercept(ast2, d, ast)
				ast = ast2
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}

		fs <- <-ast
		fileCount++

		close(ms)
		summary := <-summ
		summary.Print(x)
	}

	close(fs)
	for i := 0; i < fileCount; i++ {
		<-done
	}

	fmt.Printf("Total time: %d\n", time.Since(t).Nanoseconds())
	fmt.Printf("Total files count: %d\n", fileCount)
}
