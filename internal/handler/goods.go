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
	rate, perr := strconv.ParseFloat(r.FormValue("rate"), 64)
	unit := r.FormValue("unit")
	user_id := r.Context().Value("user_id")

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if perr != nil {
		components.GeneralToastError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}
	if name == "" {
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	good := types.Good{
		Name:   name,
		Rate:   rate,
		Unit:   unit,
		UserID: userID,
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
	name := r.FormValue("good_name")
	good_rate := strings.TrimSpace(r.FormValue("good_rate"))
	user_id := r.Context().Value("user_id")

	// h.h.logger.Info("UpdateGood", "good_id", good_id, "name", name, "rate", rate, "user_id", user_id)
	goodID, _ := primitive.ObjectIDFromHex(good_id)
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	g, _ := h.srv.GetGoodByID(h.h.ctx, userID, goodID)

	if strings.Contains(good_rate, " /"+g.Unit) {
		good_rate = strings.ReplaceAll(good_rate, " /"+g.Unit, "")
	}

	rate, err := strconv.ParseFloat(good_rate, 64)
	fmt.Println(rate, good_rate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}
	if name == "" || rate == 0 {
		w.WriteHeader(http.StatusBadRequest)
		components.GoodInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}

	good := types.UpdateGood{
		Name: name,
		Rate: rate,
		Unit: g.Unit,
		// Price:    g.Price,
		Quantity: g.Quantity,
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
