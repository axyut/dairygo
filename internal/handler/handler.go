package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/axyut/dairygo/internal/config"
	m "github.com/axyut/dairygo/internal/middleware"
	"github.com/axyut/dairygo/internal/service"
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
	h.HomeHandler = &HomeHandler{h}
	h.ReportsHandler = &ReportsHandler{h}
	h.ProductionHandler = &ProductionHandler{h, srv.ProductionService}

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
	// router.HandleFunc("GET /audience/buying_rate", h.AudienceHandler.BuyingRate)
	router.HandleFunc("GET /goods/refresh", h.GoodsHandler.RefreshGoods)

	router.HandleFunc("POST /audience", h.AudienceHandler.NewAudience)
	router.HandleFunc("POST /goods", h.GoodsHandler.NewGood)
	router.HandleFunc("POST /transaction", h.TransactionHandler.NewTransaction)
	router.HandleFunc("POST /production", h.ProductionHandler.NewProduction)

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
	// calculate
	// unavailable all other routes

	server := &http.Server{
		Addr:    ":" + conf.PORT,
		Handler: middleware(router),
	}
	return server
}
