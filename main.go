package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var blacklistWords arrayFlags
	var whitelistWords arrayFlags
	flag.Var(&blacklistWords, "b", "blacklist term")
	flag.Var(&whitelistWords, "w", "whitelist term")
	flag.Parse()
	tail := flag.Args()

	var wg sync.WaitGroup

	blacklistRegexps := wordsToRegexps(blacklistWords)
	whitelistRegexps := wordsToRegexps(whitelistWords)

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
	for _, url := range tail {
		// TODO: Warning messages on bad URLs
		if validateUrl(url) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				outfeed.Items = append(outfeed.Items, parseFeedWithFilters(url, blacklistRegexps, whitelistRegexps)...)
			}()
		}
	}
	wg.Wait()
	atom, err := outfeed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	_ = atom
	fmt.Println(atom)
}

func wordsToRegexps(words []string) []*regexp.Regexp {
	outslice := []*regexp.Regexp{}
	for _, word := range words {
		var re, err = regexp.Compile(`(?i)\b` + word + `\b`)
		if err == nil {
			outslice = append(outslice, re)
		}
	}
	return outslice
}

func validateUrl(toTest string) bool {
	uri, err := url.Parse(toTest)
	return err == nil && (uri.Scheme == "http" || uri.Scheme == "https")
}

func itemMatchesRegexp(item *gofeed.Item, re *regexp.Regexp) bool {
	// Does not currently handle item.Categories
	found := false
	if re.MatchString(item.Content) {
		found = true
	} else if re.MatchString(item.Title) {
		found = true
	} else if re.MatchString(item.Description) {
		found = true
	}
	return found
}

func parseFeedWithFilters(url string, blacklistRegexps []*regexp.Regexp, whitelistRegexps []*regexp.Regexp) []*feeds.Item {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseURL(url)
	items := []*feeds.Item{}
	if err != nil {
		fmt.Println("unable to parse", url, "-- skipping with error", err)
	} else {
		for _, item := range parsed.Items {
			blacklisted := false
			whitelisted := false
			for _, blacklistRe := range blacklistRegexps {
				if itemMatchesRegexp(item, blacklistRe) {
					blacklisted = true
					break
				}
			}
			for _, whitelistWord := range whitelistRegexps {
				if itemMatchesRegexp(item, whitelistWord) {
					whitelisted = true
					break
				}
			}
			if !blacklisted || whitelisted {
				// Currently unhandled:
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
				created, err := time.Parse(time.RFC3339, item.Published)
				if err == nil {
					newitem.Created = created
				}
				updated, err := time.Parse(time.RFC3339, item.Updated)
				if err == nil {
					newitem.Updated = updated
				}
				items = append(items, newitem)
			}
		}
	}
	return items
}
