package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/data/model"
	"github.com/mofe64/iyaloja/inventory/dto/response"
	"github.com/mofe64/iyaloja/inventory/handler"
	"github.com/mofe64/iyaloja/inventory/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

var objIds = []primitive.ObjectID{
	primitive.NewObjectID(),
	primitive.NewObjectID(),
	primitive.NewObjectID(),
	primitive.NewObjectID(),
}

func populateDB(ctx *gin.Context) ([]interface{}, error) {
	inventoryObjs := []interface{}{
		model.Inventory{
			Name:        "inventory1",
			Description: "inventory1",
			OwnerId:     "1",
			Id:          objIds[0],
		},
		model.Inventory{
			Name:        "inventory2",
			Description: "inventory2",
			OwnerId:     "1",
			Id:          objIds[1],
		},
		model.Inventory{
			Name:        "inventory3",
			Description: "inventory3",
			OwnerId:     "2",
			Id:          objIds[2],
		},
		model.Inventory{
			Name:        "inventory4",
			Description: "inventory4",
			OwnerId:     "3",
			Id:          objIds[3],
		},
	}

	inventoryCollection := config.GetCollection(config.DATABASE, "inventories")
	insertManyResult, err := inventoryCollection.InsertMany(ctx, inventoryObjs)
	if err != nil {
		return nil, err
	}
	return insertManyResult.InsertedIDs, nil

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
		"ownerId":     "1234",
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
		Method: "POST",
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

func Test_getSingleInventory(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	_, dbErr := populateDB(ctx)
	if dbErr != nil {
		t.Fatalf("Error performing insert op %v\n", dbErr)
	}
	queryIdByteSlice, _ := objIds[0].MarshalText()
	queryId := string(queryIdByteSlice)
	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
	}

	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "inventoryId",
		Value: queryId,
	})
	getSingleInventory := handler.GetSingleInventory()
	getSingleInventory(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	var res response.APIResponse
	responseString := w.Body.String()
	err := json.Unmarshal([]byte(responseString), &res)
	if err != nil {
		util.ApplicationLog.Printf("ERROR UNMARSHALLING RESPONSE %v\n", err)
	}
	assert.Equal(t, http.StatusOK, res.Status)
	assert.NotNil(t, res.Data)
	assert.Equal(t, "Success", res.Message)
}

func Test_getSingleInventory_returns404IfIdIncorrect(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	_, dbErr := populateDB(ctx)
	if dbErr != nil {
		t.Fatalf("Error performing insert op %v\n", dbErr)
	}

	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
	}

	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "inventoryId",
		Value: "12345",
	})
	getSingleInventory := handler.GetSingleInventory()
	getSingleInventory(ctx)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var res response.APIResponse
	responseString := w.Body.String()
	err := json.Unmarshal([]byte(responseString), &res)
	if err != nil {
		util.ApplicationLog.Printf("ERROR UNMARSHALLING RESPONSE %v\n", err)
	}
	assert.Equal(t, http.StatusNotFound, res.Status)
}

func Test_getUserInventories(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
	}
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "ownerId",
		Value: "1",
	})
	_, dbErr := populateDB(ctx)
	if dbErr != nil {
		t.Fatalf("Error performing insert op %v\n", dbErr)
	}
	getInventoriesHandler := handler.GetInventories()
	getInventoriesHandler(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	var res response.APIResponse
	responseString := w.Body.String()
	err := json.Unmarshal([]byte(responseString), &res)
	if err != nil {
		util.ApplicationLog.Printf("ERROR UNMARSHALLING RESPONSE %v\n", err)
	}

	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, float64(2), res.Data["inventoryCount"])
	t.Cleanup(func() {
		CleanUpDbOps(ctx)
	})

}

func Test_deleteInventory(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	_, dbErr := populateDB(ctx)
	if dbErr != nil {
		t.Fatalf("Error performing insert op %v\n", dbErr)
	}
	queryIdByteSlice, _ := objIds[0].MarshalText()
	queryId := string(queryIdByteSlice)
	ctx.Request = &http.Request{
		Method: "DELETE",
		Header: make(http.Header),
	}

	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "inventoryId",
		Value: queryId,
	})
	deleteInventoryHandler := handler.DeleteInventory()
	deleteInventoryHandler(ctx)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
