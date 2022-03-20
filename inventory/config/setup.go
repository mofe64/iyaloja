package config

import (
	"context"
	"github.com/mofe64/iyaloja/inventory/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var DATABASE *mongo.Client = ConnectDB()

func ConnectDB() *mongo.Client {
	// Create our context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create Mongo Client and connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		util.ApplicationLog.Fatalln("Error connecting to Mongodb %v", err)
	}

	// Defer disconnection of client if we are not in test mode
	if EnvProfile() != "test" {
		defer func() {
			util.ApplicationLog.Println("Mongo disconnecting from client ")
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
	}

	// Ping Mongo db Database
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		util.ApplicationLog.Fatalln(err)
	}
	util.ApplicationLog.Println("Connected to MongoDB ")
	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(EnvDatabaseName()).Collection(collectionName)
	return collection
}
