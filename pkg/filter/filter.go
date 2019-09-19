package filter

import (
	"github.com/mmcdole/gofeed"
	"regexp"
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

func (filter RegexpFilter) Match(i gofeed.Item) bool {
	// TODO: Does not currently handle item.Categories
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
			// TODO: Allow case-sensitive matching?
			var re, err = regexp.Compile(`(?i)\b` + word + `\b`)
			// TODO: Log error
			if err == nil {
				reSlice = append(reSlice, re)
			}
		}
		return RegexpFilter{regexps: reSlice}
	}
}
