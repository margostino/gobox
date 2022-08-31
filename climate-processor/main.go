package main

import (
	"github.com/mmcdole/gofeed"
	"log"
)

func main() {
	var itemsSelected []*gofeed.Item
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("some-feed")
	for _, item := range feed.Items {
		// TODO: ask user
		itemsSelected = append(itemsSelected, item)
		log.Println(item.UpdatedParsed.UTC().String() + " - " + item.Title)
		log.Println(item.Link)
		log.Println(item.Content)
	}
}
