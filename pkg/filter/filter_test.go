package filter_test

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"os"
	"reflect"
	"testing"
	"time"
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

func TestNewSinceDuration(t *testing.T) {
	origin := "2019-10-12T16:25:00Z"
	fakeNow, _ := time.Parse(time.RFC3339, origin)
	var tests = []struct {
		duration string
		expected string
	}{
		{"30S", "2019-10-12T16:24:30Z"},
		{"15M", "2019-10-12T16:10:00Z"},
		{"12H", "2019-10-12T04:25:00Z"},
		{"1d", "2019-10-11T16:25:00Z"},
		{"3m", "2019-07-12T16:25:00Z"},
		{"1Y", "2018-10-12T16:25:00Z"},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Generating Since filter for %s", tt.duration)
		t.Run(testname, func(t *testing.T) {
			f, _ := filter.NewSince(&tt.duration, fakeNow)
			expectedTime, _ := time.Parse(time.RFC3339, tt.expected)
			beforeTime := expectedTime.Add(-1 * time.Second)
			afterTime := expectedTime.Add(1 * time.Second)
			beforeItem := gofeed.Item{
				Title:           "Test",
				PublishedParsed: &beforeTime,
			}
			afterItem := gofeed.Item{
				Title:           "Test",
				PublishedParsed: &afterTime,
			}
			if f.Match(beforeItem) {
				t.Errorf("Since filter for %s matched one second before %s with now at %s", tt.duration, tt.expected, origin)
			}
			if !f.Match(afterItem) {
				t.Errorf("Since filter for %s did not match one second after %s with now at %s", tt.duration, tt.expected, origin)
			}
		})
	}
}

func TestNewSinceAgainstNoTimeItem(t *testing.T) {
	i := gofeed.Item{
		Title: "Test",
	}
	var tests = []struct {
		when     string
		expected bool
	}{
		{"2019-10-30", false},
		{"7d", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Generating Since filter for %s", tt.when)
		t.Run(testname, func(t *testing.T) {
			f, _ := filter.NewSince(&tt.when, time.Now())
			ans := f.Match(i)
			if ans != tt.expected {
				t.Errorf("Since filter for %s expected %t, got %t for timestampless item", tt.when, tt.expected, ans)
			}
		})
	}
}

func TestNewSinceTimestamp(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2019-10-12T16:25:00Z")
	i := gofeed.Item{
		Title:           "Test",
		PublishedParsed: &timestamp,
	}
	i2 := gofeed.Item{
		Title:         "Test",
		UpdatedParsed: &timestamp,
	}
	var tests = []struct {
		when     string
		expected bool
	}{
		{"2019-10-12", true},
		{"2019-10-13", false},
		{"2019-10-12T04:24:00Z", true},
		{"2019-10-12T04:26:00Z", true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Generating Since filter for %s", tt.when)
		t.Run(testname, func(t *testing.T) {
			f, _ := filter.NewSince(&tt.when, time.Now())
			ans := f.Match(i)
			if ans != tt.expected {
				t.Errorf("Since filter for %s expected %t, got %t for 2019-10-12T16:25:00Z", tt.when, tt.expected, ans)
			}
			ans2 := f.Match(i2)
			if ans2 != tt.expected {
				t.Errorf("Since filter for %s expected %t, got %t for 2019-10-12T16:25:00Z", tt.when, tt.expected, ans2)
			}
		})
	}
}

func TestNewSinceUnparseable(t *testing.T) {
	nogood := []string{"", "100", "32x", "2019-10-12T", "2019", "2019-13-01", "2019-10-32", "2019-10-01-01"}
	for _, which := range nogood {
		_, err := filter.NewSince(&which, time.Now())
		if err == nil {
			t.Errorf("did not throw error on unparseable since string %s", which)
		}
	}
}
