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

type ProductionService struct {
	service    *Service
	collection db.Collection
}

func (s *ProductionService) InsertProduction(ctx context.Context, production types.Production, userID primitive.ObjectID) (insertedProduction types.Production, err error) {
	prod := *s.collection
	insertedProduction = types.Production{}
	production.Profit = math.Abs(production.Profit)
	production.Loss = math.Abs(production.Loss)
	production.ProducedQuantity = math.Abs(production.ProducedQuantity)
	production.ChangeQuantity = math.Abs(production.ChangeQuantity)

	res, err := prod.InsertOne(ctx, production)
	if err != nil {
		s.service.logger.Error("Error while inserting new production", err)
		return
	}
	err = prod.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&insertedProduction)
	if err != nil {
		s.service.logger.Error("Error while decoding new production", err)
		return
	}
	return
}

func (s *ProductionService) GetAllProductions(ctx context.Context, userID primitive.ObjectID) (productions []types.Production, err error) {
	productions = []types.Production{}
	prod := *s.collection

	lastweek := primitive.NewDateTimeFromTime(time.Now().Local().AddDate(0, 0, -7))
	options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: lastweek}}}}

	cursor, err := prod.Find(ctx, filter, options)
	if err != nil {
		s.service.logger.Error("Error while finding productions", err)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &productions)
	if err != nil {
		s.service.logger.Error("Error while decoding productions", err)
		return
	}
	return
}

func (s *ProductionService) GetProductionByID(ctx context.Context, prodID primitive.ObjectID, userID primitive.ObjectID) (production types.Production, err error) {
	production = types.Production{}
	prod := *s.collection
	err = prod.FindOne(ctx, bson.M{"_id": prodID, "userID": userID}).Decode(&production)
	return
}

func (s *ProductionService) DeleteProduction(ctx context.Context, prodID primitive.ObjectID, userID primitive.ObjectID) (err error) {
	prod := *s.collection
	_, err = prod.DeleteOne(ctx, bson.M{"_id": prodID, "userID": userID})
	return
}

func (s *ProductionService) DeleteAllProductions(ctx context.Context, userID primitive.ObjectID) (err error) {
	prod := *s.collection
	_, err = prod.DeleteMany(ctx, bson.M{"userID": userID})
	return
}
