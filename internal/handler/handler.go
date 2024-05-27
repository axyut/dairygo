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
	router.HandleFunc("/register", h.UserHandler.CreateUser)
	router.HandleFunc("/login", h.UserHandler.LoginUser)
	router.HandleFunc("GET /sold", h.TransactionHandler.GetSold)
	router.HandleFunc("GET /bought", h.TransactionHandler.GetBought)
	router.HandleFunc("POST /logout", h.UserHandler.LogoutUser)
	router.HandleFunc("POST /audience", h.AudienceHandler.NewAudience)
	router.HandleFunc("POST /goods", h.GoodsHandler.NewGood)
	router.HandleFunc("POST /transaction", h.TransactionHandler.NewTransaction)
	// calculate
	// unavailable all other routes

	server := &http.Server{
		Addr:    ":" + conf.PORT,
		Handler: chained(router),
	}
	return server
}
