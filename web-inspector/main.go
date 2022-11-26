package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/gobox/common"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Attributes struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attributes
}

func main() {
	url := "https://earth.org/sea-level-rise-nyc/" //"https://earth.org/"

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	data := parse(string(text))
	data = nil
	fmt.Println(data)

}

func isValidTokenType(tokenType html.TokenType) bool {
	return tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken || tokenType == html.TextToken
}

func isValidData(token html.Token) bool {
	return !strings.Contains(token.Data, "\n")
}

func parse(text string) []*Token {

	var tokens = make([]*Token, 0)
	tkn := html.NewTokenizer(strings.NewReader(text))

	for {
		tokenType := tkn.Next()
		currentToken := tkn.Token()

		if isValidTokenType(tokenType) && isValidData(currentToken) {
			attrs := make([]*Attributes, 0)
			for _, attr := range currentToken.Attr {
				att := &Attributes{
					Key:   attr.Key,
					Value: attr.Val,
				}
				attrs = append(attrs, att)
			}
			token := &Token{
				Type:       tokenType,
				Data:       currentToken.Data,
				Attributes: attrs,
			}
			tokens = append(tokens, token)
		}

		if tokenType == html.ErrorToken {
			return tokens
		}

	}
}

func scrapper() {
	//url := "https://earth.org/impact-and-reach/"
	url := "https://en.wikipedia.org/wiki/The_Lord_of_the_Rings"
	//url := "https://en.wikipedia.org/wiki/Dune_(novel)"
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := browser.MustPage(url).MustWaitLoad()

	//heading, err := page.Element("#firstHeading")
	//common.Check(err)
	//title, err := heading.Text()
	//common.Check(err)
	//println(title)

	//bodyElements, err := page.Elements("#bodyContent")
	//common.Check(err)
	//for _, element := range bodyElements {
	//	text, _ := element.Text()
	//	println(text)
	//}

	bodyElements, err := page.Elements(".infobox.vcard")
	common.Check(err)
	for _, element := range bodyElements {
		text, _ := element.Text()
		println(text)
	}
}
