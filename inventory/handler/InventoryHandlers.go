package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/data/model"
	"github.com/mofe64/iyaloja/inventory/dto/response"
	"github.com/mofe64/iyaloja/inventory/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

var inventoryCollection = config.GetCollection(config.DATABASE, "inventories")

func CreateInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inventory model.Inventory
		if err := c.BindJSON(&inventory); err != nil {
			util.ApplicationLog.Printf("Error binding Json Obj %v\n", err)
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Status:    http.StatusBadRequest,
				Message:   err.Error(),
				Timestamp: time.Now(),
				Data:      gin.H{},
			})
			return
		}
		inventory.Id = primitive.NewObjectID()
		inventory.DateCreated = time.Now()
		util.ApplicationLog.Printf("Inventory obj %v\n", inventory)

		saveResult, err := inventoryCollection.InsertOne(ctx, inventory)
		if err != nil {
			util.ApplicationLog.Printf("Error Saving Obj %v\n", err)
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				Timestamp: time.Now(),
				Data:      gin.H{},
			})
			return
		}
		c.JSON(http.StatusCreated, response.APIResponse{
			Status:    http.StatusCreated,
			Message:   "Success",
			Timestamp: time.Now(),
			Data: gin.H{
				"inventory": saveResult,
			},
		})
	}
}
