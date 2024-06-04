package handler

import (
	"net/http"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
)

type ReportsHandler struct {
	h *Handler
}

func (h *ReportsHandler) GetReportsPage(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	w.WriteHeader(http.StatusAccepted)
	p := pages.ReportsPage(goods)
	client.Layout(p, "Reports").Render(r.Context(), w)
}

func (h *ReportsHandler) GenerateReports(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	w.WriteHeader(http.StatusAccepted)
	components.GeneralToastError("Under Construction").Render(r.Context(), w)
}

func (h *ReportsHandler) GenerateProductionReports(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	w.WriteHeader(http.StatusAccepted)
	components.GeneralToastError("Under Construction").Render(r.Context(), w)
}

func (h *ReportsHandler) GenerateTransactionReports(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	w.WriteHeader(http.StatusAccepted)
	components.GeneralToastError("Under Construction").Render(r.Context(), w)
}

func (h *ReportsHandler) DownloadReports(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	h.GenerateReports(w, r)
	// TODO: download report
}
