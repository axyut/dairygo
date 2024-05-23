package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	GoodID     primitive.ObjectID `json:"goodID" bson:"goodID"`
	Quantity   float64            `json:"quantity" bson:"quantity"`
	Price      float64            `json:"price" bson:"price"`
	BoughtFrom primitive.ObjectID `json:"boughtFrom,omitempty" bson:"boughtFrom,omitempty"`
	SoldTo     primitive.ObjectID `json:"soldTo,omitempty" bson:"soldTo,omitempty"`
	Type       TransactionType    `json:"type" bson:"type"`
	UserID     primitive.ObjectID `json:"userID" bson:"userID"`
}

type TransactionType string

const (
	Import TransactionType = "import"
	Export TransactionType = "export"
)
