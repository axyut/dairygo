package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/db"
	"go.mongodb.org/mongo-driver/bson"
)

type UserService struct {
	global     *db.Collections
	collection db.Collection
	ctx        context.Context
}

func (s *UserService) NewUser(name string, email string, password string) error {
	user := *s.collection
	res, err := user.InsertOne(s.ctx, bson.M{
		"name":     name,
		"email":    email,
		"password": password,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
