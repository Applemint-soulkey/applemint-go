package crud

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GalleryItem struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Text      string             `json:"text" bson:"text"`
	Link      string             `json:"link" bson:"link"`
	Origin    string             `json:"origin" bson:"origin"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func GetGalleryItems(cursor int64) (map[string]interface{}, error) {
	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection("gallery")
	findOption := options.Find().SetSort(bson.M{"createdAt": -1}).SetSort(bson.M{"_id": -1}).SetSkip(cursor).SetLimit(PAGE_SIZE)
	dbCursor, err := coll.Find(context.TODO(), bson.M{}, findOption)
	if err != nil {
		return nil, err
	}

	// Get Items
	var items []GalleryItem = make([]GalleryItem, 0)
	for dbCursor.Next(context.TODO()) {
		var item GalleryItem
		err = dbCursor.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	// Get Items Count
	count, err := coll.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Make Response
	response := make(map[string]interface{})
	response["items"] = items
	response["count"] = count
	response["cursor"] = cursor

	dbclient.Disconnect(context.TODO())

	return response, nil
}
