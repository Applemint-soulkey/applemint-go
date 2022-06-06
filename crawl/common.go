package crawl

import "time"

var PAGE_SIZE = 1
var BASE_URL_BP = "https://v12.battlepage.com"

type Item struct {
	TextContent string    `json:"text_content" bson:"text_content"`
	Url         string    `json:"url" bson:"url"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Domain		string	  `json:"domain" bson:"domain"`
	Tags		[]string  `json:"tags" bson:"tags"`
	Path		string	  `json:"path" bson:"path"`
	Source      string    `json:"source" bson:"source"`
}