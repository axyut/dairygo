package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Good struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Unit       string             `json:"unit" bson:"unit"`
	KharidRate float64            `json:"kharidRate" bson:"kharidRate"`
	BikriRate  float64            `json:"bikriRate" bson:"bikriRate"`
	Quantity   float64            `json:"quantity" bson:"quantity"`
	BoughtFrom primitive.ObjectID `json:"boughtFrom,omitempty" bson:"boughtFrom,omitempty"`
	SoldTo     primitive.ObjectID `json:"soldTo,omitempty" bson:"soldTo,omitempty"`
	UserID     primitive.ObjectID `json:"userID" bson:"userID"`
}

type UpdateGood struct {
	Name       string  `json:"name" bson:"name"`
	Unit       string  `json:"unit" bson:"unit"`
	KharidRate float64 `json:"kharidRate" bson:"kharidRate"`
	BikriRate  float64 `json:"bikriRate" bson:"bikriRate"`
	Quantity   float64 `json:"quantity" bson:"quantity"`
}
