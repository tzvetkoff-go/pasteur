package stringutil

import (
	"regexp"
	"strings"
)

// LeadingIdentRegExp ...
var LeadingIdentRegExp = regexp.MustCompile("[\r\n]*([ \t]+)")

// TrimHeredoc ...
func TrimHeredoc(s string) string {
	leadingIndentMatches := LeadingIdentRegExp.FindStringSubmatch(s)
	if len(leadingIndentMatches) != 2 {
		return s
	}
	leadingIndent := leadingIndentMatches[1]

	return strings.TrimSpace(strings.Replace(s, "\n"+leadingIndent, "\n", -1))
}
