package filter

import (
	"github.com/mmcdole/gofeed"
	"regexp"
	"strings"
)

type ItemFilter interface {
	Match(gofeed.Item) bool
}

type RegexpFilter struct {
	regexps []*regexp.Regexp
}

type TrueFilter struct {
}

func (filter TrueFilter) Match(i gofeed.Item) bool {
	return true
}

// TODO: Does not currently handle item.Categories
func (filter RegexpFilter) Match(i gofeed.Item) bool {
	found := false
	for _, re := range filter.regexps {
		if re.MatchString(i.Content) {
			found = true
		} else if re.MatchString(i.Title) {
			found = true
		} else if re.MatchString(i.Description) {
			found = true
		}
	}
	return found
}

// TODO: Allow case-sensitive matching?
// TODO: Log error
func NewRegexpFilter(words []string) ItemFilter {
	wildcard := false
	for _, word := range words {
		if word == "*" {
			wildcard = true
			break
		}
	}
	if wildcard {
		return TrueFilter{}
	} else {
		reSlice := []*regexp.Regexp{}
		for _, word := range words {
			word = strings.TrimSpace(word)
			// Silently discard empty/whitespace strings
			if word != "" {
				var re, err = regexp.Compile(`(?i)\b` + word + `\b`)
				if err == nil {
					reSlice = append(reSlice, re)
				}
			}
		}
		return RegexpFilter{regexps: reSlice}
	}
}
