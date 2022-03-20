package model

import "time"

type Order struct {
	Date time.Time `json:"date"`
}
