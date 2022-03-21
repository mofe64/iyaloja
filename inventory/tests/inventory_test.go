package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/dto/response"
	"github.com/mofe64/iyaloja/inventory/handler"
	"github.com/mofe64/iyaloja/inventory/util"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper method to test post requests, we pass a test context as well as the content
func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	requestBody, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
}

func CleanUpDbOps(ctx *gin.Context) {
	err := config.GetCollection(config.DATABASE, "inventories").Database().Drop(ctx)
	if err != nil {
		return
	}
	err = config.DATABASE.Disconnect(ctx)
	if err != nil {
		return
	}
}

func Test_createInventoryHandler(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	reqBody := gin.H{
		"name":        "test inventory",
		"description": "test description",
	}

	MockJsonPost(ctx, reqBody)
	createInventoryHandler := handler.CreateInventory()
	createInventoryHandler(ctx)

	var res response.APIResponse
	responseString := w.Body.String()
	err := json.Unmarshal([]byte(responseString), &res)
	if err != nil {
		util.ApplicationLog.Printf("ERROR UNMARSHALLING RESPONSE %v\n", err)
	}
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusCreated, res.Status)
	assert.Equal(t, "Success", res.Message)
	assert.NotNil(t, res.Timestamp)
	assert.NotNil(t, res.Data)
	t.Cleanup(func() {
		CleanUpDbOps(ctx)
	})
}
func Test_createInventoryHandler_failsWhenRequiredFieldMissing(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	reqBody := gin.H{}
	MockJsonPost(ctx, reqBody)
	createInventoryHandler := handler.CreateInventory()
	createInventoryHandler(ctx)

	var res response.APIResponse
	responseString := w.Body.String()
	err := json.Unmarshal([]byte(responseString), &res)
	if err != nil {
		util.ApplicationLog.Printf("ERROR UNMARSHALLING RESPONSE %v\n", err)
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusBadRequest, res.Status)
	assert.NotNil(t, res.Timestamp)
	t.Cleanup(func() {
		CleanUpDbOps(ctx)
	})

}
