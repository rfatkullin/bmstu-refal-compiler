package tokens

import (
	"bytes"
	//"container/vector"
	"fmt"
	"os"
	"text/scanner"
	"strconv"
	"unicode"
)

type table struct {
	ascii  [128]func()
	letter func()
}

func (t *table) bind(c rune, handler func()) {
	if t.ascii[c] != nil {
		fmt.Fprintf(os.Stderr, "internal error in table.bind: %s is already bound\n",
			string([]rune{c}))
	}

	t.ascii[c] = handler
}

func (t *table) bindRange(low, high rune, handler func()) {
	for x := low; x <= high; x++ {
		t.bind(x, handler)
	}
}

type builder struct {
	modes  []string
	cur    *table
	tables map[string]*table
}

func newBuilder(modes []string) (b *builder) {
	b = new(builder)
	b.modes = modes
	b.tables = make(map[string]*table)

	for _, t := range modes {
		b.tables[t] = new(table)
	}

	b.cur = b.tables[modes[0]]
	return
}

func (b *builder) setMode(mode string) {
	b.cur = b.tables[mode]
}

func (b *builder) dispatch(char rune) (handler func()) {
	if char&0x1FFF80 == 0 {
		if handler = b.cur.ascii[char]; handler != nil {
			return
		}
	}

	if b.cur.letter != nil && (unicode.IsLetter(char) || char == '_') {
		handler = b.cur.letter
	}
	return
}

func (b *builder) bind(script map[string]func()) {
	var sc scanner.Scanner
	var modes []string
	for targets, handler := range script {
		sc.Init(bytes.NewBufferString(targets))

		err := func(msg string) {
			fmt.Fprintf(os.Stderr, "internal error in builder.bind: %s in %q at %v\n",
				msg, targets, sc.Pos())
		}

		expect := func(tok rune) {
			if sc.Scan() != tok {
				err(fmt.Sprintf("expected %q", string([]rune{tok})))
			}
		}

	L3:
		for {
			// Reading list of modes.
			modes = modes[:0]
			expect('(')
		L1:
			for {
				if sc.Scan() != scanner.Ident {
					err("identifier expected")
					return
				}
				modes = append(modes, sc.TokenText())
				switch sc.Scan() {
				case ',': // missing comma
				case ')':
					break L1
				default:
					err("expected ',' or ')'")
					return
				}
			}

			// Reading characters
		L2:
			for {
				var tok rune
				switch sc.Scan() {
				case scanner.Ident:
					if sc.TokenText() == "letter" {
						for _, m := range modes {
							b.tables[m].letter = handler
						}
						tok = sc.Scan()
					} else {
						err("unknown keyword")
						return
					}
				case scanner.Char:
					c1, _, _, _ := strconv.UnquoteChar(sc.TokenText()[1:], '\'')
					if tok = sc.Scan(); tok == '-' {
						if sc.Scan() == scanner.Char {
							c2, _, _, _ := strconv.UnquoteChar(sc.TokenText()[1:], '\'')
							for _, m := range modes {
								b.tables[m].bindRange(c1, c2, handler)
							}
							tok = sc.Scan()
						} else {
							err("expected character literal")
						}
					} else {
						for _, m := range modes {
							b.tables[m].bind(c1, handler)
						}
					}
				case scanner.String:
					s, _ := strconv.Unquote(sc.TokenText())
					tok = sc.Scan()
					for _, c := range []rune(s) {
						for _, m := range modes {
							b.tables[m].bind(c, handler)
						}
					}
				default:
					err("unexpected token")
					return
				}

				switch tok {
				case scanner.EOF:
					break L3
				case ';':
					break L2
				case ',': // missing comma
				default:
					err("expected ','")
					return
				}
			}
		}
	}
}
