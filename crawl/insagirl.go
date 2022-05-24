package crawl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/linkify"
)

type respType struct {
	V []string `json:"v"`
}

func CrawlISG() int {
	targetList := []string{
		"http://insagirl-hrm.appspot.com/json2/1/1/2/",
		"http://insagirl-hrm.appspot.com/json2/2/1/2/",
	}

	items := []Item{}
	// Get Items
	log.Print("Get Items")
	for _, targetURL := range targetList {
		items = append(items, getIsgData(targetURL)...)
	}

	log.Print("Insert Items")
	insertedCount := InsertItems(items)

	return insertedCount
}

func getIsgData(url string) []Item {
	resp, err := http.Get(url)
	checkError(err)
	checkResponseCode(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	data := respType{}
	err = json.Unmarshal(body, &data)
	checkError(err)
	items := []Item{}
	for _, v := range data.V {
		dataSet := strings.Split(v, "|")
		if dataSet[1] != "syncwatch" {
			items = append(items, getItemFromRawData(dataSet[2]))
		}
	}
	return items
}

func getItemFromRawData(rawString string) Item {
	item := Item{}
	links := linkify.Links(rawString)
	for _, link := range links {
		linkString := rawString[link.Start:link.End]
		item.Url = linkString
		item.Timestamp = time.Now()
		item.Source = "insagirl"
		item.TextContent =  strings.Trim(strings.Replace(rawString, linkString, "", -1), " ")
	}
	return item
}
