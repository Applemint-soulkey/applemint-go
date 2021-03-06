package crud

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollectionList() ([]string, error) {
	// Connect to DB
	dbclient := connectDB()
	db := dbclient.Database("Item")
	collection_list, err := db.ListCollectionNames(context.TODO(), bson.M{})
	checkError(err)

	dbclient.Disconnect(context.TODO())

	return collection_list, nil
}

func GetCollectionInfo(collectionName string) (int64, []GroupInfo, error) {
	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection(collectionName)

	// Set Group Stage
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$domain"},
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

	var groupInfos []GroupInfo = make([]GroupInfo, 0)

	etcGroupInfo := GroupInfo{
		Domain: "etc",
		Count:  0,
	}
	var totalCount int64 = 0
	for _, result := range results {
		groupInfo := GroupInfo{}
		bytes, _ := bson.Marshal(result)
		bson.Unmarshal(bytes, &groupInfo)
		totalCount += groupInfo.Count
		if groupInfo.Count > GROUP_SIZE {
			groupInfos = append(groupInfos, groupInfo)
		} else {
			etcGroupInfo.Count += groupInfo.Count
		}
	}
	groupInfos = append(groupInfos, etcGroupInfo)

	dbclient.Disconnect(context.TODO())

	return totalCount, groupInfos, nil
}

func ClearCollection(coll string) int64 {
	dbclient := connectDB()
	coll_item := dbclient.Database("Item").Collection(coll)
	result, err := coll_item.DeleteMany(context.TODO(), bson.M{})
	checkError(err)

	dbclient.Disconnect(context.TODO())

	return result.DeletedCount
}

func GetItems(collectionName string, cursor int64, domain string, path string) ([]Item, error) {
	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection(collectionName)

	// Set Filter
	var domainFilterMap bson.M = bson.M{}

	if domain == "etc" {
		_, groupInfos, err := GetCollectionInfo(collectionName)
		if err != nil {
			return nil, err
		}
		var filterDomain []string = make([]string, 0)

		for _, groupInfo := range groupInfos {
			if groupInfo.Count > GROUP_SIZE {
				filterDomain = append(filterDomain, groupInfo.Domain)
			}
		}
		domainFilterMap = bson.M{"domain": bson.M{"$nin": filterDomain}}
	} else if domain != "" {
		domainFilterMap["domain"] = domain
	}

	if path != "" {
		domainFilterMap["path"] = path
	}

	log.Println(domainFilterMap)

	// Make Stage
	stage := bson.A{
		bson.D{
			{"$match",
				domainFilterMap,
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"timestamp", -1},
					{"_id", -1},
				},
			},
		},
		bson.D{{"$skip", cursor}},
		bson.D{{"$limit", PAGE_SIZE}},
	}
	dbCursor, err := coll.Aggregate(context.TODO(), stage)
	checkError(err)

	// Get Items
	var items []Item = make([]Item, 0)
	for dbCursor.Next(context.TODO()) {
		var item Item
		err = dbCursor.Decode(&item)
		checkError(err)
		items = append(items, item)
	}

	dbclient.Disconnect(context.TODO())

	return items, nil
}

func GetItem(itemId string, collectionName string) (Item, error) {
	// itemId Length Check
	if len(itemId) != 24 {
		return Item{}, errors.New("itemId length error")
	}

	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection(collectionName)

	// itemId to ObjectId
	bsonItemId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		return Item{}, errors.New("itemId to ObjectId error")
	}

	// Get
	item := Item{}
	err = coll.FindOne(context.TODO(), bson.M{"_id": bsonItemId}).Decode(&item)

	dbclient.Disconnect(context.TODO())

	return item, err
}

func UpdateItem(itemId string, collectionName string, item Item) int64 {
	// Connect to DB
	dbclient := connectDB()
	coll := dbclient.Database("Item").Collection(collectionName)
	bsonItemId, err := primitive.ObjectIDFromHex(itemId)

	// Update
	result, err := coll.UpdateOne(context.TODO(), bson.M{"_id": bsonItemId}, bson.M{"$set": item})
	checkError(err)

	dbclient.Disconnect(context.TODO())

	return result.ModifiedCount
}

func DeleteItem(itemId string, collectionName string) int64 {
	// itemId Length Check
	if len(itemId) != 24 {
		return 0
	}

	log.Println("DeleteItem: " + itemId)

	// Connect to DB
	dbclient := connectDB()

	// itemId to ObjectId
	bsonItemId, err := primitive.ObjectIDFromHex(itemId)
	checkError(err)

	// Delete
	result, err := dbclient.Database("Item").Collection(collectionName).DeleteOne(context.TODO(), bson.M{"_id": bsonItemId})
	checkError(err)

	dbclient.Disconnect(context.TODO())
	return result.DeletedCount
}

func MoveItem(itemId string, coll_origin string, coll_dest string) error {

	// itemId Length Check
	if len(itemId) != 24 {
		return errors.New("itemId length error")
	}

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
		if err != nil {
			return err
		}

		// Update Timestamp when item is moved to Trash
		if coll_dest == "trash" {
			item.Timestamp = time.Now()
		}

		_, err = dest_coll.InsertOne(context.TODO(), item)
		if err != nil {
			return err
		}
		_, err = origin_coll.DeleteOne(context.TODO(), bson.M{"_id": bsonItemId})
		if err != nil {
			return err
		}
		defer sessionContext.EndSession(sessionContext)
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}
		return nil
	})

	dbclient.Disconnect(context.TODO())

	checkError(err)
	return err
}
