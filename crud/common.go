package crud

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Id          string    `json:"id" bson:"_id"`
	TextContent string    `json:"text_content" bson:"text_content"`
	Url         string    `json:"url" bson:"url"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Domain      string    `json:"domain" bson:"domain"`
	Tags        []string  `json:"tags" bson:"tags"`
	Path        string    `json:"path" bson:"path"`
	Source      string    `json:"source" bson:"source"`
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func connectDB() *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + os.Getenv("ENV_MONGO_ID") + ":" + os.Getenv("ENV_MONGO_PWD") + "@cluster0.ptfrm.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	checkError(err)

	return client
}
