package handler

import (
	"net/http"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/pages"
)

type HomeHandler struct {
	h *Handler
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	home := pages.GuestIndex()
	client.Layout(home, "DairyGo").Render(h.h.ctx, w)
}
