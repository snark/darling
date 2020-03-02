package darling

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/snark/darling/pkg/feed"
	"github.com/snark/darling/pkg/filter"
	"github.com/snark/darling/pkg/output"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

func FilterFeeds(blacklistWords []string, whitelistWords []string, since *string, limit *int, outputType *string, feedTokens []string) {
	var wg sync.WaitGroup

	blacklist := filter.NewRegexp(blacklistWords)
	whitelist := filter.NewRegexp(whitelistWords)
	wordMatch := filter.Or{&filter.Not{blacklist}, whitelist}
	var sinceMatch filter.ItemFilter
	var err error
	if *since != "" {
		sinceMatch, err = filter.NewSince(since, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}

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
	filterList := []filter.ItemFilter{&wordMatch}
	if sinceMatch != nil {
		filterList = append(filterList, sinceMatch)
	}
	for _, token := range feedTokens {
		if validateUrl(token) {
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
			}(token, filterList)
		} else if validatePath(token) {
			wg.Add(1)
			go func(fp string, filters []filter.ItemFilter) {
				if *limit > 0 {
					filters = append(filters, &filter.Count{Limit: *limit})
				}
				defer wg.Done()
				file, err := os.Open(fp)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to read %s: %s", fp, err)
					return
				} else {
					defer file.Close()
					buf, err := ioutil.ReadAll(file)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Unable to read %s: %s", fp, err)
						return
					}
					f, err := feed.ParseFromString(string(buf))
					if err != nil {
						fmt.Fprintf(os.Stderr, "Unable to read %s: %s", fp, err)
					} else {
						outfeed.Items = append(outfeed.Items, feed.ProcessItems(f.Items, filters)...)
					}
				}
			}(token, filterList)
		} else {
			fmt.Fprintf(os.Stderr, "Unable to process %s as a feed location", token)
		}
		// TODO: Warning messages on bad tokens
	}
	wg.Wait()
	// Reverse chronological order
	sort.SliceStable(outfeed.Items, func(a, b int) bool {
		return outfeed.Items[a].Created.After(outfeed.Items[b].Created)
	})

	var result string
	if *outputType == "atom" {
		result, err = output.FeedToAtom(outfeed)
	} else {
		result, err = output.FeedToRss(outfeed)
	}
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

func validatePath(toTest string) bool {
	var res bool
	if _, err := os.Stat(toTest); err == nil {
		res = true
	} else if os.IsNotExist(err) {
		res = false
	} else {
		// An error other than non-existence; let's return true and report
		// the error when we try to open it.
		res = true
	}
	return res
}
