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

func (h *TransactionHandler) GetTransactionPage(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	page := pages.TransactionPage()
	client.Layout(page, "Transactions").Render(r.Context(), w)
}

func (h *TransactionHandler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	goodID := r.FormValue("goodID")
	quantity := r.FormValue("quantity")
	audienceID := r.FormValue("audienceID")
	trans_type := r.FormValue("type")
	payment := r.FormValue("payment")
	user_id := r.Context().Value("user_id")

	var boughtFrom primitive.ObjectID
	var soldTo primitive.ObjectID
	var payment_b bool = false

	if goodID == "" || quantity == "" || audienceID == "" || trans_type == "" || user_id == "" {
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	// fmt.Println(goodID, quantity, price, audienceID, trans_type, payment, userID)

	userID, Uerr := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	good_id, Gerr := primitive.ObjectIDFromHex(goodID)
	aud_id, Aerr := primitive.ObjectIDFromHex(audienceID)
	quantity_f, ferr := strconv.ParseFloat(quantity, 64)

	// fmt.Println(id, good_id, aud_id, quantity_f, price_f, payment_b)
	if Uerr != nil || Gerr != nil || Aerr != nil || ferr != nil {
		h.h.logger.Error("Error while Parsing in Handler.", Uerr, Gerr, Aerr, ferr)
		components.GeneralToastError("Error while Parsing in Handler.").Render(r.Context(), w)
		return
	}

	if trans_type != string(types.Bought) && trans_type != string(types.Sold) {
		components.GeneralToastError("Invalid Transaction Type.").Render(r.Context(), w)
		return
	}

	var good_q float64
	trans_good, err := h.h.srv.GoodsService.GetGoodByID(r.Context(), userID, good_id)
	if err != nil {
		components.GeneralToastError("That Good Doesn't Exist!").Render(r.Context(), w)
		return
	}

	trans_aud, err := h.h.srv.AudienceService.GetAudienceByID(r.Context(), userID, aud_id)
	if err != nil {
		components.GeneralToastError("That Audience Doesn't Exist!").Render(r.Context(), w)
		return
	}
	var trans_price float64
	if trans_type == string(types.Bought) {
		boughtFrom = aud_id
		good_q = trans_good.Quantity + quantity_f
		trans_price = trans_good.KharidRate * quantity_f
	} else if trans_type == string(types.Sold) {
		soldTo = aud_id
		if quantity_f > trans_good.Quantity {
			components.GeneralToastError("Not Enough Quantity!").Render(r.Context(), w)
			return
		}
		good_q = trans_good.Quantity - quantity_f
		trans_price = trans_good.BikriRate * quantity_f
	}

	if payment == "on" {
		payment_b = true
	}

	transaction := types.Transaction{
		GoodID:     good_id,
		Quantity:   quantity_f,
		Price:      trans_price,
		BoughtFrom: boughtFrom,
		SoldTo:     soldTo,
		Type:       types.TransactionType(trans_type),
		Payment:    payment_b,
		UserID:     userID,
	}
	_, err = h.srv.InsertTransaction(h.h.ctx, transaction)
	if err != nil {
		components.GeneralToastError("Error with transaction service.").Render(r.Context(), w)
		return
	}
	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), userID, trans_good.ID, types.UpdateGood{
		Name:       trans_good.Name,
		Unit:       trans_good.Unit,
		KharidRate: trans_good.KharidRate,
		BikriRate:  trans_good.BikriRate,
		Quantity:   good_q,
	})
	if err != nil {
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	if payment_b {
		trans_aud.Paid += transaction.Price
	} else if !payment_b && trans_type == string(types.Sold) {
		trans_aud.ToReceive += transaction.Price
	} else if !payment_b && trans_type == string(types.Bought) {
		trans_aud.ToPay += transaction.Price
	}
	_, err = h.h.srv.AudienceService.UpdateAudience(r.Context(), trans_aud)
	if err != nil {
		components.GeneralToastError("Error with audience service.").Render(r.Context(), w)
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

	for i := len(soldTrans) - 1; i >= 0; i-- {
		v := soldTrans[i]
		good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), v.UserID, v.GoodID)
		soldAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.UserID, v.SoldTo)
		// boughtAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.BoughtFrom)
		client_Trans = append(client_Trans, types.Transaction_Client{
			TransactionID: v.ID,
			GoodName:      good.Name,
			GoodUnit:      good.Unit,
			Quantity:      strconv.FormatFloat(v.Quantity, 'f', 2, 64),
			Price:         strconv.FormatFloat(v.Price, 'f', 2, 64),
			SoldTo:        soldAudience.Name,
			Payment:       v.Payment,
			// BoughtFrom: boughtAudience.Name,
		})
	}
	pages.Sold(client_Trans).Render(r.Context(), w)
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
	for i := len(boughts) - 1; i >= 0; i-- {
		v := boughts[i]
		good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), v.UserID, v.GoodID)
		// soldAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.SoldTo)
		boughtAudience, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), v.UserID, v.BoughtFrom)
		client_T = append(client_T, types.Transaction_Client{
			TransactionID: v.ID,
			GoodName:      good.Name,
			GoodUnit:      good.Unit,
			Quantity:      strconv.FormatFloat(v.Quantity, 'f', 2, 64),
			Price:         strconv.FormatFloat(v.Price, 'f', 2, 64),
			// SoldTo:     soldAudience.Name,
			Payment:    v.Payment,
			BoughtFrom: boughtAudience.Name,
		})
	}
	pages.Bought(client_T).Render(r.Context(), w)
}

func (h *TransactionHandler) GetInternal(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	internals, err := h.h.srv.TransactionService.GetInternalTransactions(r.Context(), user.ID)
	if err != nil {
		components.GeneralToastError("Error with service.")
		return
	}
	client_T := []types.Transaction_Client{}
	for i := len(internals) - 1; i >= 0; i-- {
		v := internals[i]
		good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), v.UserID, v.GoodID)
		client_T = append(client_T, types.Transaction_Client{
			TransactionID: v.ID,
			GoodName:      good.Name,
			GoodUnit:      good.Unit,
			Quantity:      strconv.FormatFloat(v.Quantity, 'f', 2, 64),
			Price:         strconv.FormatFloat(v.Price, 'f', 2, 64),
			Payment:       v.Payment,
		})
	}
	pages.Internal(client_T).Render(r.Context(), w)
}

func (h *TransactionHandler) InternalTransaction(w http.ResponseWriter, r *http.Request) {
	before_good_id := r.URL.Query().Get("id")
	after_good_id := r.FormValue("after_good_id")
	quantity := r.FormValue("convert_quantity")
	user_id := r.Context().Value("user_id")

	userID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", user_id))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		components.GeneralToastError("User Not Found. Please Login Again.").Render(r.Context(), w)
		return
	}
	beforeGoodID, err := primitive.ObjectIDFromHex(before_good_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Improper Good/Product.").Render(r.Context(), w)
		return
	}
	afterGoodID, err := primitive.ObjectIDFromHex(after_good_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Improper Good/Product.").Render(r.Context(), w)
		return
	}
	quantity_f, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Quantity is not a number.").Render(r.Context(), w)
		return
	}

	beforeGood, err := h.h.srv.GoodsService.GetGoodByID(r.Context(), userID, beforeGoodID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Before Good Not Found.").Render(r.Context(), w)
		return
	}
	afterGood, err := h.h.srv.GoodsService.GetGoodByID(r.Context(), userID, afterGoodID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("After Good Not Found.").Render(r.Context(), w)
		return
	}
	if beforeGoodID == afterGoodID {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Same Good/Product.").Render(r.Context(), w)
		return
	}

	if beforeGood.Quantity < quantity_f {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Not Enough Quantity.").Render(r.Context(), w)
		return
	}

	transaction := types.Transaction{
		GoodID:   afterGoodID,
		Quantity: quantity_f,
		Price:    afterGood.KharidRate * quantity_f,
		Type:     types.Internal,
		Payment:  true,
		UserID:   userID,
	}
	_, err = h.srv.InsertTransaction(h.h.ctx, transaction)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with transaction service.").Render(r.Context(), w)
		return
	}

	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), userID, beforeGood.ID, types.UpdateGood{
		Name:       beforeGood.Name,
		Unit:       beforeGood.Unit,
		KharidRate: beforeGood.KharidRate,
		BikriRate:  beforeGood.BikriRate,
		Quantity:   beforeGood.Quantity - quantity_f,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), userID, afterGood.ID, types.UpdateGood{
		Name:       afterGood.Name,
		Unit:       afterGood.Unit,
		KharidRate: afterGood.KharidRate,
		BikriRate:  afterGood.BikriRate,
		Quantity:   afterGood.Quantity + quantity_f,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), userID)
	components.TodaysGoods(goods, true).Render(r.Context(), w)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	trans_id := r.URL.Query().Get("id")

	transID, err := primitive.ObjectIDFromHex(trans_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Deletion transaction ID parse error.")
		return
	}

	err = h.srv.DeleteTransaction(r.Context(), user.ID, transID)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Deletion unsuccessful. Service Failure.")
		return
	}
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	trans_id := r.URL.Query().Get("id")
	payment_b := r.URL.Query().Get("payment")

	if payment_b == "" || trans_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	transID, err := primitive.ObjectIDFromHex(trans_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction ID parse error.")
		return
	}
	payment, err := strconv.ParseBool(payment_b)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction payment parse error.")
		return
	}

	transUpdated, err := h.srv.UpdateTransaction(r.Context(), user.ID, transID, payment)
	// fmt.Println("\nunedited:", payment_b, "recvieved:", payment, "updated:", transUpdated.Payment, "\n")
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction service error.")
		return
	}
	pages.CheckboxBoolPayment(transUpdated.ID.Hex(), transUpdated.Payment).Render(r.Context(), w)

}
