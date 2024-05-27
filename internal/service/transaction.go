package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionService struct {
	service    *Service
	collection db.Collection
}

func (s *TransactionService) NewTransaction(ctx context.Context, trans types.Transaction) (insertedTrans types.Transaction, err error) {
	transaction := *s.collection
	res, err := transaction.InsertOne(ctx, bson.M{
		"goodID":     trans.GoodID,
		"quantity":   trans.Quantity,
		"price":      trans.Price,
		"boughtFrom": trans.BoughtFrom,
		"soldTo":     trans.SoldTo,
		"userID":     trans.UserID,
		"type":       trans.Type,
		"payment":    trans.Payment,
	})
	if err != nil {
		s.service.logger.Error("Error while inserting new transaction", err)
		return trans, err
	}
	err = transaction.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&trans)
	if err != nil {
		s.service.logger.Error("Error while decoding new transaction", err)
		return trans, err
	}
	return trans, nil
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
