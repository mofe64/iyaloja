package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Inventory struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	InventoryTags []string           `json:"inventoryTags,omitempty"`
	Description   string             `json:"description,omitempty"`
	Name          string             `json:"name" validate:"required"`
	DateCreated   time.Time          `json:"dateCreated"`
	DateModified  time.Time          `json:"dateModified"`
	ItemCount     int64              `json:"itemCount" default:"0"`
}
