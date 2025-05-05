package spender

import "regexp"

type set struct{}

var (
	trimBegin *regexp.Regexp = regexp.MustCompile(
		`(?i)(online\s+payment\s+\d+\s+to\s+|pwp|tst)\*?\s*`,
	)
	trimEnd *regexp.Regexp = regexp.MustCompile(
		`(?i)\s+(privacycom\s+tn:\s+\d+|(ppd|web)\s+id:\s+\S+)`,
	)
)

// Version is the package version.
const Version string = "0.2.2"
