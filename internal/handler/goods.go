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

func (h *GoodsHandler) GetGoods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	goodID := r.URL.Query().Get("id")
	if goodID == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.GetGoods(h.h.ctx, goodID)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}

func (h *GoodsHandler) NewGood(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	price, perr := strconv.ParseFloat(r.FormValue("price"), 64)
	quantity, qerr := strconv.ParseFloat(r.FormValue("quantity"), 64)
	userID := r.Context().Value("user_id")

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}

	if perr != nil || qerr != nil {
		components.GoodInsertError("Enter Numeric Value for Price and Quantity.").Render(r.Context(), w)
		return
	}
	if name == "" {
		components.GoodInsertError("Empty Fields!").Render(r.Context(), w)
		return
	}
	good := types.Good{
		Name:     name,
		Price:    price,
		Quantity: quantity,
		UserID:   id,
	}

	insertedGood, err := h.srv.InsertGood(h.h.ctx, good)
	if err != nil {
		components.GoodInsertError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	components.GoodInsertSuccess(insertedGood).Render(r.Context(), w)
}
