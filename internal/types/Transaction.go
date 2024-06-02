package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	GoodID       primitive.ObjectID `json:"goodID" bson:"goodID"`
	GoodName     string             `json:"goodName" bson:"goodName"`
	GoodUnit     string             `json:"goodUnit" bson:"goodUnit"`
	Quantity     float64            `json:"quantity" bson:"quantity"`
	Price        float64            `json:"price" bson:"price"`
	BoughtFromID primitive.ObjectID `json:"boughtFromID,omitempty" bson:"boughtFromID,omitempty"`
	BoughtFrom   string             `json:"boughtFrom,omitempty" bson:"boughtFrom,omitempty"`
	SoldToID     primitive.ObjectID `json:"soldToID,omitempty" bson:"soldToID,omitempty"`
	SoldTo       string             `json:"soldTo,omitempty" bson:"soldTo,omitempty"`
	Type         TransactionType    `json:"type" bson:"type"`
	Payment      bool               `json:"payment" bson:"payment"`
	UserID       primitive.ObjectID `json:"userID" bson:"userID"`
	CreationTime primitive.DateTime `json:"creationTime" bson:"creationTime"`
}

type TransactionType string

const (
	Sold   TransactionType = "sold"
	Bought TransactionType = "bought"
)
