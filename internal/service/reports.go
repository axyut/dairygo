package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/axyut/dairygo/internal/types"
)

type ReportsService struct {
	service *Service
}

func (s *ReportsService) GetReportDate(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate time.Time, toDate time.Time) (report types.Report, err error) {
	return
}

func (s *ReportsService) GetReportPerDay(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, ofDate time.Time) (report types.Report, err error) {
	return
}

func (s *ReportsService) GetReportPerGood(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, ofDate time.Time) (report types.ProductionReport, err error) {

	return
}
func (s *ReportsService) GetProductionReportPerDate(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate time.Time, toDate time.Time) (report types.ProductionReport, err error) {

	return
}
func (s *ReportsService) GetProductionReportPerDay(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, ofDate time.Time) (report types.ProductionReport, err error) {

	return
}
func (s *ReportsService) GetProductionReportPerGood(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate time.Time, toDate time.Time) (report types.ProductionReport, err error) {

	return
}
