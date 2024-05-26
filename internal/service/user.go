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

func (s *UserService) GetUser(ctx context.Context, email string) error {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.service.logger.Info("User Found")
	return nil
}

func (s *UserService) LoginUser(ctx context.Context, email_username string, password string) error {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"email": email_username}).Decode(&result)
	if err != nil {
		err := user.FindOne(ctx, bson.M{"username": email_username}).Decode(&result)
		if err != nil {
			s.service.logger.Error("User not found")
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		s.service.logger.Error("Password not matched")
		return err
	}
	return nil
}
