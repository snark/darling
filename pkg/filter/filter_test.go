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

func TestTrue(t *testing.T) {
	// Matches anything
	trueF := filter.True{}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := trueF.Match(*feed.Items[i])
			if ans != true {
				t.Errorf("got %t, want true", ans)
			}
		})
	}
}

func TestNot(t *testing.T) {
	// Matches anything
	trueF := filter.True{}
	notF := filter.Not{&trueF}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := notF.Match(*feed.Items[i])
			if ans != false {
				t.Errorf("got %t, want false", ans)
			}
		})
	}
}

func TestOr(t *testing.T) {
	// Matches anything
	trueF := filter.True{}
	notF := filter.Not{&trueF}
	orF1 := filter.Or{&trueF, &notF}
	orF2 := filter.Or{&notF, &trueF}
	orF3 := filter.Or{&notF, &notF}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := orF1.Match(*feed.Items[i])
			if ans != true {
				t.Errorf("got %t, want true", ans)
			}
			ans2 := orF2.Match(*feed.Items[i])
			if ans2 != true {
				t.Errorf("got %t, want true", ans)
			}
			ans3 := orF3.Match(*feed.Items[i])
			if ans3 != false {
				t.Errorf("got %t, want false", ans)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	// Matches anything
	trueF := filter.True{}
	notF := filter.Not{&trueF}
	andF1 := filter.And{&trueF, &notF}
	andF2 := filter.And{&notF, &trueF}
	andF3 := filter.And{&trueF, &trueF}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans := andF1.Match(*feed.Items[i])
			if ans != false {
				t.Errorf("got %t, want false ", ans)
			}
			ans2 := andF2.Match(*feed.Items[i])
			if ans2 != false {
				t.Errorf("got %t, want false ", ans)
			}
			ans3 := andF3.Match(*feed.Items[i])
			if ans3 != true {
				t.Errorf("got %t, want true", ans)
			}
		})
	}
}

func TestCount(t *testing.T) {
	// Matches anything
	count1 := filter.Count{Limit: 4}
	count2 := filter.Count{Limit: 0}
	file, _ := os.Open("../../testdata/lobste.rs.rss")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	for i, _ := range feed.Items {
		testname := fmt.Sprintf("Lobste.rs Item #%d", i)
		t.Run(testname, func(t *testing.T) {
			ans1 := count1.Match(*feed.Items[i])
			expected := i < 4
			if ans1 != expected {
				t.Errorf("got %t, want %t", ans1, expected)
			}
			ans2 := count2.Match(*feed.Items[i])
			if ans2 != true {
				t.Errorf("got %t, want true", ans2)
			}
		})
	}
}

func TestRegexpBasic(t *testing.T) {
	// Create a filter and match it against some real data.
	// We exercise case-insensitivity and word boundaries
	// in our regexps. This test is matching titles only.
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	re := filter.NewRegexp([]string{
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
			ans := re.Match(*feed.Items[tt.i])
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestRegexpTitle(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wc := filter.NewRegexp([]string{"gerry"})
	ans := wc.Match(*feed.Items[0])
	if ans != true {
		t.Errorf("got false from title match")
	}
}

func TestRegexpDescription(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wc := filter.NewRegexp([]string{"delightfully"})
	ans := wc.Match(*feed.Items[13])
	if ans != true {
		t.Errorf("got false from description match")
	}
}

func TestRegexpContent(t *testing.T) {
	feed := parsedFeedFromFile("../../testdata/waxy.org.rss")
	wc := filter.NewRegexp([]string{"amazing creators"})
	ans := wc.Match(*feed.Items[13])
	if ans != true {
		t.Errorf("got false from content match")
	}
}

func TestRegexpMultiWord(t *testing.T) {
	wc := filter.NewRegexp([]string{
		"this weekend",
	})
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	ans := wc.Match(*feed.Items[1])
	if ans != true {
		t.Errorf("got false from multiword token ('this weekend') against full match")
	}
	ans = wc.Match(*feed.Items[6])
	if ans != false {
		t.Errorf("got true from multiword token ('this weekend') against partial match ('this')")
	}
}

func TestRegexpWhiteSpace(t *testing.T) {
	wc := filter.NewRegexp([]string{
		"",
		"  ",
		"\t",
	})
	feed := parsedFeedFromFile("../../testdata/lobste.rs.rss")
	for i, _ := range feed.Items {
		ans := wc.Match(*feed.Items[i])
		if ans != false {
			t.Errorf("got true from whitespace match on Lobste.rs item %d", i)
		}
	}
}

func TestRegexpWildcard(t *testing.T) {
	wc := filter.NewRegexp([]string{
		"foo",
		"*",
		"bar",
	})
	filterType := reflect.TypeOf(wc)
	if filterType.String() != "*filter.True" {
		t.Errorf("Wildcarded regexp filter yielded %s", filterType.String())
	}
}
