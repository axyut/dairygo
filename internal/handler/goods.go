package handler

import (
	"net/http"

	"github.com/axyut/dairygo/internal/service"
)

type GoodsHandler struct {
	global *service.Service
	srv    *service.GoodsService
}

func (h *GoodsHandler) GetGoods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	goodID := r.URL.Query().Get("id")
	if goodID == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.GetGoods(goodID)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}
