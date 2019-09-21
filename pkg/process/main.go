package process

import (
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/snark/darling/pkg/filter"
	"time"
)

func ProcessItems(parsedItems []*gofeed.Item, blacklistFilters []filter.ItemFilter, whitelistFilters []filter.ItemFilter) []*feeds.Item {
	outitems := []*feeds.Item{}
	for _, item := range parsedItems {
		blacklisted := false
		whitelisted := false
		for _, filter := range blacklistFilters {
			if filter.Match(*item) {
				blacklisted = true
			}
			break
		}
		for _, filter := range whitelistFilters {
			if filter.Match(*item) {
				whitelisted = true
			}
			break
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
			outitems = append(outitems, newitem)
		}
	}
	return outitems
}
