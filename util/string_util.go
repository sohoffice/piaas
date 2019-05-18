package util

import (
	"github.com/sohoffice/piaas/stringarrays"
	"sort"
)

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
	return stringarrays.Compare(*set, that)
}

func (set *StringSet) IndexOf(s string) int {
	sorted := SortedStringArray(*set)
	return sorted.IndexOf(s)
}

func (set StringSet) ToString() string {
	return stringarrays.ToString(set)
}
