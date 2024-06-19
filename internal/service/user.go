package service

import (
	"context"
	"fmt"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
	user.FindOne(ctx, bson.M{"_id": result.InsertedID}, opts).Decode(&newUser)
	s.service.logger.Info("User Created")
	return newUser, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) error {
	user := *s.collection
	var result types.User
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
	err := user.FindOne(ctx, bson.M{"email": email}, opts).Decode(&result)
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
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
	err := user.FindOne(ctx, bson.M{"username": username}, opts).Decode(&result)
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
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
	err := user.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&result)
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
	result.Password = ""
	return result, nil
}

func (s *UserService) UpdateUser(ctx context.Context, updateUser types.User) (User types.User, error error) {
	user := *s.collection
	res, err := user.UpdateOne(ctx, bson.M{"_id": updateUser.ID}, bson.M{"$set": updateUser})
	if err != nil {
		s.service.logger.Error("Error updating user", "Error", err)
		return User, err
	}
	opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}})
	user.FindOne(ctx, bson.M{"_id": res.UpsertedID}, opts).Decode(&updateUser)
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
