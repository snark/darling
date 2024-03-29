package feed

import (
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"time"
)

func Fetch(url string) (*gofeed.Feed, error) {
	// TODO: Caching support
	fp := gofeed.NewParser()
	parsed, err := fp.ParseURL(url)
	return parsed, err
}

func ParseFromString(s string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseString(s)
	return parsed, err
}

func ProcessItems(parsedItems []*gofeed.Item, filters []filter.ItemFilter) []*feeds.Item {
	outitems := []*feeds.Item{}
	for _, item := range parsedItems {
		missed := false
		for i := range filters {
			if !filters[i].Match(*item) {
				missed = true
				break
			}
		}
		if !missed {
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
			outitems = append(outitems, newitem)
		}
	}
	return outitems
}
