package service

import (
	"context"
	"math"

	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AudienceService struct {
	service    *Service
	collection db.Collection
}

func (s *AudienceService) InsertAudience(ctx context.Context, aud types.Audience) (insertedAud types.Audience, err error) {
	audience := *s.collection
	insertedAud = types.Audience{}
	res, err := audience.InsertOne(ctx, bson.M{
		"name":      aud.Name,
		"contact":   aud.Contact,
		"userID":    aud.UserID,
		"toPay":     math.Abs(aud.ToPay),
		"toReceive": math.Abs(aud.ToReceive),
		"mapRates":  make(map[string]float64),
	})
	if err != nil {
		s.service.logger.Error("Error while inserting new audience", err)
		return
	}
	audience.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&insertedAud)
	return
}

func (s *AudienceService) GetAudienceByID(ctx context.Context, userID primitive.ObjectID, audID primitive.ObjectID) (aud types.Audience, err error) {
	audience := *s.collection
	aud = types.Audience{}
	err = audience.FindOne(ctx, bson.M{"_id": audID, "userID": userID}).Decode(&aud)
	if err != nil {
		s.service.logger.Error("Error while fetching audience", err)
		return
	}
	return
}

func (s *AudienceService) GetAllAudiences(ctx context.Context, userID primitive.ObjectID) (auds []types.Audience, err error) {
	audience := *s.collection
	auds = []types.Audience{}

	cursor, err := audience.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		s.service.logger.Error("Error while fetching all audiences", err)
		return
	}
	defer cursor.Close(ctx)
	if cursor.All(ctx, &auds) != nil {
		s.service.logger.Error("Error while decoding all audiences", err)
		return
	}

	return
}

func (s *AudienceService) UpdateAudience(ctx context.Context, update types.Audience) (aud types.Audience, err error) {
	audience := *s.collection
	aud = types.Audience{}
	res, err := audience.UpdateOne(ctx, bson.M{"_id": update.ID}, bson.M{"$set": bson.M{
		"name":      update.Name,
		"contact":   update.Contact,
		"toPay":     math.Abs(update.ToPay),
		"toReceive": math.Abs(update.ToReceive),
		"mapRates":  update.MapRates,
		"userID":    update.UserID,
	}})
	if err != nil {
		s.service.logger.Error("Error while updating audience", err)
		return
	}
	audience.FindOne(ctx, bson.M{"_id": res.UpsertedID}).Decode(&aud)
	return
}

func (s *AudienceService) DeleteAudience(ctx context.Context, userID primitive.ObjectID, audID primitive.ObjectID) (err error) {
	audience := *s.collection
	_, err = audience.DeleteOne(ctx, bson.M{"_id": audID, "userID": userID})
	if err != nil {
		s.service.logger.Error("Error while deleting audience", err)
		return
	}
	return
}

func (s *AudienceService) DeleteAllAudiences(ctx context.Context, userID primitive.ObjectID) (err error) {
	audience := *s.collection
	_, err = audience.DeleteMany(ctx, bson.M{"userID": userID})
	if err != nil {
		s.service.logger.Error("Error while deleting all audiences", err)
		return
	}
	return
}

// func (s *AudienceService) GetBuyingRate(ctx context.Context, userID primitive.ObjectID, audienceID primitive.ObjectID, goodID primitive.ObjectID) (rate types.Rate, err error) {
// 	aud := *s.collection
// 	audience := types.Audience{}
// 	err = aud.FindOne(ctx, bson.M{"_id": audienceID, "userID": userID}).Decode(&audience)
// 	if err != nil {
// 		s.service.logger.Error("Error while fetching audience", err)
// 		return
// 	}
// 	for _, rate := range audience.Rates {
// 		if rate.GoodID == goodID {
// 			return rate, nil
// 		}
// 	}
// 	return
// }
