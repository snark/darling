package filter_test

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"os"
	"reflect"
	"testing"
)

func parsedFeedFromFile(fpath string) *gofeed.Feed {
	file, _ := os.Open(fpath)
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	return feed
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

func TestNotFilter(t *testing.T) {
	// Matches anything
	trueFilter := filter.TrueFilter{}
	notFilter := filter.NotFilter{trueFilter}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := notFilter.Match(*feed.Items[i])
			if ans != false {
				t.Errorf("got %t, want false", ans)
			}
		})
	}
}

func TestOrFilter(t *testing.T) {
	// Matches anything
	trueFilter := filter.TrueFilter{}
	notFilter := filter.NotFilter{trueFilter}
	orFilter1 := filter.OrFilter{trueFilter, notFilter}
	orFilter2 := filter.OrFilter{notFilter, trueFilter}
	orFilter3 := filter.OrFilter{notFilter, notFilter}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := orFilter1.Match(*feed.Items[i])
			if ans != true {
				t.Errorf("got %t, want true", ans)
			}
			ans2 := orFilter2.Match(*feed.Items[i])
			if ans2 != true {
				t.Errorf("got %t, want true", ans)
			}
			ans3 := orFilter3.Match(*feed.Items[i])
			if ans3 != false {
				t.Errorf("got %t, want false", ans)
			}
		})
	}
}

func TestAndFilter(t *testing.T) {
	// Matches anything
	trueFilter := filter.TrueFilter{}
	notFilter := filter.NotFilter{trueFilter}
	andFilter1 := filter.AndFilter{trueFilter, notFilter}
	andFilter2 := filter.AndFilter{notFilter, trueFilter}
	andFilter3 := filter.AndFilter{trueFilter, trueFilter}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := andFilter1.Match(*feed.Items[i])
			if ans != false {
				t.Errorf("got %t, want false ", ans)
			}
			ans2 := andFilter2.Match(*feed.Items[i])
			if ans2 != false {
				t.Errorf("got %t, want false ", ans)
			}
			ans3 := andFilter3.Match(*feed.Items[i])
			if ans3 != true {
				t.Errorf("got %t, want true", ans)
			}
		})
	}
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
		{12, true},  // package
		{13, false}, // "potentially" -- word boundaries
		{14, false},
		{15, true}, // package
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

func TestRegexpFilterTitle(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wcFilter := filter.NewRegexpFilter([]string{"gerry"})
	ans := wcFilter.Match(*feed.Items[0])
	if ans != true {
		t.Errorf("got false from title match")
	}
}

func TestRegexpFilterDescription(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wcFilter := filter.NewRegexpFilter([]string{"delightfully"})
	ans := wcFilter.Match(*feed.Items[13])
	if ans != true {
		t.Errorf("got false from description match")
	}
}

func TestRegexpFilterContent(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wcFilter := filter.NewRegexpFilter([]string{"amazing creators"})
	ans := wcFilter.Match(*feed.Items[13])
	if ans != true {
		t.Errorf("got false from content match")
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
