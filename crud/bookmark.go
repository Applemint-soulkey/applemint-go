package crud

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookmarkInfo struct {
	Path  string `bson:"_id"`
	Count int64  `bson:"count"`
}

func SendToBookmark(item Item, coll_origin string, path string) (Item, error) {
	// connect to raindrop
	log.Print("SendToBookmark")
	item.Path = path
	target_id_string := item.ID.Hex()
	log.Print("target_id_string: " + target_id_string)
	UpdateItem(target_id_string, coll_origin, item)

	if coll_origin != "bookmark" {
		err := MoveItem(target_id_string, coll_origin, "bookmark")
		checkError(err)
	}

	resultItem, err := GetItem(target_id_string, "bookmark")
	checkError(err)

	return resultItem, err
}

func GetBookmarkList() ([]BookmarkInfo, error) {
	log.Print("GetBookmarkList")
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection("bookmark")

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$path"},
			{"count", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	// Set Sort Stage
	sortStage := bson.D{
		{"$sort", bson.D{
			{"count", -1},
		}},
	}

	// Aggregate Stage From Collection
	dbCursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage, sortStage})
	checkError(err)

	// Load Domain Count
	var results []bson.M
	if err := dbCursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var pathInfos []BookmarkInfo = make([]BookmarkInfo, 0)
	for _, result := range results {
		pathInfo := BookmarkInfo{}
		bytes, _ := bson.Marshal(result)
		bson.Unmarshal(bytes, &pathInfo)
		pathInfos = append(pathInfos, pathInfo)
	}
	dbclient.Disconnect(context.TODO())

	return pathInfos, nil
}
