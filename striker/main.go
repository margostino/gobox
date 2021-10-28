package main

import (
	"github.com/go-rod/rod"
	"github.com/margostino/gobox/common"
	"log"
	"strings"
)

func main() {
	var config = common.GetStrikerConfig("./configuration/configuration.yml")
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()
	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()
	// Create a new page
	page := browser.MustPage(config.Url)

	// Process
	e := page.MustElement(".Summary_toggleAllComments__1M7FV").MustElement("*").MustClick()
	println(e.MustText()) //--> "Less commentary"
	rawEvents := page.MustElement(".Summary_blockWrapper__1P4fu").MustText()
	events := strings.Split(rawEvents, "\n")
	for _, event := range events {
		log.Println(event)
	}

	//exists := page.MustElements("#match-detail__summary__342__time").Empty())
	//time := page.MustElement("#match-detail__summary__339__time").MustText()
	//comment := page.MustElement("#match-detail__summary__339__comment").MustText()
}
