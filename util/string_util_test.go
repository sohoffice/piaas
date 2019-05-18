package util

import (
	"github.com/sohoffice/piaas/stringarrays"
	"testing"
)

func TestStringArray_IndexOf(t *testing.T) {
	ar := []string{"aa", "bb", "cc"}
	if stringarrays.IndexOf(ar, "z") != -1 {
		t.Error("'z' should not be found.")
	}
	if stringarrays.IndexOf(ar, "bb") != 1 {
		t.Error("'bb' should be found at position 1.")
	}
}

func TestStringArray_Compare(t *testing.T) {
	var ar1 = []string{"a", "b"}
	var ar2 = []string{"a", "b", "c"}
	var ar3 = []string{"0", "a", "b"}
	var ar4 = []string{"a", "b"}
	if stringarrays.Compare(ar1, ar2) {
		t.Errorf("%s and %s should not be the same", ar1, ar2)
	}
	if stringarrays.Compare(ar1, ar3) {
		t.Errorf("%s and %s should not be the same", ar1, ar3)
	}
	if !stringarrays.Compare(ar1, ar4) {
		t.Errorf("%s and %s should be the same", ar1, ar4)
	}
	if !stringarrays.Compare(ar1, ar1) {
		t.Errorf("%s and %s should be the same", ar1, ar1)
	}
}

func TestSortedStringArray_Insert(t *testing.T) {
	var ar SortedStringArray = []string{"a", "c"}
	ar1 := *ar.Insert("0")
	if !stringarrays.Compare(ar1, []string{"0", "a", "c"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"0", "a", "c"}, ar1)
	}
	ar2 := *ar.Insert("b")
	if !stringarrays.Compare(ar2, []string{"a", "b", "c"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"a", "b", "c"}, ar2)
	}
	ar3 := *ar.Insert("z")
	if !stringarrays.Compare(ar3, []string{"a", "c", "z"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"a", "c", "z"}, ar3)
	}
	ar4 := *ar.Insert("c")
	if !stringarrays.Compare(ar4, []string{"a", "c", "c"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"a", "c", "c"}, ar4)
	}
}

func TestSortedStringArray_IndexOf(t *testing.T) {
	var ar SortedStringArray = []string{"a", "b", "c"}
	if p := ar.IndexOf("a"); p != 0 {
		t.Errorf("Expected to find '%s' at position %d, but get %d.", "a", 0, p)
	}
	if p := ar.IndexOf("b"); p != 1 {
		t.Errorf("Expected to find '%s' at position %d, but get %d.", "b", 1, p)
	}
	if p := ar.IndexOf("c"); p != 2 {
		t.Errorf("Expected to find '%s' at position %d, but get %d.", "c", 2, p)
	}
	if p := ar.IndexOf("0"); p != -1 {
		t.Errorf("Expected to find '%s' at position %d, but get %d.", "0", -1, p)
	}
	if p := ar.IndexOf("z"); p != -1 {
		t.Errorf("Expected to find '%s' at position %d, but get %d.", "z", -1, p)
	}
}

func TestStringSet_Add(t *testing.T) {
	var set StringSet = make([]string, 0)
	var set1 StringSet
	set1 = *set.Add("a")
	if !set1.Compare(StringSet{"a"}) {
		t.Errorf("Error adding '%s'", "a")
	}

	set1 = *set1.Add("b")
	if !set1.Compare(StringSet{"a", "b"}) {
		t.Errorf("Error adding '%s'", "a, b")
	}

	set1 = *set1.Add("a")
	if !set1.Compare(StringSet{"a", "b"}) {
		t.Errorf("Should not duplicate 'a', got %s", set1)
	}
}
