package crud

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func ClearCollection(coll string) int64 {
	dbclient := connectDB()
	coll_item := dbclient.Database("Item").Collection(coll)
	result, err := coll_item.DeleteMany(context.TODO(), bson.M{})
	checkError(err)

	return result.DeletedCount
}

func GetItem(itemId string, collectionName string) (Item, error) {
	// itemId Length Check
	if len(itemId) != 24 {return Item{}, errors.New("itemId length error")}

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
	return result.ModifiedCount
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

func MoveItem(itemId string, coll_origin string, coll_dest string) error {

	// itemId Length Check
	if len(itemId) != 24 {return errors.New("itemId length error")}

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

	checkError(err)
	return err
}