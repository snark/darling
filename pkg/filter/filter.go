package filter

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"regexp"
	"strconv"
	"strings"
	"time"
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

type Since struct {
	When time.Time
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

// TODO: Allow additional arguments allowing us to also check
// updated times.
func (filter *Since) Match(i gofeed.Item) bool {
	if i.PublishedParsed != nil {
		return i.PublishedParsed.After(filter.When)
	}
	if i.UpdatedParsed != nil {
		return i.UpdatedParsed.After(filter.When)
	}
	return false
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

// We want to accept a couple different options here:
// * RFC3339 (2006-01-02T15:04:05Z07:00)
// * RFC3339, date only (2006-01-02)
// * nx, where n is an integer and x is an numeric indicator
//   from standard Unix date formatting (one of Y, m, d, H, M, S)
//   Note that time.Duration cannot be used because it caps at hours
var whenFormat = regexp.MustCompile(`^(?P<Num>\d+)(?P<Dur>[YmdHMS]{1})$`)

func NewSince(when *string, now time.Time) (ItemFilter, error) {
	whenMatch := whenFormat.FindStringSubmatch(*when)
	dateLayout := "2006-01-02"
	if len(whenMatch) == 3 {
		num, err := strconv.Atoi(whenMatch[1])
		if err != nil {
			return &True{}, fmt.Errorf("Unable to parse %s", *when)
		}
		var d time.Time
		switch whenMatch[2] {
		case "S":
			d = now.Add(time.Second * -1 * time.Duration(num))
		case "M":
			d = now.Add(time.Minute * -1 * time.Duration(num))
		case "H":
			d = now.Add(time.Hour * -1 * time.Duration(num))
		case "d":
			d = now.AddDate(0, 0, -1*num)
		case "m":
			d = now.AddDate(0, -1*num, 0)
		case "Y":
			d = now.AddDate(-1*num, 0, 0)
		default:
			return &True{}, fmt.Errorf("Unable to parse %s", *when)
		}
		return &Since{When: d}, nil
	} else {
		d, err := time.Parse(time.RFC3339, *when)
		if err != nil {
			d, err = time.Parse(dateLayout, *when)
		}
		if err != nil {
			return &True{}, fmt.Errorf("Unable to parse %s", *when)
		}
		return &Since{When: d}, nil
	}
}
