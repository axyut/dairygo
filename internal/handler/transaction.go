package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"github.com/axyut/dairygo/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionHandler struct {
	h   *Handler
	srv *service.TransactionService
}

func (h *TransactionHandler) GetTransactionPage(w http.ResponseWriter, r *http.Request) {
	soldTrans := h.getSoldTransClient(w, r)
	page := pages.TransactionPage(soldTrans)
	client.Layout(page, "Transactions").Render(r.Context(), w)
}

func (h *TransactionHandler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	goodID := r.FormValue("goodID")
	quantity := r.FormValue("quantity")
	audienceID := r.FormValue("audienceID")
	trans_type := r.FormValue("type")
	payment := r.FormValue("payment")
	date := r.FormValue("date")
	user_id := r.Context().Value("user_id")

	var boughtFromID primitive.ObjectID
	var boughtFrom string
	var soldToID primitive.ObjectID
	var soldTo string
	var payment_b bool = false

	if goodID == "" || quantity == "" || audienceID == "" || trans_type == "" || user_id == "" {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Invalid Transaction Type.").Render(r.Context(), w)
		return
	}

	var good_q float64
	trans_good, err := h.h.srv.GoodsService.GetGoodByID(r.Context(), userID, good_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("That Good Doesn't Exist!").Render(r.Context(), w)
		return
	}

	trans_aud, err := h.h.srv.AudienceService.GetAudienceByID(r.Context(), userID, aud_id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("That Audience Doesn't Exist!").Render(r.Context(), w)
		return
	}
	var trans_price float64
	if trans_type == string(types.Bought) {
		boughtFromID = trans_aud.ID
		boughtFrom = trans_aud.Name
		good_q = trans_good.Quantity + quantity_f
		trans_price = trans_good.KharidRate * quantity_f
	} else if trans_type == string(types.Sold) {
		soldToID = trans_aud.ID
		soldTo = trans_aud.Name
		if quantity_f > trans_good.Quantity {
			w.WriteHeader(http.StatusNotAcceptable)
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
		ID:           primitive.NewObjectIDFromTimestamp(time.Now()),
		GoodID:       good_id,
		Quantity:     quantity_f,
		Price:        trans_price,
		BoughtFrom:   boughtFrom,
		BoughtFromID: boughtFromID,
		SoldToID:     soldToID,
		SoldTo:       soldTo,
		Type:         types.TransactionType(trans_type),
		Payment:      payment_b,
		UserID:       userID,
		CreationTime: utils.GetMongoTimeFromHTMLDate(date),
	}

	_, err = h.srv.InsertTransaction(h.h.ctx, transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	if !payment_b && trans_type == string(types.Sold) {
		trans_aud.ToReceive += transaction.Price
	} else if !payment_b && trans_type == string(types.Bought) {
		trans_aud.ToPay += transaction.Price
	}
	_, err = h.h.srv.AudienceService.UpdateAudience(r.Context(), trans_aud)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with audience service.").Render(r.Context(), w)
		return
	}

	components.GeneralToastSuccess("Transaction Added Successfully!").Render(r.Context(), w)
}

func (h *TransactionHandler) getSoldTransClient(w http.ResponseWriter, r *http.Request) (soldTrans []types.Transaction) {
	user := h.h.UserHandler.GetUser(w, r)

	soldTrans, err := h.h.srv.TransactionService.GetSoldTransactions(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with service.")
		return
	}
	return
}

func (h *TransactionHandler) GetSold(w http.ResponseWriter, r *http.Request) {
	soldTrans := h.getSoldTransClient(w, r)
	pages.Sold(soldTrans).Render(r.Context(), w)
}

func (h *TransactionHandler) GetBought(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		w.WriteHeader(http.StatusUnauthorized)
		components.GeneralToastError("Couldn't Identify User. Please Login Again.")
	}

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Couldn't Fullfill Your Request. Please Try Again.")
	}
	boughts, err := h.h.srv.TransactionService.GetBoughtTransactions(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with service.")
		return
	}

	pages.Bought(boughts).Render(r.Context(), w)
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
	// update audience
	audID := transUpdated.SoldToID
	if transUpdated.Type == types.Bought {
		audID = transUpdated.BoughtFromID
	}
	aud, err := h.h.srv.AudienceService.GetAudienceByID(r.Context(), user.ID, audID)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction audience service error.")
		return
	}
	fmt.Println("before price:", transUpdated.Price, aud.Name, "toget:", aud.ToReceive, "topay", aud.ToPay)
	if payment {
		fmt.Println("payed")
		if transUpdated.Type == types.Sold {
			aud.ToReceive -= transUpdated.Price
		} else if transUpdated.Type == types.Bought {
			aud.ToPay -= transUpdated.Price
		}
	} else if !payment {
		fmt.Println("not payed")
		if transUpdated.Type == types.Sold {
			aud.ToReceive += transUpdated.Price
		} else if transUpdated.Type == types.Bought {
			aud.ToPay += transUpdated.Price
		}
	}
	fmt.Println("after price:", transUpdated.Price, aud.Name, "to recieve:", aud.ToReceive, "topay", aud.ToPay)
	_, err = h.h.srv.AudienceService.UpdateAudience(r.Context(), aud)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction audience service error.")
		return
	}

	pages.CheckboxBoolPayment(transUpdated.ID.Hex(), transUpdated.Payment).Render(r.Context(), w)

}

func (h *TransactionHandler) DeleteAllTransactions(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	trans_type := r.URL.Query().Get("type")
	transType := types.TransactionType(trans_type)
	if trans_type == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	if transType != types.Sold && transType != types.Bought {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}
	err := h.srv.DeleteAllTransactions(r.Context(), user.ID, transType)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Deletion unsuccessful. Service Failure.")
		return
	}

	w.WriteHeader(http.StatusOK)
	components.GeneralToastSuccess("Deleted Successfully").Render(r.Context(), w)
}
