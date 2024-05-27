package service

import (
	"context"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AudienceService struct {
	service    *Service
	collection db.Collection
}

func (s *AudienceService) NewAudience(ctx context.Context, aud types.Audience) (insertedAud types.Audience, err error) {
	audience := *s.collection
	insertedAud = types.Audience{}
	res, err := audience.InsertOne(ctx, bson.M{
		"name":      aud.Name,
		"contact":   aud.Contact,
		"email":     aud.Email,
		"userID":    aud.UserID,
		"toPay":     aud.ToPay,
		"toReceive": aud.ToReceive,
		"paid":      aud.Paid,
	})
	if err != nil {
		s.service.logger.Error("Error while inserting new audience", err)
		return
	}
	audience.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&insertedAud)
	return
}

func (s *AudienceService) GetAudienceByID(ctx context.Context, audID primitive.ObjectID) (aud types.Audience, err error) {
	audience := *s.collection
	aud = types.Audience{}
	err = audience.FindOne(ctx, bson.M{"_id": audID}).Decode(&aud)
	if err != nil {
		s.service.logger.Error("Error while fetching audience", err)
		return
	}
	return
}

func (s *AudienceService) GetAllAudiences(ctx context.Context, userID primitive.ObjectID) (aud []types.Audience, err error) {
	audience := *s.collection
	aud = []types.Audience{}

	cursor, err := audience.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		s.service.logger.Error("Error while fetching all audiences", err)
		return
	}
	defer cursor.Close(ctx)
	if cursor.All(ctx, &aud) != nil {
		s.service.logger.Error("Error while decoding all audiences", err)
		return
	}
	return
}
