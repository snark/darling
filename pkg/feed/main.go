package feed

import (
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"time"
)

func Fetch(url string) (*gofeed.Feed, error) {
	// TODO: Caching support
	// TODO: Load from file
	fp := gofeed.NewParser()
	parsed, err := fp.ParseURL(url)
	return parsed, err
}

func ProcessItems(parsedItems []*gofeed.Item, filters []filter.ItemFilter) []*feeds.Item {
	outitems := []*feeds.Item{}
	for _, item := range parsedItems {
		matched := false
		for _, filter := range filters {
			if filter.Match(*item) {
				matched = true
				break
			}
		}
		if !matched {
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
