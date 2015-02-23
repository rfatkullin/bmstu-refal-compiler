// Bauman Refal Compiler message processing package
package messages

import (
	"fmt"
	"os"
	"sort"
)

import "bmstu-refal-compiler/coords"

type Importance int

const (
	RECOVERY_WARNING Importance = iota
	WARNING
	ERROR
)

var importances []string = []string{"recovery warning", "warning", "error"}

func (x Importance) String() string {
	return importances[x]
}

type Data struct {
	coords.Pos
	Importance
	Msg string
}

type DataArray []Data

func (a DataArray) Len() int { return len(a) }

func (a DataArray) Less(i, j int) bool {
	return a[i].Offs < a[j].Offs
}

func (a DataArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Summary struct {
	DataArray
	Stat map[Importance]int
}

const initialCapacity = 16

func NewSummary() *Summary {
	return &Summary{make(DataArray, 0, initialCapacity), make(map[Importance]int)}
}

func (s *Summary) Put(d Data) {
	s.DataArray = append(s.DataArray, d)
	n, _ := s.Stat[d.Importance]
	s.Stat[d.Importance] = n + 1
}

func (s Summary) Print(fileName string) {
	for _, x := range s.DataArray {
		fmt.Fprintf(os.Stderr, "%s: %v: %v: %s\n",
			fileName, x.Pos, x.Importance, x.Msg)
	}
}

func Handle(ms <-chan Data, result chan<- Summary) {
	s := NewSummary()
	for d := range ms {
		s.Put(d)
	}

	if s.Len() > 0 {
		sort.Sort(s)

		// Removing duplicates.
		i := 0
		for {
			length := s.Len() - 1
			if i == length {
				break
			}

			m1, m2 := s.DataArray[i], s.DataArray[i+1]
			if m1.Offs == m2.Offs && m1.Importance == m2.Importance &&
				m1.Msg == m2.Msg {
				copy(s.DataArray[i:length], s.DataArray[i+1:])
				s.DataArray = s.DataArray[0:length]
			} else {
				i++
			}
		}
	}

	result <- *s
}
