package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	DayTime                primitive.DateTime `json:"dayTime"`
	GoodName               string             `json:"goodName"`
	TotalRemainingQuantity float64            `json:"totalRemainingQuantity"`
	TotalProfit            float64            `json:"totalProfit"`
	TotalLoss              float64            `json:"totalLoss"`
	TotalPaidQuantity      float64            `json:"totalPaidQuantity"`
	TotalUnpaidQuantity    float64            `json:"totalUnpaidQuantity"`
	TotalQuantity          float64            `json:"totalQuantity"`
	TotalRemainingPrice    float64            `json:"totalRemainingPrice"`
	TotalBoughtPrice       float64            `json:"totalBoughtPrice"`
	TotalSoldPrice         float64            `json:"totalSoldPrice"`
	TotalToPay             float64            `json:"totalToPay"`
	TotalToReceive         float64            `json:"totalToReceive"`
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
