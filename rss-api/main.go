package main

import (
	"github.com/margostino/gobox/common"
	"github.com/mmcdole/gofeed"
	"log"
	"net/url"
	"strings"
)

func main() {
	feedUrl := "https://www.earthdata.nasa.gov/topics/rss/Human%20Dimensions"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedUrl)

	if feed != nil {
		for _, entry := range feed.Items {
			var link, source string
			rawLink, err := url.Parse(entry.Link)

			if common.IsError(err, "when parsing feed link") {
				link = entry.Link
			} else {
				link = rawLink.Query().Get("url")
				sourceUrl, err := url.Parse(link)
				if !common.IsError(err, "when parsing source link") {
					source = strings.ReplaceAll(sourceUrl.Hostname(), "www.", "")
				}
			}
			log.Printf("Link: %s\nSource: %s\nItem: %s\n", link, source, entry)
		}
	} else {
		log.Printf("There are no feeds")
	}

}
