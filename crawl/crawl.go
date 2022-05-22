package crawl

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	TextContent string
	Url   string
	Timestamp time.Time
	Source string
}

func Crawl(url string) {
	fmt.Println(("I'm crawl"))
	CrawlBP()
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkResponseCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Response code: ", res.StatusCode)
	}
}

func getPageDocument(targetURL string) *goquery.Document {
	res, err := http.Get(targetURL)
	checkError(err)
	checkResponseCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	return doc
}
