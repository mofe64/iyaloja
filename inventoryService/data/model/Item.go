package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type InventoryItem struct {
	Id                        primitive.ObjectID `json:"id" bson:"_id"`
	ItemName                  string             `json:"itemName"`
	ItemCount                 int64              `json:"itemCount"`
	Inventory                 primitive.ObjectID `json:"inventory"`
	ItemNotificationThreshold int                `json:"itemNotificationThreshold,omitempty"`
	InStock                   bool               `json:"inStock"`
	CostPriceInfo             PriceInfo          `json:"costPriceInfo"`
	SellPriceInfo             PriceInfo          `json:"sellPriceInfo"`
}

type PriceInfo struct {
	BaseAmount int64  `json:"baseAmount"`
	Unit       string `json:"unit"`
	// TODO add currency enum
}
