package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/db"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionService struct {
	global     *db.Collections
	collection db.Collection
	ctx        context.Context
}

func (s *TransactionService) NewTransaction(userID string, goodsID string, quantity int) error {
	transaction := *s.collection
	res, err := transaction.InsertOne(s.ctx, bson.M{
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

func (s *TransactionService) GetTransaction(userID string) error {
	transaction := *s.collection
	res, err := transaction.Find(s.ctx, bson.M{"userID": userID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
