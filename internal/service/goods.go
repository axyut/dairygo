package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoodsService struct {
	service    *Service
	collection db.Collection
}

func (s *GoodsService) NewGoods(ctx context.Context, name string, price float64, quantity int) error {
	goods := *s.collection
	res, err := goods.InsertOne(ctx, bson.M{
		"name":     name,
		"price":    price,
		"quantity": quantity,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func (s *GoodsService) GetGoodByID(ctx context.Context, id primitive.ObjectID) (good types.Good, err error) {
	goods := *s.collection
	good = types.Good{}
	err = goods.FindOne(ctx, bson.M{"_id": id}).Decode(&good)
	if err != nil {
		s.service.logger.Error("Error while fetching good", err)
		return
	}
	return
}

func (s *GoodsService) InsertGood(ctx context.Context, good types.Good) (insertedGood types.Good, err error) {
	goods := *s.collection
	res, err := goods.InsertOne(ctx, bson.M{
		"name":     good.Name,
		"price":    good.Price,
		"quantity": good.Quantity,
		"unit":     good.Unit,
		"userID":   good.UserID,
	})
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
