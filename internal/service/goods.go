package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *GoodsService) GetGoods(ctx context.Context, name string) error {
	goods := *s.collection
	res, err := goods.Find(ctx, bson.M{"name": name})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
