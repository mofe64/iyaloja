package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/handler"
)

func InventoryRoute(router *gin.Engine) {
	inventoryRoutes := router.Group("/inventory")
	{
		inventoryRoutes.POST("", handler.CreateInventory())
	}
}
