package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type AudienceService struct {
	service    *Service
	collection db.Collection
}

func (s *AudienceService) NewAudience(ctx context.Context, name string, contact string, userID string, email string, toPay float64, toReceive float64, paid float64) error {
	audience := *s.collection
	res, err := audience.InsertOne(ctx, bson.M{
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

func (s *AudienceService) GetAudience(ctx context.Context, UserID string) error {
	audience := *s.collection
	res, err := audience.Find(ctx, bson.M{"userID": UserID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
