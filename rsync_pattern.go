package piaas

import (
	"fmt"
	"github.com/sohoffice/piaas/util"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

type RsyncPattern struct {
	// the original pattern
	orig string

	sanitized string
	// the sanitized pattern compiled into regexp
	pattern *regexp.Regexp
	// if  the  pattern  starts  with  a  / then it is anchored to a particular spot in the hierarchy of files
	anchored bool
	// if the pattern ends with a / then it will only match a directory
	directoryOnly bool
	// if the pattern contains one of  these wildcard characters: '*', '?'.
	hasWildcard bool
	// should match against the full pathname recursively
	isFull bool
}

func (rp RsyncPattern) Match(path string) bool {
	list := splitFilename(path, &util.StringArray{})
	return rp.matchParts(list, 0)
}

// Match the filename parts from the beginning to `endIndex`
func (rp *RsyncPattern) matchParts(parts *util.StringArray, endIndex int) bool {
	// extra precaution to make sure the endIndex do not grow beyond the parts length
	if endIndex >= len(*parts) {
		return false
	}
	var flag bool
	// if we're in the full match mode, join the path to match
	if rp.isFull {
		full := filepath.Join((*parts)[0 : endIndex+1]...)
		flag = rp.pattern.Match([]byte(full))
	} else {
		// otherwise match the element pointed by endIndex
		flag = rp.pattern.Match([]byte((*parts)[endIndex]))
	}

	// the pattern matches, or we have reached the end of the list
	if flag || endIndex >= len(*parts) {
		return flag
	} else {
		return rp.matchParts(parts, endIndex+1)
	}
}

func (rp RsyncPattern) ToString() string {
	return fmt.Sprintf("RsyncPattern[Orig: %s, Sanitized: %s]", rp.orig, rp.sanitized)
}

// question mark replaced by any random character
var patQuestionMark = regexp.MustCompile("\\?")
var sanitizeQuestionMark = func(bytes []byte) []byte {
	return patQuestionMark.ReplaceAll(bytes, []byte("[^/]?"))
}

// ** to any string and sub path
var patDoubleAsterisks = regexp.MustCompile("\\*\\*")
var sanitizeDoubleAsterisks = func(bytes []byte) []byte {
	return patDoubleAsterisks.ReplaceAll(bytes, []byte(".*"))
}

// * to any string but not sub path
var patAsterisk = regexp.MustCompile("([^.])\\*")
var sanitizeAsterisks = func(bytes []byte) []byte {
	return patAsterisk.ReplaceAllFunc(bytes, func(ar []byte) []byte {
		s := string(ar)
		s = s[0:1] + "[^/]*"
		return []byte(s)
	})
}

// escaping special characters
var patEscape = regexp.MustCompile("[.\\\\(){}|+^$]")
var sanitizeEscape = func(bytes []byte) []byte {
	return patEscape.ReplaceAllFunc(bytes, func(ar []byte) []byte {
		return append([]byte("\\"), ar...)
	})
}

// Collect all pattern sanitize rules
var sanitizers = []func([]byte) []byte{sanitizeEscape, sanitizeQuestionMark, sanitizeDoubleAsterisks, sanitizeAsterisks}

var anchorTest = regexp.MustCompile("^/.*")
var directoryTest = regexp.MustCompile(".*/$")
var wildcardTest = regexp.MustCompile(".*[?*].*")

// if the pattern contains a / (not counting a trailing /) or a "**"
var fullTest = regexp.MustCompile(".*(/|\\*\\*)[^/]+")

// Create a new RsyncPattern.
func NewRsyncPattern(pat string) RsyncPattern {
	bytes := []byte(filepath.ToSlash(pat))
	for _, san := range sanitizers {
		bytes = san(bytes)
	}
	sanitized := "^" + string(bytes) + "$"
	reg, err := regexp.Compile(sanitized)
	if err != nil {
		log.Fatalf("Error compiling pattern %s: %s", string(bytes), err)
	}
	return RsyncPattern{
		orig:      pat,
		sanitized: sanitized,
		pattern:   reg,
		// the below should be checked against the original pattern.
		anchored:      anchorTest.MatchString(pat),
		directoryOnly: directoryTest.MatchString(pat),
		hasWildcard:   wildcardTest.MatchString(pat),
		isFull:        fullTest.MatchString(pat),
	}
}

type RsyncPatterns struct {
	// all registered patterns
	patterns []RsyncPattern

	// The base directory of the pattern.
	// patterns will only contain path relative to the basedir.
	basedir string

	// the OS path separator == filepath.Separator
	separator string
}

type RsyncMatcher struct {
	path  string
	parts []string
}

// Create a new RsyncMatcher
func NewRsyncMatcher(path string) RsyncMatcher {

	return RsyncMatcher{
		path: path, parts: filepath.SplitList(path),
	}
}

// Build the convert rule to standardize the path separator to /
func convertSeparatorRule(sep string) func(string) string {
	if sep != "/" {
		// convert the separator, if the separator is not /
		return func(s string) string {
			return strings.Replace(s, sep, "/", -1)
		}
	} else {
		return func(s string) string {
			return s
		}
	}
}

// Split the filename into parts, each part represent one level in the directory hierarchy.
//
// The split is done using filepath.Split and be advised the separator will be part of the directory element.
//
// Ex: /foo/bar => ["/", "foo/", "bar"]
// The "/" and "foo/" are directories and they all have a trailing slash to indicate this is a directory.
func splitFilename(path string, splitted *util.StringArray) *util.StringArray {
	if path == "" {
		return splitted
	}
	var isDir bool
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
		isDir = true
	}
	dir, file := filepath.Split(path)
	var list util.StringArray
	if isDir {
		list = append(*splitted, file+"/")
	} else {
		list = append(*splitted, file)
	}
	if strings.HasSuffix(dir, "/") {
		dir = dir[0 : len(dir)-1]
		list = append(list, "/")
	}
	if dir == "" {
		ar2 := util.StringArray(list)
		return ar2.Reverse()
	} else {
		return splitFilename(dir, &list)
	}
}
