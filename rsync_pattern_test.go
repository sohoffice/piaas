package piaas

import (
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRsyncMatcher(t *testing.T) {
	// simple matcher
	p1 := NewRsyncPattern("foo")
	runRsyncMatcherTests(t, p1, []string{"/foo/123", "/123/foo", "/123/foo/123"}, []string{"/", "/12fo3o", "/foo123"})

	// absolute matcher
	p2 := NewRsyncPattern("/foo")
	runRsyncMatcherTests(t, p2, []string{"/foo/123", "/foo"}, []string{"/", "/foo123", "/123/foo", "/123/foo/456"})

	// sub directory matcher
	p3 := NewRsyncPattern("foo/")
	runRsyncMatcherTests(t, p3, []string{"/foo/", "/123/foo/"}, []string{"/", "/123", "/123/456", "123foo/"})

	// partial filename matcher
	p4 := NewRsyncPattern("foo*")
	runRsyncMatcherTests(t, p4, []string{"/foo", "/foo/bar", "/foo123", "/foo123/bar"}, []string{"/", "/123", "/123foo", "/123foobar", "/bar/fobazo"})

	// partial filename and sub directories
	p5 := NewRsyncPattern("foo**")
	runRsyncMatcherTests(t, p5, []string{"/foo", "/foo123", "/foo/abc", "/abc/foo", "/abc/foo123", "/abc/foo/bar"}, []string{"/", "/123", "/fobaro"})

	// multiple levels
	p6 := NewRsyncPattern("foo/bar")
	runRsyncMatcherTests(t, p6, []string{"/foo/bar", "/123/foo/bar", "/123/foo/bar/456"}, []string{"/foo/123/bar", "/foo/bar123", "/123/foo/456/bar"})
}

func runRsyncMatcherTests(t *testing.T, pat RsyncPattern, positives []string, negatives []string) {
	log.Debugf("Run tests of %s.", pat.ToString())
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

func TestRsyncPatterns(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "rsync-pattern-test")
	util.CheckFatal("Temp dir", err)
	defer os.RemoveAll(tempDir)

	p1 := NewRsyncPattern("foo")
	p2 := NewRsyncPattern("*.bak")
	rp := NewRsyncPatterns(tempDir, p1, p2)

	runRsyncPatternsTests(t, rp,
		[]string{path.Join(tempDir, "foo"), path.Join(tempDir, "bar", "foo"), path.Join(tempDir, "file.bak")},
		[]string{"/tmp/abc", tempDir + "zzz", path.Join(tempDir, "foobar")})
}

func runRsyncPatternsTests(t *testing.T, rp RsyncPatterns, positives []string, negatives []string) {
	for _, s := range positives {
		if rp.Match(s) != true {
			t.Errorf("%s should match by patterns.", s)
		}
	}
	for _, s := range negatives {
		if rp.Match(s) != false {
			t.Errorf("%s should not match by patterns.", s)
		}
	}
}
