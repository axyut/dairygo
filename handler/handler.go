package handler

import (
	"net/http"

	"github.com/a-h/templ"
	client "github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/service"
)

type Handler struct {
	srv                *service.Service
	AudienceHandler    *AudienceHandler
	GoodsHandler       *GoodsHandler
	TransactionHandler *TransactionHandler
	UserHandler        *UserHandler
}

func NewHandler(srv *service.Service) *Handler {
	audienceHandler := &AudienceHandler{srv, srv.AudienceService}
	goodsHandler := &GoodsHandler{srv, srv.GoodsService}
	transactionHandler := &TransactionHandler{srv, srv.TransactionService}
	userHandler := &UserHandler{srv, srv.UserService}
	return &Handler{srv, audienceHandler, goodsHandler, transactionHandler, userHandler}
}

func RootHandler(srv *service.Service) {
	h := NewHandler(srv)

	// Templ Handler
	homePage := client.Index()
	http.Handle("/", templ.Handler(homePage))

	// API Handler
	// User
	http.HandleFunc("/api/user", h.UserHandler.CreateUser)

	// Audience
	http.HandleFunc("/api/audience", h.AudienceHandler.GetAudience)

	// Goods
	http.HandleFunc("/api/goods", h.GoodsHandler.GetGoods)

	// Transaction
	http.HandleFunc("/api/transaction", h.TransactionHandler.GetTransaction)
}
