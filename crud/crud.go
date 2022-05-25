package crud

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	TextContent string `json:"text_content" bson:"text_content"` 
	Url   string `json:"url" bson:"url"`
	Path string `json:"path" bson:"path"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Source string `json:"source" bson:"source"`
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func DeleteItem(itemId string, collectionName string) int64 {
	// itemId Length Check
	if len(itemId) != 24 {return 0}

	// Connect to DB
	dbclient := connectDB()

	// itemId to ObjectId
	bsonItemId, err := primitive.ObjectIDFromHex(itemId)
	checkError(err)

	// Delete
	result, err := dbclient.Database("Item").Collection(collectionName).DeleteOne(context.TODO(), bson.M{"_id": bsonItemId})
	checkError(err)
	return result.DeletedCount
}

func MoveItem(itemId string, coll_origin string, coll_dest string) {

	// itemId Length Check
	if len(itemId) != 24 {return}

	// Connect to DB
	dbclient := connectDB()
	origin_coll := dbclient.Database("Item").Collection(coll_origin)
	dest_coll := dbclient.Database("Item").Collection(coll_dest)
	
	// itemId to ObjectId
	bsonItemId, err := primitive.ObjectIDFromHex(itemId)
	checkError(err)

	//transaction
	err = dbclient.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		item := Item{}
		err = origin_coll.FindOne(context.TODO(), bson.M{"_id": bsonItemId}).Decode(&item)
		checkError(err)
		fmt.Println(item)
		_, err = dest_coll.InsertOne(context.TODO(), item)
		checkError(err)
		_, err = origin_coll.DeleteOne(context.TODO(), bson.M{"_id": bsonItemId})
		checkError(err)
		defer sessionContext.EndSession(sessionContext)
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}
		return nil
	})

	checkError(err)
}

func KeepItem(item Item, itemPath string) {
	fmt.Println("keep item")

}

func sendToDropbox(item Item) {
	fmt.Println("send to dropbox")
}

func sendToRaindrop(item Item) {
	fmt.Println("send to raindrop")
}