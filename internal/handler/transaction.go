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

type TransactionHandler struct {
	h   *Handler
	srv *service.TransactionService
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	transID := r.URL.Query().Get("id")
	if transID == "" {
		http.Error(w, "Empty Fields!", http.StatusBadRequest)
	}
	err := h.srv.GetTransaction(h.h.ctx, transID)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	}
}

func (h *TransactionHandler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	goodID := r.FormValue("goodID")
	quantity := r.FormValue("quantity")
	price := r.FormValue("price")
	audienceID := r.FormValue("audienceID")
	trans_type := r.FormValue("type")
	payment := r.FormValue("payment")
	userID := r.Context().Value("user_id")

	var boughtFrom primitive.ObjectID
	var soldTo primitive.ObjectID
	var payment_b bool = false

	if goodID == "" || quantity == "" || price == "" || audienceID == "" || trans_type == "" || payment == "" || userID == "" {
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	// fmt.Println(goodID, quantity, price, audienceID, trans_type, payment, userID)

	id, Uerr := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	good_id, Gerr := primitive.ObjectIDFromHex(goodID)
	aud_id, Aerr := primitive.ObjectIDFromHex(audienceID)
	quantity_f, ferr := strconv.ParseFloat(quantity, 64)
	price_f, perr := strconv.ParseFloat(price, 64)

	// fmt.Println(id, good_id, aud_id, quantity_f, price_f, payment_b)
	if Uerr != nil || Gerr != nil || Aerr != nil || ferr != nil || perr != nil {
		h.h.logger.Error("Error while Parsing in Handler.", Uerr, Gerr, Aerr, ferr, perr)
		components.GeneralToastError("Error while Parsing in Handler.").Render(r.Context(), w)
		return
	}

	if trans_type != string(types.Bought) && trans_type != string(types.Sold) {
		components.GeneralToastError("Invalid Transaction Type.").Render(r.Context(), w)
		return
	}

	if trans_type == string(types.Bought) {
		boughtFrom = aud_id
	} else if trans_type == string(types.Sold) {
		soldTo = aud_id
	}

	if payment == "on" {
		payment_b = true
	}

	transaction := types.Transaction{
		GoodID:     good_id,
		Quantity:   quantity_f,
		Price:      price_f,
		BoughtFrom: boughtFrom,
		SoldTo:     soldTo,
		Type:       types.TransactionType(trans_type),
		Payment:    payment_b,
		UserID:     id,
	}
	_, err := h.srv.NewTransaction(h.h.ctx, transaction)
	if err != nil {
		components.GeneralToastError("Error with service.").Render(r.Context(), w)
		return
	}
	components.GeneralToastSuccess("Transaction Added Successfully!").Render(r.Context(), w)
}
