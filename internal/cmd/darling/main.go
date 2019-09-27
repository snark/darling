package darling

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/snark/darling/pkg/feed"
	"github.com/snark/darling/pkg/filter"
	"log"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

func FilterFeeds(blacklistWords []string, whitelistWords []string, feedUrls []string) {
	var wg sync.WaitGroup

	blacklist := filter.NewRegexpFilter(blacklistWords)
	whitelist := filter.NewRegexpFilter(whitelistWords)
	wordMatch := filter.AndFilter{blacklist, filter.NotFilter{whitelist}}

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
			go func() {
				defer wg.Done()
				f, err := feed.Fetch(url)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Unable to fetch %s: %s", url, err)
				} else {
					outfeed.Items = append(outfeed.Items, feed.ProcessItems(f.Items, []filter.ItemFilter{wordMatch})...)
				}
			}()
		}
	}
	wg.Wait()
	// Reverse chronological order
	sort.SliceStable(outfeed.Items, func(a, b int) bool {
		return outfeed.Items[a].Created.After(outfeed.Items[b].Created)
	})

	atom, err := outfeed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Handle error
	_ = atom
	fmt.Println(atom)
}

func validateUrl(toTest string) bool {
	uri, err := url.Parse(toTest)
	return err == nil && (uri.Scheme == "http" || uri.Scheme == "https")
}
