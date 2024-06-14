package handler

import (
	"fmt"
	"math"
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
	user := h.h.UserHandler.GetUser(w, r)
	buying_rate := r.FormValue("buying_rate")
	unit := r.FormValue("unit")
	good_name := r.FormValue("goodName")

	var err error
	var boughtFromID primitive.ObjectID
	var boughtFrom string
	var soldToID primitive.ObjectID
	var soldTo string
	var payment_b bool = false
	var advanced bool = false
	var trans_good types.Good
	var good_change_quantity float64
	var totalPrice float64
	var goodUnit string
	var buyingRate float64
	var goodName string
	var Rate float64

	if goodID == "" || quantity == "" || audienceID == "" || trans_type == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	aud_id, Aerr := primitive.ObjectIDFromHex(audienceID)
	quanTity, ferr := strconv.ParseFloat(quantity, 64)

	if Aerr != nil || ferr != nil {
		h.h.logger.Error("Error while Parsing in Handler.", Aerr, ferr)
		components.GeneralToastError("Error while Parsing in Handler.").Render(r.Context(), w)
		return
	}
	quanTity = math.Abs(quanTity)

	trans_aud, errA := h.h.srv.AudienceService.GetAudienceByID(r.Context(), user.ID, aud_id)
	if errA != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("That Audience Doesn't Exist!").Render(r.Context(), w)
		return
	}

	if trans_type != string(types.Bought) && trans_type != string(types.Sold) {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Invalid Transaction Type.").Render(r.Context(), w)
		return
	}

	if good_name == "" && unit == "" {
		good_id, Gerr := primitive.ObjectIDFromHex(goodID)
		if Gerr != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			components.GeneralToastError("Invalid Good ID.").Render(r.Context(), w)
			return
		}
		trans_good, err = h.h.srv.GoodsService.GetGoodByID(r.Context(), user.ID, good_id)
		goodName = trans_good.Name
		goodUnit = trans_good.Unit
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			components.GeneralToastError("That Good Doesn't Exist!").Render(r.Context(), w)
			return
		}
	} else {
		advanced = true
		goodName = good_name
		goodUnit = unit
	}

	// increase decrease goods quantity
	if trans_type == string(types.Bought) {
		boughtFromID = trans_aud.ID
		boughtFrom = trans_aud.Name

		if buying_rate != "" {
			buyingRate, err = strconv.ParseFloat(buying_rate, 64)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				components.GeneralToastError("Invalid Buying Rate.").Render(r.Context(), w)
				return
			}
			buyingRate = math.Abs(buyingRate)
		} else {
			buyingRate = trans_aud.MapRates[trans_good.ID.Hex()]
			if buyingRate == 0 {
				w.WriteHeader(http.StatusNotAcceptable)
				components.GeneralToastError("Set Buying Rate for that Good First.").Render(r.Context(), w)
				return
			}
		}
		Rate = buyingRate

		good_change_quantity = trans_good.Quantity + quanTity
	} else if trans_type == string(types.Sold) {
		soldToID = trans_aud.ID
		soldTo = trans_aud.Name
		Rate = trans_good.SellingRate

		if quanTity > trans_good.Quantity {
			w.WriteHeader(http.StatusNotAcceptable)
			components.GeneralToastError("Not Enough Quantity!").Render(r.Context(), w)
			return
		}
		good_change_quantity = trans_good.Quantity - quanTity
	}

	if payment == "on" {
		payment_b = true
	}
	totalPrice = Rate * quanTity

	// update to pay and to receive
	toPay := trans_aud.ToPay
	toReceive := trans_aud.ToReceive
	trans_aud.ToPay, trans_aud.ToReceive = utils.SetToPayToRecieve(trans_type, payment_b, totalPrice, toPay, toReceive)

	// fmt.Println("boughtFrom =", boughtFrom, "\nsoldTo =", soldTo, "\ngoodName =", goodName, "\ngoodUnit =", goodUnit, "\ngood_change_quantity =", good_change_quantity, "\ntotalPrice =", totalPrice, "\npayment_b =", payment_b, "\nquanTity =", quanTity, "\nbuyingRate =", buyingRate)
	transaction := types.Transaction{
		ID:              primitive.NewObjectIDFromTimestamp(time.Now()),
		GoodID:          trans_good.ID,
		GoodName:        goodName,
		GoodUnit:        goodUnit,
		Quantity:        quanTity,
		Price:           totalPrice,
		BoughtFrom:      boughtFrom,
		BoughtFromID:    boughtFromID,
		SoldToID:        soldToID,
		SoldTo:          soldTo,
		Type:            types.TransactionType(trans_type),
		Payment:         payment_b,
		ChangeToPay:     toPay - trans_aud.ToPay,
		ChangeToReceive: toReceive - trans_aud.ToReceive,
		AudToPay:        toPay,
		AudToReceive:    toReceive,
		UserID:          user.ID,
		CreationTime:    utils.GetMongoTimeFromHTMLDate(date, time.Now()),
	}

	_, err = h.srv.InsertTransaction(h.h.ctx, transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with transaction service.").Render(r.Context(), w)
		return
	}

	if !advanced {
		_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), user.ID, trans_good.ID, types.UpdateGood{
			Name:        trans_good.Name,
			Unit:        trans_good.Unit,
			SellingRate: trans_good.SellingRate,
			Quantity:    good_change_quantity,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
			return
		}
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

	trans, _ := h.srv.GetTransactionByID(r.Context(), transID, user.ID)
	err = h.srv.DeleteTransaction(r.Context(), user.ID, transID)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Deletion unsuccessful. Service Failure.")
		return
	}

	var audID primitive.ObjectID
	good, _ := h.h.srv.GoodsService.GetGoodByID(r.Context(), user.ID, trans.GoodID)
	if trans.Type == types.Sold {
		good.Quantity += trans.Quantity
		audID = trans.SoldToID
	} else if trans.Type == types.Bought {
		good.Quantity -= trans.Quantity
		audID = trans.BoughtFromID
	}

	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), user.ID, good.ID, types.UpdateGood{
		Name:        good.Name,
		Unit:        good.Unit,
		SellingRate: good.SellingRate,
		Quantity:    good.Quantity,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.")
		return
	}

	trans_aud, _ := h.h.srv.AudienceService.GetAudienceByID(r.Context(), user.ID, audID)

	// if paila ko change minus ma xa (toPay before was less than toPay after)
	// teslai aile ko ma ghatayesi, paila jati thiyo teti aauxa
	if trans.ChangeToPay < 0 {
		trans_aud.ToPay -= math.Abs(trans.ChangeToPay)
	} else {
		trans_aud.ToPay += trans.ChangeToPay
	}

	if trans.ChangeToReceive < 0 {
		trans_aud.ToReceive -= math.Abs(trans.ChangeToReceive)
	} else {
		trans_aud.ToReceive += trans.ChangeToReceive
	}
	_, err = h.h.srv.AudienceService.UpdateAudience(r.Context(), trans_aud)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with audience service.")
		return
	}

	w.WriteHeader(http.StatusOK)
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

	oneTrans, err := h.srv.GetTransactionByID(r.Context(), transID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction service error.")
		return
	}

	// update audience
	audID := oneTrans.SoldToID
	if oneTrans.Type == types.Bought {
		audID = oneTrans.BoughtFromID
	}
	aud, err := h.h.srv.AudienceService.GetAudienceByID(r.Context(), user.ID, audID)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction audience service error.")
		return
	}

	updateToPayToRecieveOnPayment(&oneTrans, payment, &aud)
	oneTrans.Payment = payment
	oneTrans.AudToPay = aud.ToPay
	oneTrans.AudToReceive = aud.ToReceive

	// commit updates
	_, err = h.srv.UpdateTransaction(r.Context(), user.ID, transID, oneTrans)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction service error.")
		return
	}

	_, err = h.h.srv.AudienceService.UpdateAudience(r.Context(), aud)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		components.GeneralToastError("Update transaction audience service error.")
		return
	}

	pages.CheckboxBoolPayment(oneTrans.ID.Hex(), payment).Render(r.Context(), w)

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

func updateToPayToRecieveOnPayment(oneTrans *types.Transaction, payment bool, aud *types.Audience) {
	// fmt.Println("before price:", oneTrans.Price, aud.Name, "toget:", aud.ToReceive, "topay", aud.ToPay)
	if payment {
		// fmt.Println("payed")
		if oneTrans.Type == types.Sold {
			if aud.ToPay > 0 {
				aud.ToPay -= oneTrans.Price
				if aud.ToPay < 0 {
					aud.ToReceive += math.Abs(aud.ToPay) // convert to positive
					aud.ToPay = 0
				}
			} else {
				aud.ToReceive -= oneTrans.Price
			}
		} else if oneTrans.Type == types.Bought {
			if aud.ToReceive > 0 {
				aud.ToReceive -= oneTrans.Price
				if aud.ToReceive < 0 {
					aud.ToPay += math.Abs(aud.ToReceive) // convert to positive
					aud.ToReceive = 0
				}
			} else {
				aud.ToPay -= oneTrans.Price
			}
		}
	} else if !payment {
		// fmt.Println("not payed")
		if oneTrans.Type == types.Sold {
			if aud.ToPay > 0 {
				aud.ToPay -= oneTrans.Price
				if aud.ToPay < 0 {
					aud.ToReceive += math.Abs(aud.ToPay) // convert to positive
					aud.ToPay = 0
				}
			} else {
				aud.ToReceive += oneTrans.Price
			}
		} else if oneTrans.Type == types.Bought {
			if aud.ToReceive > 0 {
				aud.ToReceive -= oneTrans.Price
				if aud.ToReceive < 0 {
					aud.ToPay += math.Abs(aud.ToReceive) // convert to positive
					aud.ToReceive = 0
				}
			} else {
				aud.ToPay += oneTrans.Price
			}
		}
	}
	// fmt.Println("after price:", transUpdated.Price, aud.Name, "to recieve:", aud.ToReceive, "topay", aud.ToPay)
}
