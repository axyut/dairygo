package service

import (
	"context"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoodsService struct {
	service    *Service
	collection db.Collection
}

func (s *GoodsService) GetGoodByID(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID) (good types.Good, err error) {
	goods := *s.collection
	good = types.Good{}
	err = goods.FindOne(ctx, bson.M{"_id": goodID, "userID": userID}).Decode(&good)
	if err != nil {
		s.service.logger.Error("Error while fetching good", err)
		return
	}
	return
}

func (s *GoodsService) InsertGood(ctx context.Context, good types.Good) (insertedGood types.Good, err error) {
	goods := *s.collection
	res, err := goods.InsertOne(ctx, good)
	if err != nil {
		s.service.logger.Error("Error inserting good", "Error", err)
		return good, err
	}
	// fmt.Println(res)
	err = goods.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&good)
	if err != nil {
		s.service.logger.Error("Error decoding good", "Error", err)
		return good, err
	}
	return good, nil
}

func (s *GoodsService) GetAllGoods(ctx context.Context, userID primitive.ObjectID) (goods []types.Good, err error) {
	goods = []types.Good{}
	goodsCollection := *s.collection

	cursor, err := goodsCollection.Find(ctx, bson.M{
		"userID": userID,
	})
	if err != nil {
		s.service.logger.Error("Error while fetching all goods", err)
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var good types.Good
		err = cursor.Decode(&good)
		if err != nil {
			s.service.logger.Error("Error while decoding all goods", err)
			return
		}
		goods = append(goods, good)
	}
	return
}

func (s *GoodsService) UpdateGood(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, update types.UpdateGood) (good types.Good, err error) {
	goods := *s.collection
	good = types.Good{}
	err = goods.FindOneAndUpdate(ctx, bson.M{"_id": goodID, "userID": userID}, bson.M{"$set": update}).Decode(&good)
	if err != nil {
		s.service.logger.Error("Error while updating good", err)
		return
	}
	good, _ = s.GetGoodByID(ctx, userID, goodID)
	return
}

func (s *GoodsService) DeleteGood(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID) (err error) {
	goods := *s.collection
	_, err = goods.DeleteOne(ctx, bson.M{"_id": goodID, "userID": userID})
	if err != nil {
		s.service.logger.Error("Error while deleting good", err)
		return
	}
	return
}

func (s *GoodsService) DeleteAllGoods(ctx context.Context, userID primitive.ObjectID) (err error) {
	goods := *s.collection
	_, err = goods.DeleteMany(ctx, bson.M{"userID": userID})
	if err != nil {
		s.service.logger.Error("Error while deleting all goods", err)
		return
	}
	return
}
