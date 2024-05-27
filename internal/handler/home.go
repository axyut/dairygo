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
	user := h.h.UserHandler.GetUser(w, r)
	goods, err := h.h.srv.GoodsService.GetAllGoods(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	home := pages.Index(user, goods)
	client.Layout(home, "DairyGo").Render(r.Context(), w)
}

func (h *HomeHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	pages.NotFound().Render(r.Context(), w)
}
