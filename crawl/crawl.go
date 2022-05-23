package crawl

import (
	"crypto/tls"
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

func GetCrawlBPResult() []Item {
	return CrawlBP()
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
		TLSClientConfig: &tls.Config{MaxVersion: tls.VersionTLS12, InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	checkError(err)
	res, err := client.Do(req)
	checkError(err)
	checkResponseCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	return doc
}
