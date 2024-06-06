package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/axyut/dairygo/internal/types"
)

type ReportsService struct {
	service *Service
}

func (s *ReportsService) GetReportPerDate(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate primitive.DateTime, toDate primitive.DateTime) (reports []types.TransactionReport, err error) {
	transColl := *s.service.Collections.Transactions
	var allTrans []types.Transaction

	startOfDate := primitive.NewDateTimeFromTime(time.Date(fromDate.Time().Year(), fromDate.Time().Month(), fromDate.Time().Day(), 0, 0, 0, 0, fromDate.Time().Location()))
	endOfDate := primitive.NewDateTimeFromTime(time.Date(toDate.Time().Year(), toDate.Time().Month(), toDate.Time().Day(), 23, 59, 59, 0, toDate.Time().Location()))

	goodsFilter := bson.E{Key: "goodID", Value: goodID}

	filter := bson.D{{Key: "userID", Value: userID}, goodsFilter, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: startOfDate}, {Key: "$lt", Value: endOfDate}}}}
	options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})

	cursor, errC := transColl.Find(ctx, filter, options)
	if errC != nil {
		s.service.logger.Error("Error while finding transactions", errC)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &allTrans)
	if err != nil {
		s.service.logger.Error("Error while decoding transactions", err)
		return
	}

	good, _ := s.service.GoodsService.GetGoodByID(ctx, userID, goodID)

	for ofDate := fromDate.Time(); ofDate.Before(toDate.Time().Add(time.Hour * 24)); ofDate = ofDate.AddDate(0, 0, 1) {
		var transofDate []types.Transaction
		var report types.TransactionReport
		// fmt.Println(ofDate.Day())
		for _, v := range allTrans {
			if v.CreationTime.Time().Year() == ofDate.Year() && v.CreationTime.Time().Month() == ofDate.Month() && v.CreationTime.Time().Day() == ofDate.Day() {
				transofDate = append(transofDate, v)
			}
		}
		report, err = s.getReportPerGood(transofDate)
		report.DayTime = primitive.NewDateTimeFromTime(ofDate)
		report.GoodName = good.Name
		report.GoodUnit = good.Unit
		// fmt.Println(report)
		reports = append(reports, report)
	}
	return
}

func (s *ReportsService) getReportPerGood(transactions []types.Transaction) (report types.TransactionReport, err error) {

	for _, v := range transactions {
		report.TotalTransactionQuantity += v.Quantity
		report.TotalTransactionPrice += v.Price
		if v.Type == types.Bought {
			report.TotalBoughtQuantity += v.Quantity
			report.TotalBoughtPrice += v.Price
		} else if v.Type == types.Sold {
			report.TotalSoldQuantity += v.Quantity
			report.TotalSoldPrice += v.Price
		}
		if v.Payment {
			report.TotalPaidPrice += v.Price
			report.TotalPaidQuantity += v.Quantity
		} else {
			report.TotalUnpaidPrice += v.Price
			report.TotalUnpaidQuantity += v.Quantity
		}
		report.TotalToPay += v.AudToPay
		report.TotalToReceive += v.AudToReceive
	}
	// total profit loss

	return
}

func (s *ReportsService) GetProductionReportPerDate(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate primitive.DateTime, toDate primitive.DateTime) (reports []types.ProductionReport, reportsPerGood []types.ProductionReportPerChangedGood, err error) {
	prodColl := *s.service.Collections.Production

	productions := []types.Production{}
	goodsFilter := bson.E{}
	var good types.Good
	if goodID != primitive.NilObjectID {
		goodsFilter = bson.E{Key: "changeGoodID", Value: goodID}
		good, _ = s.service.GoodsService.GetGoodByID(ctx, userID, goodID)

	}

	startOfDate := primitive.NewDateTimeFromTime(time.Date(fromDate.Time().Year(), fromDate.Time().Month(), fromDate.Time().Day(), 0, 0, 0, 0, fromDate.Time().Location()))
	endOfDate := primitive.NewDateTimeFromTime(time.Date(toDate.Time().Year(), toDate.Time().Month(), toDate.Time().Day(), 23, 59, 59, 0, toDate.Time().Location()))

	filter := bson.D{{Key: "userID", Value: userID}, goodsFilter, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: startOfDate}, {Key: "$lt", Value: endOfDate}}}}
	options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})

	cursor, errC := prodColl.Find(ctx, filter, options)
	if errC != nil {
		s.service.logger.Error("Error while finding sold productions", errC)
		return
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &productions)
	if err != nil {
		s.service.logger.Error("Error while decoding sold productions", err)
		return
	}
	for ofDate := fromDate.Time(); ofDate.Before(toDate.Time().Add(time.Hour * 24)); ofDate = ofDate.AddDate(0, 0, 1) {
		// fmt.Println(ofDate)
		var report types.ProductionReport
		var reportPerGood types.ProductionReportPerChangedGood

		var prodsofDate []types.Production

		for _, v := range productions {
			if v.CreationTime.Time().Year() == ofDate.Year() && v.CreationTime.Time().Month() == ofDate.Month() && v.CreationTime.Time().Day() == ofDate.Day() {
				prodsofDate = append(prodsofDate, v)
			}
		}

		if goodID == primitive.NilObjectID {
			report, err = s.getProductionReportAllGoods(prodsofDate)
			report.DayTime = primitive.NewDateTimeFromTime(ofDate)
			reports = append(reports, report)
		} else {

			reportPerGood, err = s.getProductionReportPerGood(prodsofDate)
			reportPerGood.ChangedGoodName = good.Name
			reportPerGood.ChangedGoodUnit = good.Unit
			reportPerGood.DayTime = primitive.NewDateTimeFromTime(ofDate)
			// fmt.Println(good.Name, good.Unit, reportPerGood.ChangedGoodName, reportPerGood.ChangedGoodUnit)
			reportsPerGood = append(reportsPerGood, reportPerGood)
		}

	}
	// fmt.Println(len(reports), len(reportsPerGood))
	return
}

func (s *ReportsService) getProductionReportPerGood(productions []types.Production) (report types.ProductionReportPerChangedGood, err error) {

	report = types.ProductionReportPerChangedGood{
		ProducedGoods: make(map[string]types.AGoodTotals),
	}
	var totalChangeQuantity float64
	var totalProducedQuantity float64
	var totalChangePrice float64
	var totalProducedPrice float64
	for _, v := range productions {
		// fmt.Println( v.ChangeGoodName, v.ChangeQuantity, v.ChangePrice, v.ProducedGoodName, v.ProducedQuantity, v.ProducedPrice)

		report.ProducedGoods[string(v.ProducedGoodID.Hex())] = types.AGoodTotals{
			Changed: types.Changed{
				Name:     v.ChangeGoodName,
				Unit:     v.ChangeGoodUnit,
				Quantity: v.ChangeQuantity + report.ProducedGoods[string(v.ProducedGoodID.Hex())].Changed.Quantity,
				Price:    v.ChangePrice + report.ProducedGoods[string(v.ProducedGoodID.Hex())].Changed.Price,
			},
			Produced: types.Produced{
				Name:     v.ProducedGoodName,
				Unit:     v.ProducedGoodUnit,
				Quantity: v.ProducedQuantity + report.ProducedGoods[string(v.ProducedGoodID.Hex())].Produced.Quantity,
				Price:    v.ProducedPrice + report.ProducedGoods[string(v.ProducedGoodID.Hex())].Produced.Price,
			},
		}

		totalChangeQuantity += v.ChangeQuantity
		totalProducedQuantity += v.ProducedQuantity
		totalChangePrice += v.ChangePrice
		totalProducedPrice += v.ProducedPrice
	}

	report.TotalChangedQuantity = totalChangeQuantity
	report.TotalProducedQuantity = totalProducedQuantity
	report.TotalChangePrice = totalChangePrice
	report.TotalProducedPrice = totalProducedPrice

	return
}

func (s *ReportsService) getProductionReportAllGoods(productions []types.Production) (report types.ProductionReport, err error) {

	report = types.ProductionReport{
		Goods: make(map[string]types.AGoodTotals),
	}
	for _, v := range productions {
		// fmt.Println(ofDate.Time().Day(), v.ChangeGoodName, v.ChangeQuantity, v.ChangePrice, v.ProducedGoodName, v.ProducedQuantity, v.ProducedPrice)

		report.Goods[string(v.ProducedGoodID.Hex())] = types.AGoodTotals{
			Changed: types.Changed{
				Name:     v.ChangeGoodName,
				Unit:     v.ChangeGoodUnit,
				Quantity: v.ChangeQuantity + report.Goods[string(v.ProducedGoodID.Hex())].Changed.Quantity,
				Price:    v.ChangePrice + report.Goods[string(v.ProducedGoodID.Hex())].Changed.Price,
			},
			Produced: types.Produced{
				Name:     v.ProducedGoodName,
				Unit:     v.ProducedGoodUnit,
				Quantity: v.ProducedQuantity + report.Goods[string(v.ProducedGoodID.Hex())].Produced.Quantity,
				Price:    v.ProducedPrice + report.Goods[string(v.ProducedGoodID.Hex())].Produced.Price,
			},
		}

	}

	// fmt.Println(report)
	return
}

// database heavy operation, next function is better for small scale

// func (s *ReportsService) GetProductionReportPerDate(ctx context.Context, userID primitive.ObjectID, goodID primitive.ObjectID, fromDate primitive.DateTime, toDate primitive.DateTime) (reports []types.ProductionReport, reportsPerGood []types.ProductionReportPerChangedGood, err error) {
// 	prodColl := *s.service.Collections.Production

// 	for ofDate := fromDate.Time(); ofDate.Before(toDate.Time()); ofDate = ofDate.AddDate(0, 0, 1) {

// 		productions := []types.Production{}
// 		goodsFilter := bson.E{}
// 		if goodID != primitive.NilObjectID {
// 			goodsFilter = bson.E{Key: "changeGoodID", Value: goodID}
// 		}

// 		startOfDay := primitive.NewDateTimeFromTime(time.Date(ofDate.Year(), ofDate.Month(), ofDate.Day(), 0, 0, 0, 0, ofDate.Location()))
// 		endOfDay := primitive.NewDateTimeFromTime(time.Date(ofDate.Year(), ofDate.Month(), ofDate.Day(), 23, 59, 59, 0, ofDate.Location()))

// 		filter := bson.D{{Key: "userID", Value: userID}, goodsFilter, {Key: "creationTime", Value: bson.D{{Key: "$gt", Value: startOfDay}, {Key: "$lt", Value: endOfDay}}}}
// 		options := options.Find().SetSort(bson.D{{Key: "creationTime", Value: -1}})

// 		cursor, errC := prodColl.Find(ctx, filter, options)
// 		if errC != nil {
// 			s.service.logger.Error("Error while finding sold productions", errC)
// 			return
// 		}
// 		defer cursor.Close(ctx)
// 		err = cursor.All(ctx, &productions)
// 		if err != nil {
// 			s.service.logger.Error("Error while decoding sold productions", err)
// 			return
// 		}

// 		var report types.ProductionReport
// 		var reportPerGood types.ProductionReportPerChangedGood
// 		if goodID == primitive.NilObjectID {
// 			report, err = s.GetProductionReportAllGoods(startOfDay, productions)
// 			reports = append(reports, report)
// 		} else {
// 			reportPerGood, err = s.GetProductionReportPerGood(startOfDay, productions)
// 			reportsPerGood = append(reportsPerGood, reportPerGood)
// 		}
// 	}
// 	fmt.Println(len(reports), len(reportsPerGood))
// 	return
// }
