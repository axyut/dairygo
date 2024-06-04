package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Production struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	ChangeGoodID     primitive.ObjectID `json:"changeGoodID" bson:"changeGoodID"`
	ChangeQuantity   float64            `json:"changeQuantity" bson:"changeQuantity"`
	ChangeGoodName   string             `json:"changeGoodName" bson:"changeGoodName"`
	ChangeGoodUnit   string             `json:"changeGoodUnit" bson:"changeGoodUnit"`
	ChangePrice      float64            `json:"changePrice" bson:"changePrice"`
	ProducedGoodID   primitive.ObjectID `json:"producedGoodID" bson:"producedGoodID"`
	ProducedQuantity float64            `json:"producedQuantity" bson:"producedQuantity"`
	ProducedGoodName string             `json:"producedGoodName" bson:"producedGoodName"`
	ProducedGoodUnit string             `json:"producedGoodUnit" bson:"producedGoodUnit"`
	ProducedPrice    float64            `json:"producedPrice" bson:"producedPrice"`
	CreationTime     primitive.DateTime `json:"creationTime" bson:"creationTime"`
	UserID           primitive.ObjectID `json:"userID" bson:"userID"`
}
