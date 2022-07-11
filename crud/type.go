package crud

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	TextContent string             `json:"text_content" bson:"text_content"`
	Url         string             `json:"url" bson:"url"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Domain      string             `json:"domain" bson:"domain"`
	Tags        []string           `json:"tags" bson:"tags"`
	Path        string             `json:"path" bson:"path"`
	Source      string             `json:"source" bson:"source"`
}

type BookmarkInfo struct {
	Path  string `bson:"_id"`
	Count int64  `bson:"count"`
}

type GroupInfo struct {
	Domain string `bson:"_id"`
	Count  int64  `bson:"count"`
}

type GalleryItem struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Text      string             `json:"text" bson:"text"`
	Link      string             `json:"link" bson:"link"`
	Origin    string             `json:"origin" bson:"origin"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type GalleryResponse struct {
	Items  []GalleryItem `json:"item"`
	Count  int64         `json:"count"`
	Cursor int64         `json:"cursor"`
}
