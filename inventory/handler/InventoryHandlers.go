package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/data/model"
	"github.com/mofe64/iyaloja/inventory/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var inventoryCollection = config.GetCollection(config.DATABASE, "inventories")
var validate = validator.New()

func CreateInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inventory model.Inventory
		if err := c.ShouldBindJSON(&inventory); err != nil {
			util.ApplicationLog.Printf("Error binding Json Obj %v\n", err)
			util.GenerateJSONResponse(c, http.StatusBadRequest, err.Error(), gin.H{})
			return
		}
		if validationErr := validate.Struct(&inventory); validationErr != nil {
			util.ApplicationLog.Println("validation error")
			util.GenerateBadRequestResponse(c, validationErr.Error())
			return
		}

		inventory.Id = primitive.NewObjectID()
		inventory.DateCreated = time.Now()
		util.ApplicationLog.Printf("Inventory obj %v\n", inventory)

		saveResult, err := inventoryCollection.InsertOne(ctx, inventory)
		if err != nil {
			util.ApplicationLog.Printf("Error Saving Obj %v\n", err)
			util.GenerateInternalServerErrorResponse(c, err.Error())
			return
		}
		util.GenerateJSONResponse(c, http.StatusCreated, "Success", gin.H{
			"inventory": saveResult,
		})
	}
}
func GetInventories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ownerId := c.Param("ownerId")
		util.ApplicationLog.Println("Owner id param " + ownerId)

		inventories := []model.Inventory{}
		queryOptions := options.Find()
		queryOptions.SetSort(bson.D{{"DateCreated", -1}})

		cursor, err := inventoryCollection.Find(ctx, bson.M{"ownerId": ownerId}, queryOptions)
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				util.ApplicationLog.Printf("Error closing query cursor %v\n", err)
			}
		}(cursor, ctx)

		if err != nil {
			util.GenerateInternalServerErrorResponse(c, err.Error())
			return
		}

		for cursor.Next(ctx) {
			var result model.Inventory
			err := cursor.Decode(&result)
			util.ApplicationLog.Println("Decoded inventory result")
			util.ApplicationLog.Println(result)
			if err != nil {
				util.GenerateInternalServerErrorResponse(c, err.Error())
				return
			}
			inventories = append(inventories, result)
		}

		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}

		util.GenerateJSONResponse(c, http.StatusOK, "Success", gin.H{
			"inventoryCount": len(inventories),
			"inventories":    inventories,
		})
	}
}

func GetSingleInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		inventoryId := c.Param("inventoryId")
		util.ApplicationLog.Println("received inventory id " + inventoryId)
		objId, _ := primitive.ObjectIDFromHex(inventoryId)
		var inventory model.Inventory
		filter := bson.D{{"_id", objId}}
		err := inventoryCollection.FindOne(ctx, filter).Decode(&inventory)
		if err == mongo.ErrNoDocuments {
			util.GenerateJSONResponse(c, http.StatusNotFound, "Not Found", gin.H{})
			return
		} else if err != nil {
			util.GenerateInternalServerErrorResponse(c, err.Error())
		}

		util.GenerateJSONResponse(c, http.StatusOK, "Success", gin.H{
			"inventory": inventory,
		})
	}
}
