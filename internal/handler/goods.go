package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	sellingRate, errP := strconv.ParseFloat(r.FormValue("bikri_rate"), 64)
	unit := r.FormValue("unit")
	user := h.h.UserHandler.GetUser(w, r)

	if errP != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	good := types.Good{
		ID:          primitive.NewObjectIDFromTimestamp(time.Now()),
		Name:        name,
		SellingRate: sellingRate,
		Unit:        unit,
		UserID:      user.ID,
	}

	insertedGood, err := h.srv.InsertGood(h.h.ctx, good)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	goods, _ := h.srv.GetAllGoods(r.Context(), user.ID)
	components.GoodInsertSuccess(insertedGood, goods).Render(r.Context(), w)
}

func (h *GoodsHandler) DeleteGood(w http.ResponseWriter, r *http.Request) {
	good_id := r.URL.Query().Get("id")
	user_id := r.Context().Value("user_id")

	goodID, _ := primitive.ObjectIDFromHex(good_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))

	err := h.srv.DeleteGood(h.h.ctx, userID, goodID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.GeneralToastSuccess("Deleted Successfully").Render(r.Context(), w)
}

func (h *GoodsHandler) UpdateGood(w http.ResponseWriter, r *http.Request) {
	good_id := r.URL.Query().Get("id")
	name := r.FormValue("good_name_" + good_id)
	selling_rate := strings.TrimSpace(r.FormValue("good_selling_rate_" + good_id))
	user_id := r.Context().Value("user_id")

	// h.h.logger.Info("UpdateGood", "good_id", good_id, "name", name, "rate", rate, "user_id", user_id)

	goodID, _ := primitive.ObjectIDFromHex(good_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	g, _ := h.srv.GetGoodByID(h.h.ctx, userID, goodID)

	if strings.Contains(selling_rate, " /"+g.Unit) {
		selling_rate = strings.ReplaceAll(selling_rate, " /"+g.Unit, "")
	}

	sellingRate, errB := strconv.ParseFloat(selling_rate, 64)
	if errB != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}

	if name == "" || sellingRate == 0 {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}

	good := types.UpdateGood{
		Name:        name,
		SellingRate: sellingRate,
		Unit:        g.Unit,
		Quantity:    g.Quantity,
	}

	_, err := h.srv.UpdateGood(h.h.ctx, userID, goodID, good)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GoodInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	goods, _ := h.srv.GetAllGoods(r.Context(), userID)
	w.WriteHeader(http.StatusOK)
	components.GoodsTable(goods, true).Render(r.Context(), w)
}

func (h *GoodsHandler) RefreshGoods(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)

	allGoods, err := h.srv.GetAllGoods(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("Please Re-login.").Render(r.Context(), w)
		return
	}

	components.GoodsTable(allGoods, true).Render(r.Context(), w)
}

func (h *GoodsHandler) DeleteAllGoods(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	err := h.srv.DeleteAllGoods(h.h.ctx, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	components.GeneralToastSuccess("Deleted Successfully").Render(r.Context(), w)
}
