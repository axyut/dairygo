package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID     `json:"id" bson:"_id"`
	UserName string                 `json:"username" bson:"username"`
	Email    string                 `json:"email" bson:"email"`
	Password string                 `json:"password" bson:"password"`
	Default  map[UserDefault]string `json:"default" bson:"default"`
}

type UserDefault string

const (
	SellGood    UserDefault = "sellgood"
	SellPayment UserDefault = "sellpayment"
)
