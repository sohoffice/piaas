package util

import (
	"testing"
)

func TestStringArray_IndexOf(t *testing.T) {
	ar := StringArray([]string{"aa", "bb", "cc"})
	if ar.IndexOf("z") != -1 {
		t.Error("'z' should not be found.")
	}
	if ar.IndexOf("bb") != 1 {
		t.Error("'bb' should be found at position 1.")
	}
}

func TestStringArray_Compare(t *testing.T) {
	var ar1 StringArray = []string{"a", "b"}
	var ar2 StringArray = []string{"a", "b", "c"}
	var ar3 StringArray = []string{"0", "a", "b"}
	var ar4 StringArray = []string{"a", "b"}
	if ar1.Compare(ar2) {
		t.Errorf("%s and %s should not be the same", ar1, ar2)
	}
	if ar1.Compare(ar3) {
		t.Errorf("%s and %s should not be the same", ar1, ar3)
	}
	if !ar1.Compare(ar4) {
		t.Errorf("%s and %s should be the same", ar1, ar4)
	}
	if !ar1.Compare(ar1) {
		t.Errorf("%s and %s should be the same", ar1, ar1)
	}
}

func TestSortedStringArray_Insert(t *testing.T) {
	var ar SortedStringArray = []string{"a", "c"}
	ar1 := StringArray(*ar.Insert("0"))
	if !ar1.Compare([]string{"0", "a", "c"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"0", "a", "c"}, ar1)
	}
	ar2 := StringArray(*ar.Insert("b"))
	if !ar2.Compare([]string{"a", "b", "c"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"a", "b", "c"}, ar2)
	}
	ar3 := StringArray(*ar.Insert("z"))
	if !ar3.Compare([]string{"a", "c", "z"}) {
		t.Errorf("Not the same array.\nexpected: %s\nactual: %s", []string{"a", "c", "z"}, ar3)
	}
	ar4 := StringArray(*ar.Insert("c"))
	if !ar4.Compare([]string{"a", "c", "c"}) {
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
