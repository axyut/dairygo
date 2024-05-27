package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionHandler struct {
	h   *Handler
	srv *service.TransactionService
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) (trans []types.Transaction) {
	userID := r.Context().Value("user_id")
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		components.GeneralToastError("Couldn't Fullfill Your Request. Please Try Again.").Render(r.Context(), w)
		h.h.logger.Error("Error while Parsing in Handler.", err)
		return
	}

	trans, err = h.srv.GetTransaction(r.Context(), id)
	if err != nil {
		components.GeneralToastError("Error with service.").Render(r.Context(), w)
		return
	}
	return
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

	if goodID == "" || quantity == "" || price == "" || audienceID == "" || trans_type == "" || userID == "" {
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

func (h *TransactionHandler) GetSold(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		components.GeneralToastError("Couldn't Identify User. Please Login Again.")
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		components.GeneralToastError("Couldn't Fullfill Your Request. Please Try Again.")
	}

	soldTrans, err := h.h.srv.TransactionService.GetSoldTransactions(r.Context(), id)
	if err != nil {
		components.GeneralToastError("Error with service.")
		return
	}
	client_Trans := []types.Transaction_Client{}

	for _, v := range soldTrans {
		good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), v.GoodID)
		soldAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.SoldTo)
		// boughtAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.BoughtFrom)
		client_Trans = append(client_Trans, types.Transaction_Client{
			ID:       v.ID,
			GoodName: good.Name,
			GoodUnit: good.Unit,
			Quantity: strconv.FormatFloat(v.Quantity, 'f', 2, 64),
			Price:    strconv.FormatFloat(v.Price, 'f', 2, 64),
			SoldTo:   soldAudience.Name,
			Payment:  v.Payment,
			// BoughtFrom: boughtAudience.Name,
		})
	}
	soldPage := pages.Sold(client_Trans)
	client.Layout(soldPage, "Sold Dairy Products").Render(r.Context(), w)
}

func (h *TransactionHandler) GetBought(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		components.GeneralToastError("Couldn't Identify User. Please Login Again.")
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		components.GeneralToastError("Couldn't Fullfill Your Request. Please Try Again.")
	}
	boughts, err := h.h.srv.TransactionService.GetBoughtTransactions(r.Context(), id)
	if err != nil {
		components.GeneralToastError("Error with service.")
		return
	}
	client_T := []types.Transaction_Client{}
	for _, v := range boughts {
		good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), v.GoodID)
		// soldAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.SoldTo)
		boughtAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.BoughtFrom)
		client_T = append(client_T, types.Transaction_Client{
			ID:       v.ID,
			GoodName: good.Name,
			GoodUnit: good.Unit,
			Quantity: strconv.FormatFloat(v.Quantity, 'f', 2, 64),
			Price:    strconv.FormatFloat(v.Price, 'f', 2, 64),
			// SoldTo:     soldAudience.Name,
			Payment:    v.Payment,
			BoughtFrom: boughtAudience.Name,
		})
	}
	boughtPage := pages.Bought(client_T)
	client.Layout(boughtPage, "Bought Dairy Products").Render(r.Context(), w)
}
