package crud

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetGalleryItems(cursor int64) (GalleryResponse, error) {
	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection("gallery")
	findOption := options.Find().SetSort(bson.M{"createdAt": -1}).SetSort(bson.M{"_id": -1}).SetSkip(cursor).SetLimit(PAGE_SIZE)
	dbCursor, err := coll.Find(context.TODO(), bson.M{}, findOption)
	if err != nil {
		return GalleryResponse{}, err
	}

	// Get Items
	var items []GalleryItem = make([]GalleryItem, 0)
	for dbCursor.Next(context.TODO()) {
		var item GalleryItem
		err = dbCursor.Decode(&item)
		if err != nil {
			return GalleryResponse{}, err
		}
		items = append(items, item)
	}

	// Get Items Count
	count, err := coll.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return GalleryResponse{}, err
	}

	// Make Response
	response := GalleryResponse{
		Items:  items,
		Count:  count,
		Cursor: cursor,
	}

	// response := make(map[string]interface{})
	// response["items"] = items
	// response["count"] = count
	// response["cursor"] = cursor

	dbclient.Disconnect(context.TODO())

	return response, nil
}
