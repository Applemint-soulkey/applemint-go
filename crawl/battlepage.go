package crawl

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bobesa/go-domain-util/domainutil"
)

const BASE_URL_BP = "https://v12.battlepage.com"
const PAGE_SIZE = 5

func CrawlBP() int {
	targetList := []string{}

	// Get Target List
	log.Print("Get Target List")
	for i := 0; i < PAGE_SIZE; i++ {
		targetList = append(targetList, BASE_URL_BP+"/??=Board.Humor.Table&page="+strconv.Itoa(i+1))
		targetList = append(targetList, BASE_URL_BP+"/??=Board.ETC.Table&page="+strconv.Itoa(i+1))
	}

	// Get Items
	log.Print("Get Items")
	items := []Item{}
	for _, targetURL := range targetList {
		doc := getPageDocument(targetURL)
		items = append(items, getItemsFromBP(doc)...)
	}

	// Insert Items
	log.Print("Insert Items")
	insertedCount := InsertItems(items, "new")

	return insertedCount
}

func getItemsFromBP(doc *goquery.Document) []Item {
	var items []Item
	doc.Find(".ListTable div").Each(func(i int, s *goquery.Selection) {
		items = append(items, getItemFromBP(s))
	})
	return items
}

func getItemFromBP(doc *goquery.Selection) Item {
	regexpPage := regexp.MustCompile(`&page=[0-9]`)

	item := Item{}
	item.TextContent, _ = doc.Find(".bp_subject").Attr("title")
	itemLink, _ := doc.Find("a").Attr("href")
	fmt.Println(itemLink)
	itemLink = regexpPage.ReplaceAllLiteralString(itemLink, "")
	fmt.Println(itemLink)
	item.Url = BASE_URL_BP + itemLink
	item.Domain = domainutil.Domain(item.Url)
	item.Tags = []string{}
	item.Path = ""
	item.Timestamp = time.Now()
	item.Source = "battlepage"
	return item
}
