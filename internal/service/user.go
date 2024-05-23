package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	service    *Service
	collection db.Collection
	Type       types.User
}

func (s *UserService) NewUser(ctx context.Context, name string, email string, password string) error {
	user := *s.collection
	bytePass, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println(hashErr)
		return hashErr
	}
	newUser := types.User{
		UserName: name,
		Email:    email,
		Password: string(bytePass),
	}
	_, err := user.InsertOne(ctx, bson.M{
		"username": newUser.UserName,
		"email":    newUser.Email,
		"password": newUser.Password,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.service.logger.Info("User Created")
	return nil
}
