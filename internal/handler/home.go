package handler

import (
	"fmt"
	"net/http"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/pages"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeHandler struct {
	h *Handler
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := h.h.UserHandler.GetUser(w, r)

	goods, err := h.h.srv.GoodsService.GetAllGoods(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	auds, err := h.h.srv.AudienceService.GetAllAudiences(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	home := pages.Index(user, goods, auds)
	client.Layout(home, "DairyGo").Render(r.Context(), w)
}

func (h *HomeHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	pages.NotFound().Render(r.Context(), w)
}
