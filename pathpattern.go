package cyam

import (
	"fmt"
	"regexp"
	"strings"
)

type PathPattern struct {
	regex  *regexp.Regexp
}

func (p PathPattern) String() string {
	return p.regex.String()
}

func (p PathPattern) Match(k string) bool {
	// Check we have a full match over entire string
	s := p.regex.FindStringIndex(k)
	return len(s) > 1 && s[1] - s[0] == len(k)
}

const doubleStartToken = "${__DOUBLE__}"
func NewPathPattern(p string) (*PathPattern, error) {
	sp := p
	sp = strings.ReplaceAll(p, ".", "\\.")
	// use a token (without *'s to replace the double stars, so single stars can be replaced independently
	sp = strings.ReplaceAll(sp, "**", doubleStartToken)
	sp = strings.ReplaceAll(sp, "*", "[^,.]*")
	sp = strings.ReplaceAll(sp, doubleStartToken, ".*")
	sp = strings.Join([]string{"^", sp}, "")
	rx, err := regexp.Compile(sp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path into a searchable regexp. %v", err)
	}
	return &PathPattern{regex:rx}, nil
}
