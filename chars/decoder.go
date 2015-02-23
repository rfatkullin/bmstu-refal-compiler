// Bauman Refal Compiler input handling package
package chars

import "errors"

import (
	"io"
	"strconv"
)

import (
	"bmstu-refal-compiler/coords"
	"bmstu-refal-compiler/messages"
)

const EOF = -1

type Rune struct {
	coords.Pos
	Code rune
}

func Handle(runes chan<- Rune, ms chan<- messages.Data, r io.ReadCloser, cp int) {
	const (
		bufReserve = 4
		bufSize    = 1024
		bufTotal   = bufReserve + bufSize
	)

	var (
		buf   [bufTotal]byte
		start int = -bufTotal
		index int = bufTotal
		line  int = 1
		col   int = 1
		n     int
		err   error
	)

	complain := func(i int, s string) {
		ms <- messages.Data{
			coords.Pos{start + i, 0, 0},
			messages.ERROR, s}
	}

	complainI := func(s string) { complain(index, s) }

	read := func() {
		start += bufSize
		index -= bufSize

		copy(buf[0:bufReserve], buf[bufSize:bufTotal])
		n, err = io.ReadFull(r, buf[bufReserve:bufTotal])

		if err == io.ErrUnexpectedEOF {
			err = errors.New("EOF")
		}

		if err != nil && err.Error() != "EOF" {
			complain(bufReserve+n, "i/o error: "+err.Error())
		}
	}

	align := func(s int) {
		if r := n % s; n != bufSize && r != 0 {
			if err.Error() == "EOF" {
				complain(bufReserve+n,
					"file size is not multiple of "+strconv.Itoa(s))
			}

			n -= r
		}
	}

	crFlag := false

	write := func(c rune, fix int) {
		pos := coords.Pos{start + index - fix, line, col}
		switch c {
		case '\r':
			runes <- Rune{pos, '\n'}
			line++
			col, crFlag = 1, true
		case '\n':
			if crFlag {
				break
			}
			fallthrough
		case 0x2028 /* Unicode Line Separator */, 0x2029 /* Unicode Paragraph Separator */ :
			runes <- Rune{pos, '\n'}
			line++
			col, crFlag = 1, false
		default:
			runes <- Rune{pos, c}
			col++
			crFlag = false
		}
	}

	read()

	var cpIndex int
	switch c := buf[index:]; {
	case n >= 4 && c[0] == 0 && c[1] == 0 && c[2] == 0xFE && c[3] == 0xFF:
		cpIndex = utf32be
		index += 4
	case n >= 4 && c[0] == 0xFF && c[1] == 0xFE && c[2] == 0 && c[3] == 0:
		cpIndex = utf32le
		index += 4
	case n >= 3 && c[0] == 0xEF && c[1] == 0xBB && c[2] == 0xBF:
		cpIndex = utf8
		index += 3
	case n >= 2 && c[0] == 0xFE && c[1] == 0xFF:
		cpIndex = utf16be
		index += 2
	case n >= 2 && c[0] == 0xFF && c[1] == 0xFE:
		cpIndex = utf16le
		index += 2
	default:
		cpIndex = cpMap[cp]
	}

	var decode func() (r uint32)

	switch cpIndex {
	case utf8:
		flag := true

		complainF := func(s string) {
			if flag {
				complainI(s)
				flag = false
			}
		}

		complain8 := func() { complainF("corrupted UTF-8 code point") }

		for {
			if index >= bufTotal-3 {
				if err == nil {
					read()
				} else {
					complainF("end of file breaks code point")
					break
				}
			}

			tail := bufReserve + n - index

			if tail == 0 {
				break
			}
			c0 := buf[index]

			// 1-byte, 7-bit sequence?
			if c0 < 0x80 /* 1000 0000 */ {
				write(rune(c0), 0)
				index++
				flag = true
				continue
			}

			if c0 < 0xC0 /* 1100 0000 */ {
				complainF("unexpected UTF-8 continuation byte")
				continue
			}

			// need first continuation byte
			if tail == 1 {
				break
			}
			index++
			c1 := buf[index]

			if c1 < 0x80 /* 1000 0000 */ || c1 >= 0xC0 /* 1100 0000 */ {
				complain8()
				continue
			}

			// 2-byte, 11-bit sequence?
			if c0 < 0xE0 /* 1110 0000 */ {
				x := rune(c0&0x3F /* 0011 1111 */)<<6 |
					rune(c1&0x3F /* 0011 1111 */)
				if x <= (1<<7)-1 {
					complain8()
				} else {
					write(x, 1)
					index++
					flag = true
				}
				continue
			}

			// need second continuation byte
			if tail == 2 {
				break
			}
			index++
			c2 := buf[index]

			if c2 < 0x80 /* 1000 0000 */ || c2 >= 0xC0 /* 1100 0000 */ {
				complain8()
				continue
			}

			// 3-byte, 16-bit sequence?
			if c0 < 0xF0 /* 1111 0000 */ {
				x := rune(c0&0x1F /* 0001 1111 */)<<12 |
					rune(c1&0x3F /* 0011 1111 */)<<6 |
					rune(c2&0x3F /* 0011 1111 */)
				if x <= (1<<11)-1 {
					complain8()
				} else {
					write(x, 2)
					index++
					flag = true
				}
				continue
			}

			// need third continuation byte
			if tail == 3 {
				break
			}
			index++
			c3 := buf[index]

			if c3 < 0x80 /* 1000 0000 */ || c3 >= 0xC0 /* 1100 0000 */ {
				complain8()
				continue
			}

			// 4-byte, 21-bit sequence?
			if c0 < 0xF8 /* 1111 1000 */ {
				x := rune(c0&0x07 /* 0000 0111 */)<<18 |
					rune(c1&0x3F /* 0011 1111 */)<<12 |
					rune(c2&0x3F /* 0011 1111 */)<<6 |
					rune(c3&0x3F /* 0011 1111 */)
				if x <= (1<<16)-1 {
					complain8()
				} else {
					write(x, 3)
					index++
					flag = true
				}
				continue
			}

			// error
			complain8()
		}

	case utf16le:
		decode = func() (r uint32) {
			r = (uint32(buf[index]) | uint32(buf[index+1])<<8)
			index += 2
			return
		}
		fallthrough

	case utf16be:
		if decode == nil {
			decode = func() (r uint32) {
				r = (uint32(buf[index]) << 8) | uint32(buf[index+1])
				index += 2
				return
			}
		}

		align(2)

		for {
			if index == bufReserve+n {
				if index != bufTotal || err != nil {
					break
				}
				read()
				align(2)
			}

			x := decode()
			if x >= 0xD800 && x < 0xE000 {
				/* x is in the surrogate range. */
				if x >= 0xDC00 {
					complainI("surrogate pair cannot begin with low surrogate")
					continue
				}

				if index == bufReserve+n {
					if index != bufTotal || err.Error() == "EOF" {
						complainI("end of file breaks surrogate pair")
						continue
					}
					read()
					align(2)
				}

				y := decode()
				if y < 0xDC00 || y >= 0xE000 {
					complainI("low surrogate expected")
					continue
				}

				write(rune((x-0xD800)*0x400+(y-0xDC00)+0x10000), 4)
			} else {
				write(rune(x), 2)
			}
		}

	case utf32le:
		decode = func() (r uint32) {
			a := buf[index : index+4]
			index += 4
			r = (((((uint32(a[3]) << 8) | uint32(a[2])) << 8) |
				uint32(a[1])) << 8) | uint32(a[0])
			return
		}
		fallthrough

	case utf32be:
		if decode == nil {
			decode = func() (r uint32) {
				a := buf[index : index+4]
				index += 4
				r = (((((uint32(a[0]) << 8) | uint32(a[1])) << 8) |
					uint32(a[2])) << 8) | uint32(a[3])
				return
			}
		}

		align(4)

		for {
			if index == bufReserve+n {
				if index != bufTotal || err != nil {
					break
				}
				read()
				align(4)
			}

			switch x := decode(); {
			case x > 0x10FFFF:
				complainI("corrupted UTF-32 code point")
			case x >= 0xD800 && x < 0xE000:
				complainI("surrogates not allowed in UTF-32")
			default:
				write(rune(x), 4)
			}
		}

	default:
		var t [256]rune
		offs := cpOffs[cpIndex]
		ucpSlice := ucpData[offs : offs+256]

		for i := 0; i < 256; i++ {
			t[i] = rune(ucpSlice[i])
		}

		chPos := cpChindex[cpIndex]
		if chPos != 0 {
			t[chPos] = rune(cpChval[cpIndex])
		}

		for {
			for ; index < bufReserve+n; index++ {
				write(t[buf[index]], 0)
			}

			if err != nil {
				break
			}
			read()
		}
	}

	write(EOF, 0)
	r.Close()
	close(runes)
}
