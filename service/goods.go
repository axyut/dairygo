package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/db"
	"go.mongodb.org/mongo-driver/bson"
)

type GoodsService struct {
	global     *db.Collections
	collection db.Collection
	ctx        context.Context
}

func (s *GoodsService) NewGoods(name string, price float64, quantity int) error {
	goods := *s.collection
	res, err := goods.InsertOne(s.ctx, bson.M{
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

func (s *GoodsService) GetGoods(name string) error {
	goods := *s.collection
	res, err := goods.Find(s.ctx, bson.M{"name": name})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
