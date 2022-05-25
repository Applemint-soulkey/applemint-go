package crawl

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://"+os.Getenv("env_mongo_id")+":"+os.Getenv("env_mongo_pwd")+"@cluster0.ptfrm.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	checkError(err)

	return client
}

func getHistory(dbclient *mongo.Client) []string {
	coll := dbclient.Database("Item").Collection("History")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	checkError(err)

	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	checkError(err)

	
	var history []string
	for _, result := range results {
		history = append(history, result["url"].(string))
	}
	
	return history
}

func filterItems(dbclient *mongo.Client, items []Item) ([]Item) {
	history := getHistory(dbclient)
	var result []Item
	for _, item := range items {
		if !Contains(history, item.Url) {
			result = append(result, item)
		}
	}
	return result
}

func makeBsonSet(item []Item) []interface{} {
	bsonSet := make([]interface{}, len(item))
	log.Println("Make bson set")
	for i, item := range item {
		data, err := bson.Marshal(item)
		checkError(err)
		bsonSet[i] = data
	}
	return bsonSet
}

func InsertItems(items []Item) int {
	dbclient := ConnectDB()
	coll_new := dbclient.Database("Item").Collection("New")
	coll_history := dbclient.Database("Item").Collection("History")
	filtered := filterItems(dbclient, items)
	itemBsonSet := makeBsonSet(filtered)

	if len(itemBsonSet) > 0 {
		result_item, err := coll_new.InsertMany(context.TODO(), itemBsonSet)
		log.Println("Inserted items to Collection::New")
		checkError(err)
		_, err = coll_history.InsertMany(context.TODO(), itemBsonSet)
		log.Println("Inserted items to Collection::History ")
		checkError(err)

		log.Printf("Insert Conunt: %d\n" , len(result_item.InsertedIDs))
		return len(result_item.InsertedIDs) 
	}

	return 0
}
