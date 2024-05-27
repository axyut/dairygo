package handler

import (
	"fmt"
	"net/http"

	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AudienceHandler struct {
	h   *Handler
	srv *service.AudienceService
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
	err := h.srv.GetAudience(h.h.ctx, user_id)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}

func (h *AudienceHandler) NewAudience(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	contact := r.FormValue("contact")
	email := r.FormValue("email")
	userID := r.Context().Value("user_id")

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if name == "" || contact == "" || email == "" {
		components.AudienceInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}
	audience := types.Audience{
		Name:    name,
		Contact: contact,
		Email:   email,
		UserID:  id,
	}
	inserted, err := h.srv.NewAudience(h.h.ctx, audience)
	if err != nil {
		components.AudienceInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.AudienceInsertSuccess(inserted).Render(r.Context(), w)
}
