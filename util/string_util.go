package util

import (
	"bytes"
	"fmt"
	"sort"
)

type StringArray []string

// Find the position of target in a string array
func (ar *StringArray) IndexOf(target string) int {
	for i, s := range *ar {
		if s == target {
			return i
		}
	}
	return -1
}

func (ar *StringArray) Compare(that StringArray) bool {
	if len(*ar) != len(that) {
		return false
	}
	for i := 0; i < len(*ar); i++ {
		if (*ar)[i] != (that)[i] {
			return false
		}
	}
	return true
}

func (ar StringArray) ToString() string {
	var buf bytes.Buffer
	for i, s := range ar {
		buf.WriteString(fmt.Sprintf("%2d. ", i))
		buf.WriteString(s)
		buf.WriteRune('\n')
	}
	return buf.String()
}

type SortedStringArray []string

// Insert s into the array, and keep it sorted.
// Return the reference to a new sorted array
func (ar SortedStringArray) Insert(s string) *SortedStringArray {
	var sorted SortedStringArray
	if len(ar) == 0 {
		sorted = []string{s}
		return &sorted
	}
	pos := sort.SearchStrings(ar, s)
	if pos == 0 {
		sorted = append([]string{s}, ar...)
	} else if len(ar) == pos {
		sorted = append(ar, s)
	} else {
		after := (ar)[pos:]
		sorted = append((ar)[0:pos], append([]string{s}, after...)...)
	}
	return &sorted
}

func (ar *SortedStringArray) IndexOf(s string) int {
	if len(*ar) == 0 {
		return -1
	}
	pos := sort.SearchStrings(*ar, s)
	if pos >= len(*ar) {
		return -1
	}
	if (*ar)[pos] == s {
		return pos
	}
	return -1
}

type StringSet SortedStringArray

// Add a new string to set.
func (set StringSet) Add(s string) *StringSet {
	var ar = SortedStringArray(set)
	if ar.IndexOf(s) == -1 {
		set2 := StringSet(*ar.Insert(s))
		return &set2
	}
	return &set
}

// Check whether 2 sets are the same.
func (set *StringSet) Compare(that StringSet) bool {
	ar1 := StringArray(*set)
	ar2 := StringArray(that)
	return ar1.Compare(ar2)
}

func (set *StringSet) IndexOf(s string) int {
	sorted := SortedStringArray(*set)
	return sorted.IndexOf(s)
}

func (set StringSet) ToString() string {
	var ar = StringArray(set)
	return ar.ToString()
}
