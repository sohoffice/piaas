package piaas

import (
	"log"
	"testing"
)

func TestRsyncMatcher(t *testing.T) {
	// simple matcher
	p1 := NewRsyncPattern("foo")
	runTests(t, p1, []string{"/foo/123", "/123/foo", "/123/foo/123"}, []string{"/", "/12fo3o", "/foo123"})

	// absolute matcher
	p2 := NewRsyncPattern("/foo")
	runTests(t, p2, []string{"/foo/123", "/foo"}, []string{"/", "/foo123", "/123/foo", "/123/foo/456"})

	// sub directory matcher
	p3 := NewRsyncPattern("foo/")
	runTests(t, p3, []string{"/foo/", "/123/foo/"}, []string{"/", "/123", "/123/456", "123foo/"})

	// partial filename matcher
	p4 := NewRsyncPattern("foo*")
	runTests(t, p4, []string{"/foo", "/foo/bar", "/foo123", "/foo123/bar"}, []string{"/", "/123", "/123foo", "/123foobar", "/bar/fobazo"})

	p5 := NewRsyncPattern("foo**")
	runTests(t, p5, []string{"/foo", "/foo123", "/foo/abc", "/abc/foo", "/abc/foo123", "/abc/foo/bar"}, []string{"/", "/123", "/fobaro"})
}

func runTests(t *testing.T, pat RsyncPattern, positives []string, negatives []string) {
	log.Printf("Run tests of %s.", pat.ToString())
	for _, s := range positives {
		if !pat.Match(s) {
			t.Errorf("%s should match %s.", pat.ToString(), s)
		}
	}

	for _, s := range negatives {
		if pat.Match(s) {
			t.Errorf("%s should not match %s.", pat.ToString(), s)
		}
	}
}
