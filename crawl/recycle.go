package crawl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func ClearOldData(collectionName string) {
	// Connect to DB
	dbclient := ConnectDB()
	coll := dbclient.Database("Item").Collection(collectionName)

	// Delete Old Data
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := coll.DeleteMany(ctx, bson.D{{"timestamp", bson.D{{"$lt", time.Now().AddDate(0, 0, -3)}}}})
	checkError(err)

	dbclient.Disconnect(context.TODO())
}
