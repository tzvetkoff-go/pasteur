package stringutil

import (
	"regexp"
	"strconv"
	"strings"
)

// LeadingIdentRegExp ...
var LeadingIdentRegExp = regexp.MustCompile("[\r\n]*([ \t]+)")

// FormatQuery ...
func FormatQuery(s string) string {
	result := ""
	insideStringApos := false
	insideStringQuot := false
	prevSpace := false

	s = strings.Trim(s, " \t\r\n")

	for _, ch := range s {
		switch ch {
		case '\'':
			prevSpace = false
			if !insideStringQuot {
				insideStringApos = !insideStringApos
			}
			result += string(ch)
		case '"':
			prevSpace = false
			if !insideStringApos {
				insideStringQuot = !insideStringQuot
			}
			result += string(ch)
		case ' ':
			fallthrough
		case '	':
			fallthrough
		case '\n':
			if prevSpace {
				continue
			}
			if !insideStringApos && !insideStringQuot {
				prevSpace = true
				result += " "
			} else {
				result += string(ch)
			}
		default:
			prevSpace = false
			result += string(ch)
		}
	}

	return result
}

// ParseInt ...
func ParseInt(x interface{}) int {
	switch x := x.(type) {
	case int:
		return x
	case int16:
		return int(x)
	case uint16:
		return int(x)
	case int32:
		return int(x)
	case uint32:
		return int(x)
	case int64:
		return int(x)
	case uint64:
		return int(x)
	case string:
		i, err := strconv.ParseInt(x, 0, 64)
		if err == nil {
			return int(i)
		}
	}

	return 0
}
