package service

import (
	"context"
	"math"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionService struct {
	service    *Service
	collection db.Collection
}

func (s *TransactionService) InsertTransaction(ctx context.Context, trans types.Transaction) (insertedTrans types.Transaction, err error) {
	transaction := *s.collection
	trans.Price = math.Abs(trans.Price)
	trans.Quantity = math.Abs(trans.Quantity)
	res, err := transaction.InsertOne(ctx, trans)
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

func (s *TransactionService) GetAllTransactions(ctx context.Context, userID primitive.ObjectID) (transactions []types.Transaction, err error) {
	transactions = []types.Transaction{}
	transaction := *s.collection
	cursor, err := transaction.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		s.service.logger.Error("Error while finding transactions", err)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &transactions)
	if err != nil {
		s.service.logger.Error("Error while decoding transactions", err)
		return
	}
	return transactions, nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, transID primitive.ObjectID, userID primitive.ObjectID) (transaction types.Transaction, err error) {
	transaction = types.Transaction{}
	transactionCollection := *s.collection
	err = transactionCollection.FindOne(ctx, bson.M{"_id": transID, "userID": userID}).Decode(&transaction)
	if err != nil {
		s.service.logger.Error("Error while decoding transaction", err)
		return
	}
	return transaction, nil
}

func (s *TransactionService) GetSoldTransactions(ctx context.Context, userID primitive.ObjectID) (transactions []types.Transaction, err error) {
	transactions = []types.Transaction{}
	transaction := *s.collection
	cursor, err := transaction.Find(ctx, bson.M{"userID": userID, "type": types.Sold})
	if err != nil {
		s.service.logger.Error("Error while finding sold transactions", err)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &transactions)
	if err != nil {
		s.service.logger.Error("Error while decoding sold transactions", err)
		return
	}
	return transactions, nil
}

func (s *TransactionService) GetBoughtTransactions(ctx context.Context, userID primitive.ObjectID) (transactions []types.Transaction, err error) {
	transactions = []types.Transaction{}
	transaction := *s.collection
	cursor, err := transaction.Find(ctx, bson.M{"userID": userID, "type": types.Bought})
	if err != nil {
		s.service.logger.Error("Error while finding bought transactions", err)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &transactions)
	if err != nil {
		s.service.logger.Error("Error while decoding bought transactions", err)
		return
	}
	return transactions, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, userID primitive.ObjectID, transID primitive.ObjectID) error {
	transaction := *s.collection
	_, err := transaction.DeleteOne(ctx, bson.M{"userID": userID, "_id": transID})
	if err != nil {
		s.service.logger.Error("Error deleteting provided transaction.", err)
	}
	return nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, userID primitive.ObjectID, transID primitive.ObjectID, payment bool) (types.Transaction, error) {
	transaction := *s.collection
	trans := types.Transaction{}
	err := transaction.FindOneAndUpdate(ctx, bson.M{"userID": userID, "_id": transID}, bson.M{
		"$set": bson.M{"payment": payment},
	}).Decode(&trans)
	if err != nil {
		s.service.logger.Error("Error updating transaction", err)
		return trans, err
	}
	trans, err = s.GetTransactionByID(ctx, transID, userID)
	if err != nil {
		s.service.logger.Error("Error getting updated transaction", err)
		return trans, err
	}
	return trans, nil
}

func (s *TransactionService) DeleteAllTransactions(ctx context.Context, userID primitive.ObjectID, transType types.TransactionType) error {
	transaction := *s.collection
	_, err := transaction.DeleteMany(ctx, bson.M{"userID": userID, "type": transType})
	if err != nil {
		s.service.logger.Error("Error deleting all transactions", err)
	}
	return nil
}
