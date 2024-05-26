package handler

import (
	"context"
	"log/slog"
	"net/http"

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

func RootHandler(ctx context.Context, srv *service.Service, logger *slog.Logger) *Handler {

	h := NewHandler(ctx, srv, logger)

	// Templ Handler
	http.HandleFunc("/", h.HomeHandler.GetHome)

	// User
	http.HandleFunc("/register", h.UserHandler.CreateUser)
	http.HandleFunc("/login", h.UserHandler.LoginUser)

	// Audience
	http.HandleFunc("/api/audience", h.AudienceHandler.GetAudience)

	// Goods
	http.HandleFunc("/api/goods", h.GoodsHandler.GetGoods)

	// Transaction
	http.HandleFunc("/api/transaction", h.TransactionHandler.GetTransaction)

	// calculate

	return h
}
