package service

import (
	"context"
	"math"
	"time"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	date, _ := ctx.Value(types.CtxDate).(string)
	pay, _ := ctx.Value(types.CtxPayment).(string)
	aud, _ := ctx.Value(types.CtxAudID).(string)
	g, _ := ctx.Value(types.CtxGoodID).(string)
	start, end := GetDateTime(date)
	payment := GetPayment(pay)
	audience := s.GetAudience(aud, userID, types.Sold)
	good := s.GetGood(g, userID)

	options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "type", Value: types.Sold}, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: start}, {Key: "$lt", Value: end}}}, payment, audience, good}

	cursor, err := transaction.Find(ctx, filter, options)
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
	date, _ := ctx.Value(types.CtxDate).(string)
	pay, _ := ctx.Value(types.CtxPayment).(string)
	aud, _ := ctx.Value(types.CtxAudID).(string)
	g, _ := ctx.Value(types.CtxGoodID).(string)
	start, end := GetDateTime(date)
	payment := GetPayment(pay)
	audience := s.GetAudience(aud, userID, types.Bought)
	good := s.GetGood(g, userID)

	options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "type", Value: types.Bought}, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: start}, {Key: "$lt", Value: end}}}, payment, audience, good}

	cursor, err := transaction.Find(ctx, filter, options)
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

func (s *TransactionService) UpdateTransaction(ctx context.Context, userID primitive.ObjectID, transID primitive.ObjectID, trans types.Transaction) (types.Transaction, error) {
	transaction := *s.collection
	err := transaction.FindOneAndUpdate(ctx, bson.M{"userID": userID, "_id": transID}, bson.M{
		"$set": trans,
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

func GetDateTime(date string) (start primitive.DateTime, end primitive.DateTime) {
	end = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Now().Location()))
	if date == "lastweek" { // last 7 days
		start = primitive.NewDateTimeFromTime(time.Now().Local().AddDate(0, 0, -7))
	} else if date == "alltime" {
		start = primitive.NewDateTimeFromTime(time.UnixMicro(0))
	} else if date == "yesterday" {
		start = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Now().Location()))
		end = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 23, 59, 59, 0, time.Now().Location()))
	} else if date == "thismonth" {
		start = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()))
	} else if date == "lastmonth" {
		start = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location()))
		end = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month()-1, 30, 23, 59, 59, 0, time.Now().Location()))
	} else { // today
		start = primitive.NewDateTimeFromTime(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()))
	}
	return
}

func GetPayment(payment string) (pay primitive.E) {
	pay = primitive.E{}
	if payment == "paid" {
		pay = primitive.E{Key: "payment", Value: true}
	} else if payment == "unpaid" {
		pay = primitive.E{Key: "payment", Value: false}
	}
	return
}

func (s *TransactionService) GetAudience(audience string, userID primitive.ObjectID, typee types.TransactionType) (aud primitive.E) {
	aud = primitive.E{}

	audID, err := primitive.ObjectIDFromHex(audience)
	if err != nil {
		// s.service.logger.Error("Error while getting audience in refresh transaction", err)
		return
	}
	audExists, err := s.service.AudienceService.GetAudienceByID(context.Background(), userID, audID)
	if err != nil {
		// s.service.logger.Error("Error while getting audience in refresh transaction", err)
		return
	}
	if typee == types.Bought {
		aud = primitive.E{Key: "boughtFromID", Value: audExists.ID}
	} else if typee == types.Sold {
		aud = primitive.E{Key: "soldToID", Value: audExists.ID}
	}
	return

}

func (s *TransactionService) GetGood(good string, userID primitive.ObjectID) (g primitive.E) {
	g = primitive.E{}
	goodID, err := primitive.ObjectIDFromHex(good)
	if err != nil {
		// s.service.logger.Error("Error while getting good in refresh transaction", err)
		return
	}
	goodExists, err := s.service.GoodsService.GetGoodByID(context.Background(), userID, goodID)
	if err != nil {
		// s.service.logger.Error("Error while getting good in refresh transaction", err)
		return
	}
	g = primitive.E{Key: "goodID", Value: goodExists.ID}
	return
}
