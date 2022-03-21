package config

import (
	"context"
	"github.com/mofe64/iyaloja/inventory/util"
	"go.mongodb.org/mongo-driver/bson"
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
	// Ping Mongo db Database
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		util.ApplicationLog.Fatalln(err)
	}
	util.ApplicationLog.Println("Connected to MongoDB ")
	return client
}
func SetUpTablesAndIndexes(configProperties map[string]map[string]bool) {
	if DATABASE == nil {
		panic("Database has not been initialized yet ...")
	}
	for collection := range configProperties {
		collectionDetails := configProperties[collection]
		for field := range collectionDetails {
			CreateIndex(collection, field, collectionDetails[field])
		}
	}
}

// GetCollection - helper function to retrieve collection from database
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(EnvDatabaseName()).Collection(collectionName)
	return collection
}

// CreateIndex - creates an index for a specific field in a collection
func CreateIndex(collectionName string, field string, unique bool) bool {

	// 1. Let's define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := DATABASE.Database(EnvDatabaseName()).Collection(collectionName)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		util.ApplicationLog.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
