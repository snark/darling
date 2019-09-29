package filter

import (
	"github.com/mmcdole/gofeed"
	"regexp"
	"strings"
)

type ItemFilter interface {
	Match(gofeed.Item) bool
}

// Some basic Boolean logic types
type True struct {
}

type And struct {
	Left  ItemFilter
	Right ItemFilter
}

type Or struct {
	Left  ItemFilter
	Right ItemFilter
}

type Not struct {
	Base ItemFilter
}

type Count struct {
	Limit int
	count int
}

type Regexp struct {
	regexps []*regexp.Regexp
}

func (filter *True) Match(i gofeed.Item) bool {
	return true
}

func (filter *And) Match(i gofeed.Item) bool {
	return filter.Left.Match(i) && filter.Right.Match(i)
}

func (filter *Or) Match(i gofeed.Item) bool {
	return filter.Left.Match(i) || filter.Right.Match(i)
}

func (filter *Not) Match(i gofeed.Item) bool {
	return !filter.Base.Match(i)
}

// Definitionally not idempotent!
func (filter *Count) Match(i gofeed.Item) bool {
	if filter.Limit == 0 {
		return true
	}
	filter.count = filter.count + 1
	return filter.count <= filter.Limit
}

// TODO: Does not currently handle item.Categories
func (filter *Regexp) Match(i gofeed.Item) bool {
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
func NewRegexp(words []string) ItemFilter {
	wildcard := false
	for _, word := range words {
		if word == "*" {
			wildcard = true
			break
		}
	}
	if wildcard {
		return &True{}
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
		return &Regexp{regexps: reSlice}
	}
}
