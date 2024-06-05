package handler

import (
	"net/http"
	"time"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportsHandler struct {
	h   *Handler
	srv *service.ReportsService
}

func (h *ReportsHandler) GetReportsPage(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	w.WriteHeader(http.StatusAccepted)
	p := pages.ReportsPage(goods)
	client.Layout(p, "Reports").Render(r.Context(), w)
}

func (h *ReportsHandler) GenerateReports(w http.ResponseWriter, r *http.Request) {
	from_date := r.FormValue("from_date")
	to_date := r.FormValue("to_date")
	good_id := r.FormValue("goodID")
	reportType := r.FormValue("type")

	goodID, errGood := primitive.ObjectIDFromHex(good_id)
	fromDate := utils.GetMongoTimeFromHTMLDate(from_date, time.Now().Local().AddDate(0, -1, 0))
	toDate := utils.GetMongoTimeFromHTMLDate(to_date, time.Now().Local())

	if fromDate > toDate {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("From date can't be bigger than to date.").Render(r.Context(), w)
		return
	}

	if good_id == "" || good_id == "all" || errGood != nil {
		goodID = primitive.NilObjectID
	}

	if reportType == "production" {
		h.GenerateProductionReports(w, r, goodID, fromDate, toDate)
	} else {
		h.GenerateTransactionReports(w, r, goodID, fromDate, toDate)
	}
}

func (h *ReportsHandler) GenerateProductionReports(w http.ResponseWriter, r *http.Request, goodID primitive.ObjectID, fromDate primitive.DateTime, toDate primitive.DateTime) {
	_, _, err := h.srv.GetProductionReportPerDate(h.h.ctx, h.h.UserHandler.GetUser(w, r).ID, goodID, fromDate, toDate)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	pages.ProdReport().Render(r.Context(), w)
}

func (h *ReportsHandler) GenerateTransactionReports(w http.ResponseWriter, r *http.Request, goodID primitive.ObjectID, fromDate primitive.DateTime, toDate primitive.DateTime) {
	h.h.UserHandler.GetUser(w, r)
	w.WriteHeader(http.StatusAccepted)
	components.GeneralToastError("Under Construction").Render(r.Context(), w)
}

func (h *ReportsHandler) DownloadReports(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	h.GenerateReports(w, r)
	// TODO: download report
}
