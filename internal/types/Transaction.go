package types

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
	Payment    bool               `json:"payment" bson:"payment"`
	UserID     primitive.ObjectID `json:"userID" bson:"userID"`
}

type TransactionType string

const (
	Sold   TransactionType = "sold"
	Bought TransactionType = "bought"
)

type Transaction_Client struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	GoodName   string             `json:"goodName" bson:"goodName"`
	GoodUnit   string             `json:"goodUnit" bson:"goodUnit"`
	Quantity   string             `json:"quantity" bson:"quantity"`
	Price      string             `json:"price" bson:"price"`
	BoughtFrom string             `json:"boughtFrom,omitempty" bson:"boughtFrom,omitempty"`
	SoldTo     string             `json:"soldTo,omitempty" bson:"soldTo,omitempty"`
	Payment    bool               `json:"payment" bson:"payment"`
}
