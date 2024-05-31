package service

import (
	"context"
	"log/slog"

	"github.com/axyut/dairygo/internal/db"
)

type Service struct {
	mongo              *db.Mongo
	Collections        *db.Collections
	logger             *slog.Logger
	UserService        *UserService
	AudienceService    *AudienceService
	GoodsService       *GoodsService
	TransactionService *TransactionService
	ProductionService  *ProductionService
}

func NewService(ctx context.Context, mongo *db.Mongo, logger *slog.Logger) *Service {
	global := db.GetCollections(ctx, mongo.DB)

	service := &Service{
		mongo:       mongo,
		Collections: global,
		logger:      logger,
	}

	service.UserService = &UserService{service: service, collection: db.GetUserColl(ctx, mongo.DB)}
	service.AudienceService = &AudienceService{service: service, collection: db.GetAudienceColl(ctx, mongo.DB)}
	service.GoodsService = &GoodsService{service: service, collection: db.GetGoodsColl(ctx, mongo.DB)}
	service.TransactionService = &TransactionService{service: service, collection: db.GetTransactionColl(ctx, mongo.DB)}
	service.ProductionService = &ProductionService{service: service, collection: db.GetProductionColl(ctx, mongo.DB)}

	return service
}
