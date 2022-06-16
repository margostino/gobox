package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/gobox/common"
	"strconv"
	"strings"
	"time"
)

var allJobPosts = make([]*JobPost, 0)

type JobPost struct {
	Position    string
	Company     string
	Location    string
	Benefit     string
	PostDate    string
	CaptureDate time.Time
}

func startFrom(url string, value int) string {
	param := strconv.Itoa(value)
	return fmt.Sprintf("%s&start=%s", url, param)
}

func main() {
	var index = 0
	var factor = 50
	var isEnd = false
	baseUrl := "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"
	params := "?keywords=software engineer&location=stockholm&trk=public_jobs_jobs-search-bar_search-submit"
	partialUrl := fmt.Sprintf("%s%s", baseUrl, params)
	//url := "https://www.linkedin.com/jobs/search?keywords=software%20engineer&location=Stockholm%2C%20Stockholm%20County%2C%20Sweden&trk=public_jobs_jobs-search-bar_search-submit&position=1&pageNum=0"
	browser := rod.New().MustConnect()
	defer browser.Close()

	for ok := true; ok; ok = !isEnd {
		url := startFrom(partialUrl, index+factor)
		page := browser.MustPage(url).MustWaitLoad()
		entries, err := page.Elements("li")
		common.Check(err)

		if len(entries) == 0 {
			isEnd = true
		}

		for _, entry := range entries {
			card, err := entry.Element(".base-search-card__info")
			if err != nil {
				println(err.Error())
				break
			} else {
				text, err := card.Text()
				common.Check(err)
				parts := strings.SplitN(text, "\n", -1)

				if len(parts) > 0 {
					jobPost := &JobPost{
						Position:    common.GetOrDefault(0, parts),
						Company:     common.GetOrDefault(1, parts),
						Location:    common.GetOrDefault(2, parts),
						Benefit:     common.GetOrDefault(3, parts),
						PostDate:    common.GetOrDefault(4, parts),
						CaptureDate: time.Now(),
					}

					allJobPosts = append(allJobPosts, jobPost)
				}
				println(text)
			}
		}
	}

	println("end")
	//for ok := true; ok; ok = !isEnd {
	//	err := page.Keyboard.Press(input.End)
	//
	//	if err != nil {
	//		println(err)
	//	}
	//
	//	moreButtonElement := page.MustElement(".infinite-scroller__show-more-button")
	//
	//	if strings.Contains(moreButtonElement.Object.Description, "infinite-scroller__show-more-button--visible") {
	//		//println(moreButtonElement.Object.Description)
	//		err = page.Keyboard.Release(input.End)
	//		wait := page.MustWaitRequestIdle("")
	//		println(wait)
	//		println(moreButtonElement.Text())
	//		a := moreButtonElement.MustClick()
	//		println(a)
	//		//isEnd = true
	//	}
	//	//moreButtonElementVisible, err := page.Element(".infinite-scroller__show-more-button .infinite-scroller__show-more-button--visible")
	//	//println(err)
	//	//println(moreButtonElementVisible)
	//}
	//
	////elements, _ := page.Elements("*")
	//element := page.MustElement(".jobs-search__results-list")
	//elements, _ := element.Elements("li")
	//for _, value := range elements {
	//	card := value.MustElement(".base-search-card__info")
	//	println(" ")
	//	println(card.Text())
	//
	//	//text, _ := value.Text()
	//	//println(text)
	//	//if text == moreCommentTextSelector {
	//	//	btn := value.MustElement("button")
	//	//	btn.Click(proto.InputMouseButtonLeft)
	//	//	break
	//	//}
	//}
}
