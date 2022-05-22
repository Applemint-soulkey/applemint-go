package crawl

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func CrawlDD() []Item {
	targetList := []string{}

	// Get Target List
	log.Print("Get Target List")
	for i := 0; i < PAGE_SIZE; i++ {
		targetList = append(targetList, BASE_URL_DD+"dogdrip?page="+strconv.Itoa(i+1))
	}
	fmt.Println(targetList)

	// Get Items
	log.Print("Get Items")
	items := []Item{}
	for _, targetURL := range targetList {
		doc := getPageDocument(targetURL)
		items = append(items, getItemsFromDD(doc)...)
	}
	for _, item := range items {
		fmt.Println(item)
	}

	return items
}

func getItemsFromDD(doc *goquery.Document) []Item {
	var items []Item
	doc.Find("table.ed.table.table-divider tbody tr").Each(func(i int, s *goquery.Selection) {
		items = append(items, getItemFromDD(s))
	})
	return items
}

func getItemFromDD(doc *goquery.Selection) Item {
	item := Item{}
	temp := doc.Find(".title a")
	item.TextContent = strings.Trim(temp.Find(".ed.title-link").Text(), " ") 
	itemLink, _ := temp.Attr("href")
	item.Url = itemLink
	item.Timestamp = time.Now()
	item.Source = "dogdrip"
	return item
}
