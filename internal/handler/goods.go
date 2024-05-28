package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	userID := r.Context().Value("user_id")

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if perr != nil {
		components.GoodInsertError("Enter Numeric Value for Rate.").Render(r.Context(), w)
		return
	}
	if name == "" {
		components.GoodInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}
	good := types.Good{
		Name:   name,
		Rate:   rate,
		Unit:   unit,
		UserID: id,
	}

	insertedGood, err := h.srv.InsertGood(h.h.ctx, good)
	if err != nil {
		components.GoodInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.GoodInsertSuccess(insertedGood).Render(r.Context(), w)
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
