package output

import (
	"github.com/gorilla/feeds"
)

// TODO: Handle errors

func FeedToAtom(outfeed *feeds.Feed) (string, error) {
	return outfeed.ToAtom()
}

func FeedToRss(outfeed *feeds.Feed) (string, error) {
	return outfeed.ToRss()
}
