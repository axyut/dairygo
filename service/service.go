package service

import (
	"context"

	"github.com/axyut/dairygo/db"
)

type Service struct {
	mongo              *db.Mongo
	Collections        *db.Collections
	Ctx                context.Context
	UserService        *UserService
	AudienceService    *AudienceService
	GoodsService       *GoodsService
	TransactionService *TransactionService
}

func NewService(mongo *db.Mongo) *Service {
	global := db.GetCollections(db.Ctx, mongo.DB)
	userSrv := &UserService{global, db.GetUserColl(db.Ctx, mongo.DB), db.Ctx}
	audSrv := &AudienceService{global, db.GetAudienceColl(db.Ctx, mongo.DB), db.Ctx}
	goodsSrv := &GoodsService{global, db.GetGoodsColl(db.Ctx, mongo.DB), db.Ctx}
	transSrv := &TransactionService{global, db.GetTransactionColl(db.Ctx, mongo.DB), db.Ctx}
	return &Service{mongo, global, db.Ctx, userSrv, audSrv, goodsSrv, transSrv}
}
