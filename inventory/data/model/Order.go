package model

import (
	"github.com/mofe64/iyaloja/inventory/data/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	DateCreated   time.Time          `json:"dateCreated"`
	DateFulfilled time.Time          `json:"dateFulfilled"`
	Type          enum.OrderType     `json:"type"`
	Status        enum.OrderStatus   `json:"status"`
	Amount        string             `json:"amount"`
}
