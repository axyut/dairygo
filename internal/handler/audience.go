package handler

import (
	"net/http"

	"github.com/axyut/dairygo/internal/service"
)

type AudienceHandler struct {
	global *service.Service
	srv    *service.AudienceService
}

func (h *AudienceHandler) GetAudience(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user_id := r.URL.Query().Get("user_id")
	if user_id == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.GetAudience(user_id)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}
