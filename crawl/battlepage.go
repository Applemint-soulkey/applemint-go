package crawl

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var PAGE_SIZE = 5

func CrawlBP() []Item {
	targetList := []string{}
	
	// Get Target List
	log.Print("Get Target List")
	for i := 0; i < PAGE_SIZE; i++ {
		targetList = append(targetList, BASE_URL_BP + "/??=Board.Humor.Table&page=" + strconv.Itoa(i+1))
		targetList = append(targetList, BASE_URL_BP + "/??=Board.ETC.Table&page=" + strconv.Itoa(i+1))
	}
	fmt.Println(targetList)

	// Get Items
	log.Print("Get Items")
	items := []Item{}
	for _, targetURL := range targetList {
		doc := getPageDocument(targetURL)
		items = append(items,  getItemsFromBP(doc)...)
	}

	return items
}

func getItemsFromBP(doc *goquery.Document) []Item {
	var items []Item
	doc.Find(".ListTable div").Each(func(i int, s *goquery.Selection) {
		items = append(items, getItemFromBP(s))
	})
	return items
}

func getItemFromBP(doc *goquery.Selection) Item {
	item := Item{}
	item.TextContent, _ = doc.Find(".bp_subject").Attr("title")
	itemLink, _ := doc.Find("a").Attr("href")
	item.Url = BASE_URL_BP+itemLink
	item.Timestamp = time.Now()
	item.Source = "battlepage"
	return item
}