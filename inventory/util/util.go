package util

import (
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/dto/response"
	"log"
	"os"
	"time"
)

var ApplicationLog = log.New(os.Stdout, "inventory-service ", log.LstdFlags)

func GenerateJSONResponse(c *gin.Context, statusCode int, message string, data map[string]interface{}) {
	c.JSON(statusCode, response.APIResponse{
		Status:    statusCode,
		Message:   message,
		Timestamp: time.Now(),
		Data:      data,
	})
}
