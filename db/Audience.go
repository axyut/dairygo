package db

type Audience struct {
	ID        string  `json:"id" bson:"_id"`
	Name      string  `json:"name" bson:"name"`
	Contact   string  `json:"contact,omitempty" bson:"contact,omitempty"`
	Email     string  `json:"email,omitempty" bson:"email,omitempty"`
	ToPay     float64 `json:"toPay,omitempty" bson:"toPay,omitempty"`
	Paid      float64 `json:"paid,omitempty" bson:"paid,omitempty"`
	ToReceive float64 `json:"toReceive,omitempty" bson:"toReceive,omitempty"`
	UserID    string  `json:"userID" bson:"userID"`
}
