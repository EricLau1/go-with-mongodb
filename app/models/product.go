package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Product struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Price     float64       `bson:"price" json:"price"`
	Quantity  int           `bson:"quantity" json:"quantity"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
