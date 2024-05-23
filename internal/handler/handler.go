package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	client "github.com/axyut/dairygo/client"
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
}

func NewHandler(ctx context.Context, srv *service.Service, logger *slog.Logger) *Handler {
	h := &Handler{ctx: ctx, srv: srv, logger: logger}

	h.AudienceHandler = &AudienceHandler{h, srv.AudienceService}
	h.GoodsHandler = &GoodsHandler{h, srv.GoodsService}
	h.TransactionHandler = &TransactionHandler{h, srv.TransactionService}
	h.UserHandler = &UserHandler{h, srv.UserService}

	return h
}

func RootHandler(ctx context.Context, srv *service.Service, logger *slog.Logger) *Handler {

	h := NewHandler(ctx, srv, logger)

	// Templ Handler
	homePage := client.Index()
	http.Handle("/", templ.Handler(homePage))

	// User
	http.HandleFunc("/api/user", h.UserHandler.CreateUser)

	// Audience
	http.HandleFunc("/api/audience", h.AudienceHandler.GetAudience)

	// Goods
	http.HandleFunc("/api/goods", h.GoodsHandler.GetGoods)

	// Transaction
	http.HandleFunc("/api/transaction", h.TransactionHandler.GetTransaction)

	// calculate

	return h
}
