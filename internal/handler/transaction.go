package handler

import (
	"net/http"

	"github.com/axyut/dairygo/internal/service"
)

type TransactionHandler struct {
	global *service.Service
	srv    *service.TransactionService
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	transID := r.URL.Query().Get("id")
	if transID == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.GetTransaction(transID)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}
