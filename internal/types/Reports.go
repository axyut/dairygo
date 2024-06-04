package types

import "time"

type Report struct {
	FromTime               time.Time `json:"fromDate"`
	ToTime                 time.Time `json:"toDate"`
	GoodName               string    `json:"goodName"`
	TotalRemainingQuantity float64   `json:"totalRemainingQuantity"`
	TotalProfit            float64   `json:"totalProfit"`
	TotalLoss              float64   `json:"totalLoss"`
	TotalPaidQuantity      float64   `json:"totalPaidQuantity"`
	TotalUnpaidQuantity    float64   `json:"totalUnpaidQuantity"`
	TotalQuantity          float64   `json:"totalQuantity"`
	TotalRemainingPrice    float64   `json:"totalRemainingPrice"`
	TotalBoughtPrice       float64   `json:"totalBoughtPrice"`
	TotalSoldPrice         float64   `json:"totalSoldPrice"`
	TotalToPay             float64   `json:"totalToPay"`
	TotalToReceive         float64   `json:"totalToReceive"`
}

type ProductionReport struct {
	FromTime         time.Time `json:"fromDate"`
	ToTime           time.Time `json:"toDate"`
	ChangedGoodName  string    `json:"changedGoodName"`
	ChangedQuantity  float64   `json:"changedQuantity"`
	ChangedUnit      string    `json:"changedUnit"`
	ProducedGoodName string    `json:"producedGoodName"`
	ProducedQuantity float64   `json:"producedQuantity"`
	ProducedUnit     string    `json:"producedUnit"`
	Profit           float64   `json:"profit"`
	Loss             float64   `json:"loss"`
	TotalQuantity    float64   `json:"totalQuantity"`
	TotalPrice       float64   `json:"totalPrice"`
	TotalProfit      float64   `json:"totalProfit"`
	TotalLoss        float64   `json:"totalLoss"`
}
