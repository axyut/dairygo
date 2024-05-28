package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	service    *Service
	collection db.Collection
	Type       types.User
}

func (s *UserService) InsertUser(ctx context.Context, name string, email string, password string) (User types.User, error error) {
	user := *s.collection
	bytePass, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println(hashErr)
		return User, hashErr
	}
	newUser := types.User{
		UserName: name,
		Email:    email,
		Password: string(bytePass),
	}
	result, err := user.InsertOne(ctx, bson.M{
		"username": newUser.UserName,
		"email":    newUser.Email,
		"password": newUser.Password,
	})
	if err != nil {
		s.service.logger.Error("Error creating user", "Error", err)
		return User, err
	}
	user.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&newUser)
	s.service.logger.Info("User Created")
	return newUser, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) error {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		s.service.logger.Error("User not found")
		return err
	}
	// s.service.logger.Info("User Found")
	return nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (User types.User, error error) {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err != nil {
		s.service.logger.Error("User not found")
		return result, err
	}
	// s.service.logger.Info("User Found")
	return result, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (User types.User, error error) {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		s.service.logger.Error("User not found")
		return result, err
	}
	// s.service.logger.Info("User Found")
	return result, nil
}

func (s *UserService) LoginUser(ctx context.Context, email_username string, password string) (User types.User, error error) {
	user := *s.collection
	var result types.User
	err := user.FindOne(ctx, bson.M{"email": email_username}).Decode(&result)
	if err != nil {
		err := user.FindOne(ctx, bson.M{"username": email_username}).Decode(&result)
		if err != nil {
			s.service.logger.Error("User not found")
			return result, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		s.service.logger.Error("Password not matched")
		return result, err
	}
	return result, nil
}

func (s *UserService) UpdateUser(ctx context.Context, updateUser types.User) (User types.User, error error) {
	user := *s.collection
	res, err := user.UpdateOne(ctx, bson.M{"_id": updateUser.ID}, bson.M{"$set": bson.M{"username": updateUser.UserName, "email": updateUser.Email, "password": updateUser.Password}})
	if err != nil {
		s.service.logger.Error("Error updating user", "Error", err)
		return User, err
	}
	if res.ModifiedCount == 0 {
		s.service.logger.Error("User not found")
		return User, fmt.Errorf("user not found")
	}
	user.FindOne(ctx, bson.M{"_id": res.UpsertedID}).Decode(&updateUser)
	return updateUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	user := *s.collection
	res, err := user.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		s.service.logger.Error("Error deleting user", "Error", err)
		return err
	}
	if res.DeletedCount == 0 {
		s.service.logger.Error("User not found")
		return fmt.Errorf("user not found")
	}
	return nil
}
