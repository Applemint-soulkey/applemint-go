package crud

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PAGE_SIZE = 20
const GROUP_SIZE = 10

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
