package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionReport struct {
	DayTime                  primitive.DateTime `json:"dayTime"`
	GoodName                 string             `json:"goodName"`
	GoodUnit                 string             `json:"goodUnit"`
	TotalPaidQuantity        float64            `json:"totalPaidQuantity"`
	TotalPaidPrice           float64            `json:"totalPaidPrice"`
	TotalUnpaidQuantity      float64            `json:"totalUnpaidQuantity"`
	TotalUnpaidPrice         float64            `json:"totalUnpaidPrice"`
	TotalTransactionQuantity float64            `json:"totalQuantity"`
	TotalTransactionPrice    float64            `json:"totalPrice"`
	// if Total Good Quantity of that day is needed, make changes in transaction table
	// TotalGoodQuantity             float64            `json:"totalGoodQuantity"`
	// TotalGoodPrice                float64            `json:"totalGoodPrice"`
	// TotalRemainingGoodQuantity float64            `json:"totalRemainingGoodQuantity"`
	// TotalRemainingGoodPrice    float64            `json:"totalRemainingGoodPrice"`
	TotalBoughtQuantity float64 `json:"totalBoughtQuantity"`
	TotalBoughtPrice    float64 `json:"totalBoughtPrice"`
	TotalSoldQuantity   float64 `json:"totalSoldQuantity"`
	TotalSoldPrice      float64 `json:"totalSoldPrice"`
	TotalToPay          float64 `json:"totalToPay"`
	TotalToReceive      float64 `json:"totalToReceive"`
	TotalProfit         float64 `json:"totalProfit"`
	TotalLoss           float64 `json:"totalLoss"`
}

type ProductionReport struct {
	DayTime     primitive.DateTime     `json:"dayTime"`
	Goods       map[string]AGoodTotals `json:"producedGoods"`
	TotalProfit float64                `json:"totalProfit"`
	TotalLoss   float64                `json:"totalLoss"`
}

type ProductionReportPerChangedGood struct {
	DayTime               primitive.DateTime     `json:"dayTime"`
	TotalChangedQuantity  float64                `json:"changedQuantity"`
	TotalChangePrice      float64                `json:"totalChangePrice"`
	ChangedGoodName       string                 `json:"changedGoodName"`
	ChangedGoodUnit       string                 `json:"changedGoodUnit"`
	TotalProducedQuantity float64                `json:"totalQuantity"`
	TotalProducedPrice    float64                `json:"totalProducedPrice"`
	TotalProfit           float64                `json:"totalProfit"`
	TotalLoss             float64                `json:"totalLoss"`
	ProducedGoods         map[string]AGoodTotals `json:"producedGoods"`
}

type AGoodTotals struct {
	Changed  Changed  `json:"changed" bson:"changed"`
	Produced Produced `json:"produced" bson:"produced"`
}

type Changed struct {
	Name     string  `json:"name" bson:"name"`
	Unit     string  `json:"unit" bson:"unit"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
}

type Produced struct {
	Name     string  `json:"name" bson:"name"`
	Unit     string  `json:"unit" bson:"unit"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
}
