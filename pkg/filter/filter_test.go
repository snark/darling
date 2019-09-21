package filter_test

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"os"
	"reflect"
	"testing"
)

// NEEDS TESTING:
// Match against Title
// Match against Content
// Match against Description

func parsedFeedFromFile(fpath string) *gofeed.Feed {
	file, _ := os.Open(fpath)
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	return feed
}

func TestRegexpFilterBasic(t *testing.T) {
	// Create a filter and match it against some real data.
	// We exercise case-insensitivity and word boundaries
	// in our regexps. This test is matching titles only.
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	reFilter := filter.NewRegexpFilter([]string{
		"package",
		"potential",
		"tcp",
		"xyzzy", // Matches nothing
	})
	var tests = []struct {
		i    int
		want bool
	}{
		{0, false},
		{1, false},
		{2, false},
		{3, true}, // TCP: Note case-insensitivity
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, false},
		{9, false},
		{10, false},
		{11, false},
		{12, true},
		{13, false}, // "potentially" -- word boundaries
		{14, false},
		{15, true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Lobste.rs Item #%d", tt.i)
		t.Run(testname, func(t *testing.T) {
			ans := reFilter.Match(*feed.Items[tt.i])
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestRegexpFilterMultiWord(t *testing.T) {
	wcFilter := filter.NewRegexpFilter([]string{
		"this weekend",
	})
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	ans := wcFilter.Match(*feed.Items[1])
	if ans != true {
		t.Errorf("got false from multiword token ('this weekend') against full match")
	}
	ans = wcFilter.Match(*feed.Items[6])
	if ans != false {
		t.Errorf("got true from multiword token ('this weekend') against partial match ('this')")
	}
}

func TestRegexpFilterWhiteSpace(t *testing.T) {
	wcFilter := filter.NewRegexpFilter([]string{
		"",
		"  ",
		"\t",
	})
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	for i, _ := range feed.Items {
		ans := wcFilter.Match(*feed.Items[i])
		if ans != false {
			t.Errorf("got true from whitespace match on Lobste.rs item %d", i)
		}
	}
}

func TestRegexpFilterWildcard(t *testing.T) {
	wcFilter := filter.NewRegexpFilter([]string{
		"foo",
		"*",
		"bar",
	})
	filterType := reflect.TypeOf(wcFilter)
	if filterType.String() != "filter.TrueFilter" {
		t.Errorf("Wildcarded regexp filter yielded %s", filterType.String())
	}
}

func TestTrueFilter(t *testing.T) {
	// Matches anything
	trueFilter := filter.TrueFilter{}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := trueFilter.Match(*feed.Items[i])
			if ans != true {
				t.Errorf("got %t, want true", ans)
			}
		})
	}
}
