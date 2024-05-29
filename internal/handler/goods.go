package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoodsHandler struct {
	h   *Handler
	srv *service.GoodsService
}

func (h *GoodsHandler) NewGood(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	kharidRate, errK := strconv.ParseFloat(r.FormValue("kharid_rate"), 64)
	bikriRate, errP := strconv.ParseFloat(r.FormValue("bikri_rate"), 64)
	unit := r.FormValue("unit")
	user_id := r.Context().Value("user_id")

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if errK != nil || errP != nil {
		components.GeneralToastError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}
	if name == "" {
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	good := types.Good{
		Name:       name,
		KharidRate: kharidRate,
		BikriRate:  bikriRate,
		Unit:       unit,
		UserID:     userID,
	}

	insertedGood, err := h.srv.InsertGood(h.h.ctx, good)
	if err != nil {
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	goods, _ := h.srv.GetAllGoods(r.Context(), userID)
	components.GoodInsertSuccess(insertedGood, goods).Render(r.Context(), w)
}

func (h *GoodsHandler) DeleteGood(w http.ResponseWriter, r *http.Request) {
	good_id := r.URL.Query().Get("id")
	user_id := r.Context().Value("user_id")

	goodID, _ := primitive.ObjectIDFromHex(good_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))

	err := h.srv.DeleteGood(h.h.ctx, userID, goodID)
	if err != nil {
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.GeneralToastSuccess("Deleted Successfully").Render(r.Context(), w)
}

func (h *GoodsHandler) UpdateGood(w http.ResponseWriter, r *http.Request) {
	good_id := r.URL.Query().Get("id")
	name := r.FormValue("td_name" + good_id)
	kharid_rate := strings.TrimSpace(r.FormValue("td_kharid_rate" + good_id))
	bikri_rate := strings.TrimSpace(r.FormValue("td_bikri_rate" + good_id))
	user_id := r.Context().Value("user_id")
	fmt.Println(kharid_rate, bikri_rate)

	// h.h.logger.Info("UpdateGood", "good_id", good_id, "name", name, "rate", rate, "user_id", user_id)

	goodID, _ := primitive.ObjectIDFromHex(good_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	g, _ := h.srv.GetGoodByID(h.h.ctx, userID, goodID)
	fmt.Println(g.Name, g.Unit)
	if strings.Contains(kharid_rate, " /"+g.Unit) {
		kharid_rate = strings.ReplaceAll(kharid_rate, " /"+g.Unit, "")
	}
	if strings.Contains(bikri_rate, " /"+g.Unit) {
		bikri_rate = strings.ReplaceAll(bikri_rate, " /"+g.Unit, "")
	}

	fmt.Println(kharid_rate, bikri_rate)

	kharidRate, errK := strconv.ParseFloat(kharid_rate, 64)
	bikriRate, errB := strconv.ParseFloat(bikri_rate, 64)
	if errB != nil || errK != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}

	if name == "" || bikriRate == 0 || kharidRate == 0 {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}

	good := types.UpdateGood{
		Name:       name,
		BikriRate:  bikriRate,
		KharidRate: kharidRate,
		Unit:       g.Unit,
		Quantity:   g.Quantity,
	}

	insertedGood, err := h.srv.UpdateGood(h.h.ctx, userID, goodID, good)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GoodInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	goods, _ := h.srv.GetAllGoods(r.Context(), userID)
	w.WriteHeader(http.StatusOK)
	components.GoodInsertSuccess(insertedGood, goods).Render(r.Context(), w)
}
