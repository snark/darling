package darling

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"log"
	"net/url"
	"sort"
	"sync"
	"time"
)

func FilterFeeds(blacklistWords []string, whitelistWords []string, feedUrls []string) {
	var wg sync.WaitGroup

	blacklistFilter := filter.NewRegexpFilter(blacklistWords)
	whitelistFilter := filter.NewRegexpFilter(whitelistWords)

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
				outfeed.Items = append(outfeed.Items, parseFeedWithFilters(url, blacklistFilter, whitelistFilter)...)
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

func parseFeedWithFilters(url string, blacklistFilter filter.ItemFilter, whitelistFilter filter.ItemFilter) []*feeds.Item {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseURL(url)
	items := []*feeds.Item{}
	if err != nil {
		fmt.Println("unable to parse %s--skipping with error %s", url, err)
	} else {
		for _, item := range parsed.Items {
			blacklisted := false
			whitelisted := false
			if blacklistFilter.Match(*item) {
				blacklisted = true
			}
			if whitelistFilter.Match(*item) {
				whitelisted = true
			}
			if !blacklisted || whitelisted {
				// TODO: Currently unhandled:
				// * Author
				// * Enclosures
				// * Categories
				// * Extensions
				newitem := &feeds.Item{
					//Author: item.Author,
					Content:     item.Content,
					Description: item.Description,
					Id:          item.GUID,
					Link:        &feeds.Link{Href: item.Link},
					Title:       item.Title,
				}
				if item.PublishedParsed != nil {
					newitem.Created = *item.PublishedParsed
				} else {
					newitem.Created = time.Now()
				}
				if item.UpdatedParsed != nil {
					newitem.Updated = *item.UpdatedParsed
				}
				items = append(items, newitem)
			}
		}
	}
	return items
}
