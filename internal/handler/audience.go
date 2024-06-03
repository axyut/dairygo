package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

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
	good_id := r.URL.Query().Get("good_id")
	name := r.FormValue("aud_name_" + audience_id)
	contact := r.FormValue("aud_contact_" + audience_id)
	buying_rate := strings.TrimSpace(r.FormValue("aud_buying_rate"))
	user := h.h.UserHandler.GetUser(w, r)

	audID, err := primitive.ObjectIDFromHex(audience_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		components.GeneralToastError("Couldn't parse Audience ID.").Render(r.Context(), w)
		return
	}
	aud, _ := h.srv.GetAudienceByID(h.h.ctx, user.ID, audID)

	var goodID primitive.ObjectID
	var buyingRate float64
	if good_id != "" && buying_rate != "" {
		goodID, err = primitive.ObjectIDFromHex(good_id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			components.GeneralToastError("Couldn't parse Good ID.").Render(r.Context(), w)
			return
		}
		buyingRate, err = strconv.ParseFloat(buying_rate, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			components.GeneralToastError("Couldn't parse Buying Rate.").Render(r.Context(), w)
			return
		}
		aud.MapRates[goodID.Hex()] = math.Abs(buyingRate)

	} else if name != "" && contact != "" {
		aud.Name = name
		aud.Contact = contact
	}

	_, err = h.srv.UpdateAudience(h.h.ctx, aud)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	if good_id == "" && buying_rate == "" {
		allAuds, _ := h.srv.GetAllAudiences(r.Context(), user.ID)
		w.WriteHeader(http.StatusOK)
		components.AudTable(allAuds, true, goods).Render(r.Context(), w)
		return
	} else if good_id != "" && buying_rate != "" {
		w.WriteHeader(http.StatusOK)
		components.BuyingRateforAudRow(aud, goods, "Rate has been set to "+strconv.FormatFloat(aud.MapRates[goodID.Hex()], 'f', -1, 64)).Render(r.Context(), w)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		components.GeneralToastSuccess("Updated Successfully").Render(r.Context(), w)
		return
	}
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

// func (h *AudienceHandler) BuyingRate(w http.ResponseWriter, r *http.Request) {
// 	user := h.h.UserHandler.GetUser(w, r)
// 	audience_id := r.URL.Query().Get("aud_id")
// 	good_id := r.URL.Query().Get("good_id")

// 	audID, errA := primitive.ObjectIDFromHex(audience_id)
// 	goodID, errG := primitive.ObjectIDFromHex(good_id)

// 	if errA != nil || errG != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		components.GeneralToastError("Couldn't find your requested data.").Render(r.Context(), w)
// 		return
// 	}
// 	rate, err := h.srv.GetBuyingRate(h.h.ctx, user.ID, audID, goodID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusExpectationFailed)
// 		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
// 		return
// 	}
// 	components.BuyingRate(strconv.FormatFloat(rate.BuyingRate, 'f', -1, 64)).Render(r.Context(), w)

// }
