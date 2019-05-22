package stringarrays

import (
	"flag"
	"os"
	"testing"
)

// The method to setup and tear down the tests of this package.
func TestMain(m *testing.M) {
	flag.Parse()
	// flag.Lookup("logtostderr").Value.Set("true")
	os.Exit(m.Run())
}

func TestIndexOf(t *testing.T) {
	check := func(ar []string, el string, expected int) {
		pos := IndexOf(ar, el)
		if pos != expected {
			t.Errorf("Expected to find %s at pos %d, but is %d.", el, expected, pos)
		}
	}

	check([]string{"a", "b", "c"}, "a", 0)
	check([]string{"a", "b", "c"}, "b", 1)
	check([]string{"a", "b", "c"}, "z", -1)
}

func TestCompare(t *testing.T) {
	pred := Compare
	verifyTrue(t, []string{"a", "b"}, []string{"a", "b"}, pred)
	verifyFalse(t, []string{"a", "b"}, []string{"b", "a"}, pred)
	if Compare([]string{"a", "b"}, []string{"a", "b"}) == false {
		t.Errorf("['a', 'b'] should match.")
	}
	if Compare([]string{"a", "b"}, []string{"b", "a"}) == true {
		t.Errorf("['a', 'b']")
	}
}

func verifyTrue(t *testing.T, ar1 []string, ar2 []string, pred func([]string, []string) bool) {
	if pred(ar1, ar2) == false {
		t.Errorf("%s and %s should match.", ToString(ar1), ToString(ar2))
	}
}

func verifyFalse(t *testing.T, ar1 []string, ar2 []string, pred func([]string, []string) bool) {
	if pred(ar1, ar2) == true {
		t.Errorf("%s and %s should not match.", ToString(ar1), ToString(ar2))
	}
}

func TestReverse(t *testing.T) {
	checkReversed := func(ar []string, expected []string) {
		rev := Reverse(ar)
		if Compare(rev, expected) == false {
			t.Errorf("Unexpected reverse result.\nExpected: %s\n  Actual: %s", ToString(expected), ToString(rev))
		}
	}

	checkReversed([]string{"a", "b", "c"}, []string{"c", "b", "a"})
}
