package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/handler"
)

func InventoryRoute(router *gin.Engine) {
	inventoryRoutes := router.Group("api/v1/inventory")
	{
		inventoryRoutes.POST("", handler.CreateInventory())
		inventoryRoutes.GET("/:inventoryId", handler.GetSingleInventory())
		inventoryRoutes.GET("/owner/:ownerId", handler.GetInventories())
	}
}
