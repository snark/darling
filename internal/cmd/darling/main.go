package darling

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/snark/darling/pkg/feed"
	"github.com/snark/darling/pkg/filter"
	"github.com/snark/darling/pkg/output"
	"log"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

func FilterFeeds(blacklistWords []string, whitelistWords []string, limit *int, feedUrls []string) {
	var wg sync.WaitGroup

	blacklist := filter.NewRegexp(blacklistWords)
	whitelist := filter.NewRegexp(whitelistWords)
	wordMatch := filter.Or{&filter.Not{blacklist}, whitelist}

	now := time.Now()
	outfeed := &feeds.Feed{
		Title:       "Darling",
		Description: "Your darlings, killfiled",
		Created:     now,
		// Link and Author are absolutetly erquired by feeds
		Link:   &feeds.Link{Href: ""},
		Author: &feeds.Author{Name: "You"},
	}
	outfeed.Items = []*feeds.Item{}
	for _, url := range feedUrls {
		// TODO: Warning messages on bad URLs
		if validateUrl(url) {
			wg.Add(1)
			go func(u string, filters []filter.ItemFilter) {
				if *limit > 0 {
					filters = append(filters, &filter.Count{Limit: *limit})
				}
				defer wg.Done()
				f, err := feed.Fetch(u)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to fetch %s: %s", u, err)
				} else {
					outfeed.Items = append(outfeed.Items, feed.ProcessItems(f.Items, filters)...)
				}
			}(url, []filter.ItemFilter{&wordMatch})
		}
	}
	wg.Wait()
	// Reverse chronological order
	sort.SliceStable(outfeed.Items, func(a, b int) bool {
		return outfeed.Items[a].Created.After(outfeed.Items[b].Created)
	})

	result, err := output.FeedToRss(outfeed)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

func validateUrl(toTest string) bool {
	uri, err := url.Parse(toTest)
	return err == nil && (uri.Scheme == "http" || uri.Scheme == "https")
}
