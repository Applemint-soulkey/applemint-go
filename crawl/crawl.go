package crawl

import (
	"crypto/tls"
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
	CrawlDD()
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{MaxVersion: tls.VersionTLS12},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", targetURL, nil)
	checkError(err)
	res, err := client.Do(req)
	checkError(err)
	checkResponseCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	return doc
}
