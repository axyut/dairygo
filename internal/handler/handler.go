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
}

func NewHandler(ctx context.Context, srv *service.Service, logger *slog.Logger) *Handler {
	h := &Handler{ctx: ctx, srv: srv, logger: logger}

	h.AudienceHandler = &AudienceHandler{h, srv.AudienceService}
	h.GoodsHandler = &GoodsHandler{h, srv.GoodsService}
	h.TransactionHandler = &TransactionHandler{h, srv.TransactionService}
	h.UserHandler = &UserHandler{h, srv.UserService}
	h.HomeHandler = &HomeHandler{h}

	return h
}

func RootHandler(ctx context.Context, conf config.Config, srv *service.Service, logger *slog.Logger) *http.Server {

	h := NewHandler(ctx, srv, logger)
	router := http.NewServeMux()
	chained := m.Chain(m.Logging, m.Auth)

	router.HandleFunc("/", h.HomeHandler.GetHome)
	router.HandleFunc("/register", h.UserHandler.NewUser)
	router.HandleFunc("/login", h.UserHandler.LoginUser)
	router.HandleFunc("GET /sold", h.TransactionHandler.GetSold)
	router.HandleFunc("GET /bought", h.TransactionHandler.GetBought)
	router.HandleFunc("GET /logout", h.UserHandler.LogoutUser)
	router.HandleFunc("GET /getUserReq", h.UserHandler.GetUserRequest)
	router.HandleFunc("GET /profile", h.UserHandler.GetProfile)

	router.HandleFunc("POST /audience", h.AudienceHandler.NewAudience)
	router.HandleFunc("POST /goods", h.GoodsHandler.NewGood)
	router.HandleFunc("POST /transaction", h.TransactionHandler.NewTransaction)
	router.HandleFunc("POST /internalTransaction", h.TransactionHandler.InternalTransaction)

	router.HandleFunc("DELETE /audience", h.AudienceHandler.DeleteAudience)
	router.HandleFunc("DELETE /goods", h.GoodsHandler.DeleteGood)

	router.HandleFunc("PATCH /audience", h.AudienceHandler.UpdateAudience)
	router.HandleFunc("PATCH /goods", h.GoodsHandler.UpdateGood)
	// calculate
	// unavailable all other routes

	server := &http.Server{
		Addr:    ":" + conf.PORT,
		Handler: chained(router),
	}
	return server
}
