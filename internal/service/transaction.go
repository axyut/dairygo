package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionService struct {
	service    *Service
	collection db.Collection
}

func (s *TransactionService) NewTransaction(ctx context.Context, userID string, goodsID string, quantity int) error {
	transaction := *s.collection
	res, err := transaction.InsertOne(ctx, bson.M{
		"userID":   userID,
		"goodsID":  goodsID,
		"quantity": quantity,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, userID string) error {
	transaction := *s.collection
	res, err := transaction.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
