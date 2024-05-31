package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Audience struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Contact   string             `json:"contact,omitempty" bson:"contact,omitempty"`
	ToPay     float64            `json:"toPay,omitempty" bson:"toPay,omitempty"`
	ToReceive float64            `json:"toReceive,omitempty" bson:"toReceive,omitempty"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
}

type UpdateAudience struct {
	Name      string  `json:"name" bson:"name"`
	Contact   string  `json:"contact,omitempty" bson:"contact,omitempty"`
	ToPay     float64 `json:"toPay,omitempty" bson:"toPay,omitempty"`
	ToReceive float64 `json:"toReceive,omitempty" bson:"toReceive,omitempty"`
}
