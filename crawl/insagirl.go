package crawl

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
	ahocorasick "github.com/petar-dambovaliev/aho-corasick"
	"go.mongodb.org/mongo-driver/bson"
	"mvdan.cc/xurls/v2"
)

type respType struct {
	V []string `json:"v"`
}

func isGalleryItem(item Item) bool {
	directKeywords := []string{
		".jpg",
		".png",
		".gif",
		".jpeg",
		".webp",
		".bmp",
		".svg",
		".mp4",
		".mp3",
		".wav",
		".flac",
		".ogg",
		".opus",
		".aac",
		".wma",
		".m4a",
		".m4v",
		".3gp",
		".3g2",
		".avi",
		".mov",
		".mkv",
		".wmv",
	}

	if item.Domain == "imgur.com" {
		return true
	}

	for _, keyword := range directKeywords {
		if strings.Contains(item.Url, keyword) {
			return true
		}
	}

	return false
}

func CrawlISG() int {
	targetList := []string{
		"http://insagirl-hrm.appspot.com/json2/1/1/2/",
		"http://insagirl-hrm.appspot.com/json2/2/1/2/",
	}

	items := []Item{}
	imageQueueItems := []Item{}
	// Get Items
	log.Print("Get Items")
	for _, targetURL := range targetList {
		newItems, imageItems := getIsgData(targetURL)
		items = append(items, newItems...)
		imageQueueItems = append(imageQueueItems, imageItems...)
	}

	log.Print("Insert Items")
	insertedCount := InsertItems(items, "new")
	log.Printf("Inserted %d Items", insertedCount)
	log.Print("Insert Image Queue Items")
	insertedCount = InsertItems(imageQueueItems, "image-queue")
	log.Printf("Inserted %d Image Queue Items", insertedCount)

	return insertedCount
}

func getIsgData(url string) ([]Item, []Item) {
	resp, err := http.Get(url)
	checkError(err)
	checkResponseCode(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	data := respType{}
	err = json.Unmarshal(body, &data)
	checkError(err)

	ignoreList := getIgnoreListFromDB()

	items := []Item{}
	imageQueueItems := []Item{}
	for _, v := range data.V {
		dataSet := strings.Split(v, "|")
		item, err := getItemFromRawData(dataSet[2])
		if !containIgnoreWord(dataSet, ignoreList) && err == nil {
			if isGalleryItem(item) {
				imageQueueItems = append(imageQueueItems, item)
			} else {
				items = append(items, item)
			}
		}
	}
	return items, imageQueueItems
}

func getIgnoreListFromDB() []string {
	dbclient := ConnectDB()
	coll_ignore := dbclient.Database("Settings").Collection("ignore")
	cursor, err := coll_ignore.Find(context.TODO(), bson.D{})
	checkError(err)

	var ignoreData []bson.M
	err = cursor.All(context.TODO(), &ignoreData)
	checkError(err)

	var ignoreList []string
	for _, ignore := range ignoreData {
		ignoreList = append(ignoreList, ignore["data"].(string))
	}

	dbclient.Disconnect(context.TODO())

	return ignoreList
}

func containIgnoreWord(dataSet []string, ignoreList []string) bool {
	builder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  true,
		MatchKind:            ahocorasick.LeftMostLongestMatch,
		DFA:                  true,
	})
	ac := builder.Build(ignoreList)
	matches := ac.FindAll(dataSet[2])

	if dataSet[1] == "syncwatch" {
		return true
	} else if len(matches) > 0 {
		return true
	} else {
		return false
	}
}

func getItemFromRawData(rawString string) (Item, error) {
	item := Item{}
	rx := xurls.Strict()
	matches := rx.FindAllString(rawString, -1)
	if len(matches) == 0 {
		return item, errors.New("no url found")
	}
	for _, match := range matches {
		item.Url = match
		item.Timestamp = time.Now()
		item.Domain = domainutil.Domain(match)
		item.Source = "insagirl"
		item.Tags = []string{}
		item.Path = ""
		item.TextContent = strings.Trim(strings.Replace(rawString, match, "", -1), " ")
	}
	return item, nil
}
