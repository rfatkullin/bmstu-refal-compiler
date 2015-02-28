// Bauman Refal Compiler lexical scanner package
package tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

import (
	chars "bmstu-refal-compiler/chars"
	coords "bmstu-refal-compiler/coords"
	messages "bmstu-refal-compiler/messages"
)

type DomainTag int64

const (
	// End-of-input mark.
	END_OF_INPUT DomainTag = 1 << iota

	// Literals.
	STRING   // Character string.
	COMPOUND // Compound symbol (identifier).
	INTEGER  // Integer number.
	FLOAT    // Floating-point number.

	// Variable.
	VAR

	// Keywords.
	ENTRY  // $ENTRY
	SE     // $SE
	EXTERN // $EXTERN, and also $EXTRN and $EXTERNAL
	L      // $L
	R      // $R
	FUNC   // $FUNC

	// Special symbols.
	QMARK     // ?
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LEVAL     // <
	REVAL     // >
	COMMA     // ,
	REPLACE   // =
	COLON     // :
	DCOLON    // ::
	TARROW    // ->
	ARROW     // =>
	TILDE     // ~
	QUOTE     // `
	QLBRACE   // quoted {
	QRBRACE   // quoted }
	QLBRACKET // quoted [
	QRBRACKET // quoted ]
	QLEVAL    // quoted <
	QREVAL    // quoted >
)

var tagnames = map[DomainTag]string{
	END_OF_INPUT: "end of input",
	STRING:       "character string",
	COMPOUND:     "compound symbol",
	INTEGER:      "integer number",
	FLOAT:        "floating-point number",
	VAR:          "variable",
	ENTRY:        "$ENTRY",
	SE:           "$SE",
	EXTERN:       "$EXTERN",
	L:            "$L",
	R:            "$R",
	FUNC:         "$FUNC",
	QMARK:        "'?'",
	SEMICOLON:    "';'",
	LPAREN:       "'('",
	RPAREN:       "')'",
	LBRACE:       "'{'",
	RBRACE:       "'}'",
	LEVAL:        "'<'",
	REVAL:        "'>'",
	COMMA:        "','",
	REPLACE:      "'='",
	COLON:        "':'",
	DCOLON:       "'::'",
	TARROW:       "'->'",
	ARROW:        "'=>'",
	TILDE:        "'~'",
	QUOTE:        "'`'",
	QLBRACE:      "quoted '{'",
	QRBRACE:      "quoted '}'",
	QLBRACKET:    "quoted '['",
	QRBRACKET:    "quoted ']'",
	QLEVAL:       "quoted '<'",
	QREVAL:       "quoted '>'",
}

func (tag DomainTag) String() string {
	return tagnames[tag]
}

var tagpairs = map[DomainTag]DomainTag{
	LPAREN:    RPAREN,
	RPAREN:    LPAREN,
	LBRACE:    RBRACE,
	RBRACE:    LBRACE,
	LEVAL:     REVAL,
	REVAL:     LEVAL,
	QLBRACE:   QRBRACE,
	QRBRACE:   QLBRACE,
	QLBRACKET: QRBRACKET,
	QRBRACKET: QLBRACKET,
	QLEVAL:    QREVAL,
	QREVAL:    QLEVAL,
}

func (tag DomainTag) Pair() DomainTag {
	return tagpairs[tag]
}

func (mask DomainTag) ExpectedString() (res string) {
	tag := DomainTag(1)
	first := true
	for i := 0; i < len(tagnames); i++ {
		if mask&1 == 1 {
			if first {
				res += tagnames[tag]
				first = false
			} else {
				res += ", or " + tagnames[tag]
			}
		}

		tag <<= 1
		mask >>= 1
	}

	res += " expected"
	return
}

type VarType int

const (
	VT_E VarType = iota
	VT_V
	VT_T
	VT_S
	vt_count
)

const VAR_TYPES_NUM = int(vt_count)

var varTypes = []string{"e", "v", "t", "s"}

func (vt VarType) String() string {
	return varTypes[vt]
}

type Value struct {
	Str     []rune   // STRING
	Name    string   // COMPOUND | VAR
	IsIdent bool     // COMPOUND
	VarType          // VAR
	Int     *big.Int // INTEGER
	Float   float64  // FLOAT
}

type Data struct {
	Comment string
	DomainTag
	Value
	coords.Fragment
}

func (t *Data) Match(mask DomainTag) bool {
	return (t.DomainTag & mask) != 0
}

var hexDigits = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
	'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
}

var varTypes_7 = map[rune]VarType{
	'e': VT_E, 'v': VT_V, 't': VT_T, 's': VT_S,
}

var varTypes_5 = map[rune]VarType{
	'e': VT_E, 't': VT_T, 's': VT_S,
}

var keywords_7 = map[string]DomainTag{
	"$EXTERN": EXTERN,
	"$ENTRY":  ENTRY,
	"$SE":     SE,
	"$L":      L,
	"$R":      R,
	"$FUNC":   FUNC,
}

var keywords_5 = map[string]DomainTag{
	"$EXTERNAL": EXTERN,
	"$EXTERN":   EXTERN,
	"$EXTRN":    EXTERN,
	"$ENTRY":    ENTRY,
}

var maxMacrodigit = big.NewInt(0xFFFFFFFF)

type charbuf []rune

func newCharbuf() charbuf { return make([]rune, 0, 128) }

func (cb *charbuf) reset() { *cb = (*cb)[:0] }

func (cb *charbuf) save(c rune) { *cb = append(*cb, c) }

func (cb *charbuf) String() string { return string(([]rune)(*cb)) }

var (
	whitespaceRunes = newBuilder([]string{"R"})
	decDigitRunes   = newBuilder([]string{"R"})
	hexDigitRunes   = newBuilder([]string{"R"})
	identTailRunes  = newBuilder([]string{"R5", "R7"})
	opTailRunes     = newBuilder([]string{"QR7", "QQR7"})
)

func init() {
	dummy := func() {}

	whitespaceRunes.bind(map[string]func(){
		`(R) " \t\n\v\f\u001A"`: dummy,
	},
	)

	decDigitRunes.bind(map[string]func(){
		`(R) '0'-'9'`: dummy,
	},
	)

	hexDigitRunes.bind(map[string]func(){
		`(R) '0'-'9', 'A'-'F', 'a'-'f'`: dummy,
	},
	)

	opTailRunes.bind(map[string]func(){
		`(QR7, QQR7) "!#$%&*+,:;=?\\^|~", '/', '-', '.';` +
			`(QR7) '<', '>'`: dummy,
	},
	)

	identTailRunes.bind(map[string]func(){
		`(R5, R7) 'A'-'Z', 'a'-'z', '0'-'9', '_';` +
			`(R5) '-';` +
			`(R7) letter`: dummy,
	},
	)
}

func Handle(ts chan<- Data, ms chan<- messages.Data, runes <-chan chars.Rune,
	dialect int, nc bool) {

	var r_prev, r, r_next chars.Rune
	r_next = <-runes

	next := func() {
		if r.Code != chars.EOF {
			r_prev, r, r_next = r, r_next, <-runes
		}
	}

	var cb charbuf
	var start coords.Pos

	newData := func(tag DomainTag, comment string) Data {
		return Data{
			Comment:   comment,
			DomainTag: tag,
			Fragment:  coords.Fragment{start, r.Pos},
		}
	}

	scanTail := func(tailRunes *builder) {
		cb.reset()
		for {
			cb.save(r.Code)
			next()

			if tailRunes.dispatch(r.Code) == nil {
				break
			}
		}
	}

	write := func(tag DomainTag) {
		ts <- newData(tag, tag.String())
	}

	writeString := func() {
		data := newData(STRING, fmt.Sprintf("%v '%s'", STRING, cb.String()))
		data.Str = make([]rune, len(cb))
		copy(data.Str, cb)
		ts <- data
	}

	nameHash := make(map[string]string, 1024)

	writeCompound := func(s string, isIdent bool) {
		if name, ok := nameHash[s]; ok {
			s = name
		} else {
			nameHash[s] = s
		}

		data := newData(COMPOUND, fmt.Sprintf("%v '%s'", COMPOUND, s))
		data.Name = s
		data.IsIdent = isIdent
		ts <- data
	}

	writeCbCompound := func(isIdent bool) {
		writeCompound(cb.String(), isIdent)
	}

	writeVar := func(vt VarType) {
		s := fmt.Sprintf("%s.%s", vt.String(), cb.String())
		data := newData(VAR, fmt.Sprintf("%v %s", VAR, s))
		data.VarType, data.Name = vt, s
		ts <- data
	}

	writeInt := func(x *big.Int) {
		data := newData(INTEGER, fmt.Sprintf("%v %v", INTEGER, x))
		data.Int = x
		ts <- data
	}

	writeFloat := func(x float64) {
		data := newData(FLOAT, fmt.Sprintf("%v %v", FLOAT, x))
		data.Float = x
		ts <- data
	}

	err := func(s string) {
		ms <- messages.Data{r.Pos, messages.ERROR, s}
	}

	errStart := func(s string) {
		ms <- messages.Data{start, messages.ERROR, s}
	}

	errUnexpectedChar := func() {
		err("unexpected character")
	}

	warn := func(s string) {
		ms <- messages.Data{r.Pos, messages.WARNING, s}
	}

	bg := newBuilder([]string{"R5", "R7", "QR7", "QQR7"})

	bg.bind(map[string]func(){
		`(R5, R7, QR7, QQR7) " \t\n\v\f\u001A"`: func() {
			for {
				next()
				if whitespaceRunes.dispatch(r.Code) == nil {
					break
				}
			}
		},
	},
	)

	skipMComment := func() {
		for depth := 0; ; next() {
			switch {
			case r.Code == '/' && r_next.Code == '*':
				next()
				depth++
			case r.Code == '*' && r_next.Code == '/':
				next()
				next()
				if depth--; depth == 0 || !nc {
					return
				}
			case r.Code == chars.EOF:
				err("'*/' expected")
				return
			}
		}
	}

	skipComment := func() {
		for r.Code != '\n' && r.Code != chars.EOF {
			next()
		}
		next()
	}

	shortArith := func(s string) {
		if r_prev.Code == '<' {
			next()
			writeCompound(s, true)
		} else {
			errUnexpectedChar()
			next()
		}
	}

	scanOp := func() {
		scanTail(opTailRunes)
		writeCbCompound(true)
	}

	bg.bind(map[string]func(){
		`(R7) '/'`: func() {
			if r_next.Code == '*' {
				skipMComment()
			} else if r_next.Code == '/' {
				skipComment()
			} else {
				errUnexpectedChar()
				next()
			}
		},

		`(QR7, QQR7) '/'`: func() {
			if r_next.Code == '*' {
				skipMComment()
			} else if r_next.Code == '/' {
				skipComment()
			} else {
				scanOp()
			}
		},

		`(R5) '*'`: func() {
			if r.Col == 1 {
				skipComment()
			} else {
				shortArith("Mul")
			}
		},

		`(R5) '%'`: func() { shortArith("Mod") },

		`(R5) '+'`: func() { shortArith("Add") },

		`(R5) '-'`: func() { shortArith("Sub") },

		`(R5) '/'`: func() {
			if r_next.Code == '*' {
				skipMComment()
			} else {
				shortArith("Div")
			}
		},
	},
	)

	bg.bind(map[string]func(){
		`(R7) '?'`: func() { next(); write(QMARK) },

		`(R5, R7) ';'`: func() { next(); write(SEMICOLON) },

		`(R5, R7, QR7, QQR7) '('`: func() { next(); write(LPAREN) },

		`(R5, R7, QR7, QQR7) ')'`: func() { next(); write(RPAREN) },

		`(R5, R7) '{'`: func() { next(); write(LBRACE) },

		`(QR7, QQR7) '{'`: func() { next(); write(QLBRACE) },

		`(R5, R7) '}'`: func() { next(); write(RBRACE) },

		`(QR7, QQR7) '}'`: func() { next(); write(QRBRACE) },

		`(QR7, QQR7) '['`: func() { next(); write(QLBRACKET) },

		`(QR7, QQR7) ']'`: func() { next(); write(QRBRACKET) },

		`(R5) '<'`: func() {
			next()
			write(LEVAL)
			if whitespaceRunes.dispatch(r.Code) != nil {
				err("white space not allowed after '<'")
			}
		},

		`(R7) '<'`: func() { next(); write(LEVAL) },

		`(QQR7) '<'`: func() { next(); write(QLEVAL) },

		`(R5, R7) '>'`: func() { next(); write(REVAL) },

		`(QQR7) '>'`: func() { next(); write(QREVAL) },

		`(R5, R7) ','`: func() { next(); write(COMMA) },

		`(R5) '='`: func() { next(); write(REPLACE) },

		`(R7) '='`: func() {
			if next(); r.Code == '>' {
				next()
				write(ARROW)
			} else {
				write(REPLACE)
			}
		},

		`(R5) ':'`: func() { next(); write(COLON) },

		`(R7) ':'`: func() {
			if next(); r.Code == ':' {
				next()
				write(DCOLON)
			} else {
				write(COLON)
			}
		},

		`(R7) '~'`: func() { next(); write(TILDE) },
	},
	)

	bg.bind(map[string]func(){
		`(R7) '$'`: func() {
			scanTail(identTailRunes)
			if tag, ok := keywords_7[strings.ToUpper(cb.String())]; ok {
				write(tag)
			} else {
				errStart("unknown keyword")
			}
		},

		`(R5) '$'`: func() {
			scanTail(identTailRunes)
			if tag, ok := keywords_5[cb.String()]; ok {
				write(tag)
			} else {
				errStart("unknown keyword")
			}
		},
	},
	)

	scanHex := func(width int, msg string) (x rune) {
		next()
		for i := 0; i < width; i++ {
			if d, ok := hexDigits[r.Code]; ok {
				x = x*16 + rune(d)
				next()
			} else {
				err(msg)
				return
			}
		}

		return
	}

	escapes := newBuilder([]string{"R5", "R7"})
	escapes.bind(map[string]func(){
		`(R5, R7) '\''`: func() { cb.save('\''); next() },
		`(R5, R7) '"'`:  func() { cb.save('"'); next() },
		`(R5, R7) '\\'`: func() { cb.save('\\'); next() },
		`(R5, R7) 'n'`:  func() { cb.save('\n'); next() },
		`(R5, R7) 'r'`:  func() { cb.save('\r'); next() },
		`(R5, R7) 't'`:  func() { cb.save('\t'); next() },
		`(R7) 'u'`:      func() { cb.save(scanHex(4, "illegal \\u escape")) },
		`(R7) 'U'`:      func() { cb.save(scanHex(8, "illegal \\U escape")) },
		`(R7) 'a'`:      func() { cb.save('\a'); next() },
		`(R7) 'v'`:      func() { cb.save('\v'); next() },
		`(R7) 'b'`:      func() { cb.save('\b'); next() },
		`(R7) 'f'`:      func() { cb.save('\f'); next() },
		`(R7) ' ', '\n'`: func() {
			for r.Code == ' ' {
				next()
			}

			if r.Code == '\n' {
				next()
			} else if r.Code != chars.EOF {
				err("newline expected after '\\'")
			}
		},
		`(R5) 'x'`: func() { cb.save(scanHex(2, "illegal \\x escape")) },
		`(R5) '('`: func() { cb.save('('); next() },
		`(R5) ')'`: func() { cb.save(')'); next() },
		`(R5) '<'`: func() { cb.save('<'); next() },
		`(R5) '>'`: func() { cb.save('>'); next() },
	},
	)

	bg.bind(map[string]func(){
		`(R5, R7, QR7, QQR7) '\'', '"';` +
			`(R7, QR7, QQR7) '@'`: func() {
			verbatim := false
			if r.Code == '@' {
				verbatim = true
				if next(); r.Code != '\'' && r.Code != '"' {
					err("' or \" expected after @")
					return
				}
			}

			tc := r.Code
			cb.reset()

			next()
			for {
				switch r.Code {
				case tc:
					if verbatim && r_next.Code == tc {
						cb.save(tc)
						next()
						next()
					} else {
						if len(cb) == 0 {
							if tc == '\'' {
								warn("empty string makes no sense")
							} else {
								err("empty compound not allowed")
							}
						}

						if next(); tc == '\'' {
							writeString()
						} else {
							writeCbCompound(false)
						}
						return
					}
				case '\\':
					if verbatim {
						cb.save('\\')
						next()
					} else {
						next()
						if h := escapes.dispatch(r.Code); h != nil {
							h()
						} else {
							err("unknown escape")
							next()
						}
					}
				case '\n':
					if tc == '"' {
						err("newline in compound")
						writeCbCompound(false)
					} else {
						err("newline in string literal")
						writeString()
					}
					return
				case chars.EOF:
					if tc == '"' {
						err("end-of-file in compound")
						writeCbCompound(false)
					} else {
						err("end-of-file in string literal")
						writeString()
					}
					return
				default:
					cb.save(r.Code)
					next()
				}
			}
		},
	},
	)

	scanIdentOrVar := func(anonymousVars bool) {
		t, ok := varTypes_7[r.Code]
		if ok && (r_next.Code == '.' ||
			(identTailRunes.dispatch(r_next.Code) == nil) && anonymousVars) {
			if next(); r.Code == '.' {
				if next(); identTailRunes.dispatch(r.Code) != nil {
					scanTail(identTailRunes)
				} else {
					cb.reset()
					err("error in var")
				}
			} else {
				cb.reset()
			}

			writeVar(t)
		} else {
			scanTail(identTailRunes)
			writeCbCompound(true)
		}
	}

	bg.bind(map[string]func(){
		`(R7, QR7, QQR7) 'A'-'Z', 'a'-'z', '_', letter`: func() { scanIdentOrVar(true) },

		`(QR7, QQR7) "!#$%&*+,:;=?\\^|~";` +
			`(QR7) '<', '>'`: scanOp,

		`(R5) 'A'-'Z', 'a'-'z', '_'`: func() { scanIdentOrVar(false) },
	},
	)

	readDigits := func(digitRunes *builder) (count int) {
		for count = 0; digitRunes.dispatch(r.Code) != nil; count++ {
			cb.save(r.Code)
			next()
		}
		return
	}

	scanNumber := func(minus bool) {
		fp := false
		intRes, floatRes := big.NewInt(0), float64(0.0)

		if r.Code == '0' && (r_next.Code == 'x' || r_next.Code == 'X') {
			// Processing hexadecimal number.
			if minus {
				err("minus not allowed in hex number")
			}

			next()
			next()
			cb.reset()
			readDigits(hexDigitRunes)
			intRes.SetString(cb.String(), 16)
		} else {
			cb.reset()
			readDigits(decDigitRunes)
			if r.Code == '.' || r.Code == 'E' || r.Code == 'e' {
				// Processing floating-point number
				fp = true
				nErr := false

				if r.Code == '.' {
					cb.save('.')
					if next(); readDigits(decDigitRunes) == 0 {
						err("digits expected after dot")
						nErr = true
					}
				}

				if r.Code == 'E' || r.Code == 'e' {
					cb.save('e')
					if next(); r.Code == '+' || r.Code == '-' {
						cb.save(r.Code)
						next()
					}

					if readDigits(decDigitRunes) == 0 {
						err("exponent expected")
						nErr = true
					}
				}

				if !nErr {
					s := cb.String()
					if minus {
						s = "-" + s
					}

					var e error
					if floatRes, e = strconv.ParseFloat(s, 64); e == strconv.ErrRange {
						err("floating-point number overflow")
					}
				}
			} else {
				// Processing decimal number
				intRes.SetString(cb.String(), 10)
				if minus {
					intRes.Neg(intRes)
				}
			}
		}

		if identTailRunes.dispatch(r.Code) != nil {
			err("error in number")
			scanTail(identTailRunes)
		}

		if fp {
			writeFloat(floatRes)
		} else {
			writeInt(intRes)
		}
	}

	bg.bind(map[string]func(){
		`(R7) '0'-'9', '-', '.'`: func() {
			minus := r.Code == '-'
			if minus {
				if next(); r.Code == '>' {
					next()
					write(TARROW)
					return
				}
			}

			scanNumber(minus)
		},

		`(QR7, QQR7) '0'-'9', '-', '.'`: func() {
			minus := r.Code == '-'
			if minus && decDigitRunes.dispatch(r_next.Code) == nil &&
				r_next.Code != '.' {
				scanOp()
			} else {
				if minus {
					next()
				}
				scanNumber(minus)
			}
		},

		`(R5) '0'-'9'`: func() {
			cb.reset()
			readDigits(decDigitRunes)
			x := big.NewInt(0)
			x.SetString(cb.String(), 10)

			if x.Cmp(maxMacrodigit) > 0 {
				err("number greater than 2^32-1=4294967295")
			}

			writeInt(x)
		},
	},
	)

	quoted := false
	var qPos coords.Pos

	bg.bind(map[string]func(){
		"(R7) '`'": func() {
			quoted, qPos = true, r.Pos

			if next(); r.Code == '`' {
				next()
				bg.setMode("QQR7")
				opTailRunes.setMode("QQR7")
			} else {
				bg.setMode("QR7")
				opTailRunes.setMode("QR7")
			}

			write(QUOTE)
		},

		"(QR7, QQR7) '`'": func() {
			bg.setMode("R7")
			next()
			quoted = false
		},
	},
	)

	var mode string
	switch dialect {
	case 7:
		mode = "R7"
	case 5:
		mode = "R5"
	}

	bg.setMode(mode)
	escapes.setMode(mode)
	identTailRunes.setMode(mode)

	for next(); r.Code != chars.EOF; {
		if handler := bg.dispatch(r.Code); handler != nil {
			start = r.Pos
			handler()
		} else {
			errUnexpectedChar()
			next()
		}
	}

	start = r.Pos
	if quoted {
		err("quoting at " + qPos.String() + " not closed")
	}

	write(END_OF_INPUT)
	close(ts)
}

func Intercept(out chan<- Data, w io.WriteCloser, in <-chan Data) {
	tokens := make([]Data, 0, 1024)

	for t := range in {
		out <- t
		tokens = append(tokens, t)
	}

	bs, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(bs)
	w.Close()
	close(out)
}
