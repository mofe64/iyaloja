package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Inventory struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	InventoryTags []string           `json:"inventoryTags,omitempty" bson:"inventoryTags"`
	Description   string             `json:"description,omitempty" bson:"description"`
	Name          string             `json:"name" validate:"required" bson:"name"`
	DateCreated   time.Time          `json:"dateCreated" bson:"dateCreated"`
	DateModified  time.Time          `json:"dateModified" bson:"dateModified"`
	ItemCount     int64              `json:"itemCount" default:"0" bson:"itemCount"`
	OwnerId       string             `json:"ownerId" validate:"required" bson:"ownerId"`
}
