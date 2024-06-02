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
	user := h.h.UserHandler.GetUser(w, r)

	if name == "" || contact == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.AudienceInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}
	audience := types.Audience{
		Name:    name,
		Contact: contact,
		UserID:  user.ID,
	}
	inserted, err := h.srv.InsertAudience(h.h.ctx, audience)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.AudienceInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	components.AudienceInsertSuccess(inserted, goods).Render(r.Context(), w)
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
	user := h.h.UserHandler.GetUser(w, r)

	audID, _ := primitive.ObjectIDFromHex(audience_id)

	if name == "" || contact == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	aud, _ := h.srv.GetAudienceByID(h.h.ctx, user.ID, audID)
	aud.Name = name
	aud.Contact = contact

	_, err := h.srv.UpdateAudience(h.h.ctx, aud)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	components.AudienceInsertSuccess(aud, goods).Render(r.Context(), w)
}

func (h *AudienceHandler) RefreshAudience(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)

	allAuds, err := h.srv.GetAllAudiences(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("Please Re-login.").Render(r.Context(), w)
		return
	}
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	components.AudTable(allAuds, true, goods).Render(r.Context(), w)
}

func (h *AudienceHandler) DeleteAllAudiences(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	err := h.srv.DeleteAllAudiences(h.h.ctx, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	components.GeneralToastSuccess("Deleted Successfully").Render(r.Context(), w)
}
