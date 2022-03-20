package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/data/model"
	"github.com/mofe64/iyaloja/inventory/dto/response"
	"net/http"
	"time"
)

var inventoryCollection = config.GetCollection(config.DATABASE, "inventories")

func createInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inventory model.Inventory
		if err := c.BindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Status:    http.StatusBadRequest,
				Message:   err.Error(),
				Timestamp: time.Now(),
				Data:      gin.H{},
			})
			return
		}
		saveResult, err := inventoryCollection.InsertOne(ctx, inventory)
		if err != nil {
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
