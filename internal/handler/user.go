package handler

import (
	"net/http"

	"github.com/axyut/dairygo/internal/service"
)

type UserHandler struct {
	global *service.Service
	srv    *service.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name, email, pass := r.URL.Query().Get("name"), r.URL.Query().Get("email"), r.URL.Query().Get("pass")
	if name == "" || email == "" || pass == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.NewUser(name, email, pass)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}
