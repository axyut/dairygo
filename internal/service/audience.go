package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type AudienceService struct {
	global     *db.Collections
	collection db.Collection
	ctx        context.Context
}

func (s *AudienceService) NewAudience(name string, contact string, userID string, email string, toPay float64, toReceive float64, paid float64) error {
	audience := *s.collection
	// goods := *s.global.Goods
	res, err := audience.InsertOne(s.ctx, bson.M{
		"name":      name,
		"contact":   contact,
		"email":     email,
		"toPay":     toPay,
		"toReceive": toReceive,
		"paid":      paid,
		"userID":    userID,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func (s *AudienceService) GetAudience(UserID string) error {
	audience := *s.collection
	res, err := audience.Find(s.ctx, bson.M{"userID": UserID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
