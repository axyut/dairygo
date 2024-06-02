package service

import (
	"context"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductionService struct {
	service    *Service
	collection db.Collection
}

func (s *ProductionService) InsertProduction(ctx context.Context, production types.Production, userID primitive.ObjectID) (insertedProduction types.Production, err error) {
	prod := *s.collection
	insertedProduction = types.Production{}
	res, err := prod.InsertOne(ctx, bson.M{
		"changeGoodID":     production.ChangeGoodID,
		"changeQuantity":   production.ChangeQuantity,
		"changeGoodName":   production.ChangeGoodName,
		"changeGoodUnit":   production.ChangeGoodUnit,
		"producedGoodID":   production.ProducedGoodID,
		"producedQuantity": production.ProducedQuantity,
		"producedGoodName": production.ProducedGoodName,
		"producedGoodUnit": production.ProducedGoodUnit,
		"profit":           production.Profit,
		"loss":             production.Loss,
		"userID":           userID,
	})
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
	cursor, err := prod.Find(ctx, bson.M{"userID": userID})
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