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

func (h *AudienceHandler) NewAudience(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	contact := r.FormValue("contact")
	userID := r.Context().Value("user_id")

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.AudienceInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	if name == "" || contact == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.AudienceInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}
	audience := types.Audience{
		Name:    name,
		Contact: contact,
		UserID:  id,
	}
	inserted, err := h.srv.InsertAudience(h.h.ctx, audience)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.AudienceInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.AudienceInsertSuccess(inserted).Render(r.Context(), w)
}

func (h *AudienceHandler) DeleteAudience(w http.ResponseWriter, r *http.Request) {
	audience_id := r.URL.Query().Get("id")
	user_id := r.Context().Value("user_id")

	audID, _ := primitive.ObjectIDFromHex(audience_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))

	err := h.srv.DeleteAudience(h.h.ctx, userID, audID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.GeneralToastSuccess("Audience Deleted.").Render(r.Context(), w)
}

func (h *AudienceHandler) UpdateAudience(w http.ResponseWriter, r *http.Request) {
	audience_id := r.URL.Query().Get("id")
	name := r.FormValue("aud_name")
	contact := r.FormValue("aud_contact")
	user_id := r.Context().Value("user_id")

	audID, _ := primitive.ObjectIDFromHex(audience_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))

	if name == "" || contact == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	aud, _ := h.srv.GetAudienceByID(h.h.ctx, userID, audID)
	aud.Name = name
	aud.Contact = contact

	_, err := h.srv.UpdateAudience(h.h.ctx, aud)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.AudienceInsertSuccess(aud).Render(r.Context(), w)
}

func (h *AudienceHandler) RefreshAudience(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)

	allAuds, err := h.srv.GetAllAudiences(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("Please Re-login.").Render(r.Context(), w)
		return
	}

	components.AudTable(allAuds, true).Render(r.Context(), w)
}
