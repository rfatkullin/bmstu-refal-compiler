package emitter

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	tab = "\t"
)

func genTabs(depth int) string {
	return strings.Repeat(tab, depth)
}

func (f *Data) printLabel(depth int, label string) {
	tabs := genTabs(depth)
	fmt.Fprintf(f, "%s%s\n", tabs, label)
}

/// Unicode string --> char string"
func getStrOfRunes(str string) string {

	runes := make([]string, 0)

	for _, rune := range str {
		runes = append(runes, fmt.Sprintf("%d", rune))
	}

	return strings.Join(runes, ",")
}

/// Bytes --> string
func getStrOfBytes(num *big.Int) (str string, sign int, count int) {
	strs := make([]string, 0)
	bytes := num.Bytes()

	if len(bytes) == 0 {
		return "0", 0, 1
	}

	for _, byteVal := range num.Bytes() {
		strs = append(strs, fmt.Sprintf("%d", byteVal))
	}

	sign = 0
	if num.Sign() < 0 {
		sign = 1
	}

	str = strings.Join(strs, ",")
	count = len(bytes)

	return
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
