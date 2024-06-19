package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/internal/config"
	m "github.com/axyut/dairygo/internal/middleware"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
)

type Handler struct {
	ctx                context.Context
	srv                *service.Service
	logger             *slog.Logger
	AudienceHandler    *AudienceHandler
	GoodsHandler       *GoodsHandler
	TransactionHandler *TransactionHandler
	UserHandler        *UserHandler
	HomeHandler        *HomeHandler
	ReportsHandler     *ReportsHandler
	ProductionHandler  *ProductionHandler
}

func NewHandler(ctx context.Context, srv *service.Service, logger *slog.Logger) *Handler {
	h := &Handler{ctx: ctx, srv: srv, logger: logger}

	h.AudienceHandler = &AudienceHandler{h, srv.AudienceService}
	h.GoodsHandler = &GoodsHandler{h, srv.GoodsService}
	h.TransactionHandler = &TransactionHandler{h, srv.TransactionService}
	h.UserHandler = &UserHandler{h, srv.UserService}
	h.ReportsHandler = &ReportsHandler{h, srv.ReportsService}
	h.ProductionHandler = &ProductionHandler{h, srv.ProductionService}
	h.HomeHandler = &HomeHandler{h}

	return h
}

func RootHandler(ctx context.Context, conf config.Config, srv *service.Service, logger *slog.Logger) *http.Server {

	h := NewHandler(ctx, srv, logger)
	router := http.NewServeMux()
	middleware := m.Chain(m.Logging, m.Auth)

	router.HandleFunc("/", h.HomeHandler.GetHome)
	router.HandleFunc("/register", h.UserHandler.NewUser)
	router.HandleFunc("/login", h.UserHandler.LoginUser)
	router.HandleFunc("GET /transaction", h.TransactionHandler.GetTransactionPage)
	router.HandleFunc("GET /sold", h.TransactionHandler.GetSold)
	router.HandleFunc("GET /bought", h.TransactionHandler.GetBought)
	router.HandleFunc("GET /production", h.ProductionHandler.GetProductionPage)
	router.HandleFunc("GET /logout", h.UserHandler.LogoutUser)
	router.HandleFunc("GET /getUserReq", h.UserHandler.GetUserRequest)
	router.HandleFunc("GET /profile", h.UserHandler.GetProfile)
	router.HandleFunc("GET /reports", h.ReportsHandler.GetReportsPage)
	router.HandleFunc("GET /audience/refresh", h.AudienceHandler.RefreshAudience)
	router.HandleFunc("GET /goods/refresh", h.GoodsHandler.RefreshGoods)

	router.HandleFunc("POST /audience", h.AudienceHandler.NewAudience)
	router.HandleFunc("POST /goods", h.GoodsHandler.NewGood)
	router.HandleFunc("POST /transaction", h.TransactionHandler.NewTransaction)
	router.HandleFunc("POST /transaction/refresh", h.RefreshTransaction)
	router.HandleFunc("POST /production", h.ProductionHandler.NewProduction)
	router.HandleFunc("POST /reports/refresh", h.ReportsHandler.GenerateReports)
	router.HandleFunc("POST /reports/download", h.ReportsHandler.DownloadReports)

	router.HandleFunc("DELETE /audience", h.AudienceHandler.DeleteAudience)
	router.HandleFunc("DELETE /goods", h.GoodsHandler.DeleteGood)
	router.HandleFunc("DELETE /transaction", h.TransactionHandler.DeleteTransaction)
	router.HandleFunc("DELETE /production", h.ProductionHandler.DeleteProduction)
	router.HandleFunc("DELETE /user", h.UserHandler.DeleteUser)
	router.HandleFunc("DELETE /audience/all", h.AudienceHandler.DeleteAllAudiences)
	router.HandleFunc("DELETE /goods/all", h.GoodsHandler.DeleteAllGoods)
	router.HandleFunc("DELETE /transaction/all", h.TransactionHandler.DeleteAllTransactions)
	router.HandleFunc("DELETE /production/all", h.ProductionHandler.DeleteAllProductions)

	router.HandleFunc("PATCH /audience", h.AudienceHandler.UpdateAudience)
	router.HandleFunc("PATCH /goods", h.GoodsHandler.UpdateGood)
	router.HandleFunc("PATCH /transaction", h.TransactionHandler.UpdateTransaction)
	router.HandleFunc("PATCH /user/defaults", h.UserHandler.UpdateUserDefaults)
	// calculate
	// unavailable all other routes

	server := &http.Server{
		Addr:    ":" + conf.PORT,
		Handler: middleware(router),
	}
	return server
}

func (h *Handler) RefreshTransaction(w http.ResponseWriter, r *http.Request) {
	trans_type := r.FormValue("type")
	date := r.FormValue("date")
	payment := r.FormValue("payment")
	audID := r.FormValue("aud_id_filter")
	goodID := r.FormValue("good_id_filter")

	if trans_type == "" || date == "" || payment == "" || audID == "" || goodID == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	h.UserHandler.GetUser(w, r)
	if date != "today" && date != "yesterday" && date != "lastweek" && date != "alltime" && date != "thismonth" && date != "lastmonth" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Wrong Date. Empty Fields!").Render(r.Context(), w)
		return
	}

	if payment != "paid" && payment != "unpaid" && payment != "all" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Wrong Payment. Empty Fields!").Render(r.Context(), w)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), types.CtxAudID, audID))
	r = r.WithContext(context.WithValue(r.Context(), types.CtxGoodID, goodID))

	r = r.WithContext(context.WithValue(r.Context(), types.CtxDate, date))
	r = r.WithContext(context.WithValue(r.Context(), types.CtxPayment, payment))

	if trans_type == "production" {
		h.ProductionHandler.GetProductionPage(w, r)
	} else if trans_type == "sold" {
		h.TransactionHandler.GetSold(w, r)
	} else if trans_type == "bought" {
		h.TransactionHandler.GetBought(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Wrong Transaction Type. Empty Fields!").Render(r.Context(), w)
		return
	}
}
