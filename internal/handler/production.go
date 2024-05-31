package handler

import (
	"net/http"
	"strconv"

	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductionHandler struct {
	h   *Handler
	srv *service.ProductionService
}

func (h *ProductionHandler) NewProduction(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	change_goodID := r.URL.Query().Get("change_good_id")
	change_goodQuantity := r.FormValue("change_quantity")
	prod_goodID := r.FormValue("prod_good_id")
	prod_goodQuantity := r.FormValue("prod_quantity")

	if change_goodID == "" || prod_goodID == "" || change_goodQuantity == "" || prod_goodQuantity == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Empty Fields!").Render(r.Context(), w)
		return
	}

	Change_goodID, errC := primitive.ObjectIDFromHex(change_goodID)
	Prod_goodID, errP := primitive.ObjectIDFromHex(prod_goodID)

	if errC != nil || errP != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Invalid ID.").Render(r.Context(), w)
		return
	}

	C_quantity, errCQ := strconv.ParseFloat(change_goodQuantity, 64)
	P_quantity, errPQ := strconv.ParseFloat(prod_goodQuantity, 64)

	if errCQ != nil || errPQ != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Enter Numeric Value for Quantity.").Render(r.Context(), w)
		return
	}

	var profit float64
	var loss float64

	Cgood, errCG := h.h.srv.GoodsService.GetGoodByID(r.Context(), user.ID, Change_goodID)
	Pgood, errPG := h.h.srv.GoodsService.GetGoodByID(r.Context(), user.ID, Prod_goodID)
	if errCG != nil || errPG != nil {
		w.WriteHeader(http.StatusConflict)
		components.GeneralToastError("Good of that ID not found.").Render(r.Context(), w)
		return
	}
	change_goodName := Cgood.Name
	change_goodUnit := Cgood.Unit
	prod_goodName := Pgood.Name
	prod_goodUnit := Pgood.Unit

	if Cgood.ID == Pgood.ID {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Cannot Produce same good.").Render(r.Context(), w)
		return
	}

	if C_quantity < P_quantity || Cgood.Quantity < C_quantity {
		w.WriteHeader(http.StatusBadRequest)
		components.GeneralToastError("Cannot Produce. Not Enough Quantity.").Render(r.Context(), w)
		return
	}

	// when kharid rate is later implemented according to bought transactions, this will be updated. idk how
	if (Cgood.KharidRate * C_quantity) < (Pgood.KharidRate * P_quantity) {
		profit = (Pgood.KharidRate * P_quantity) - (Cgood.KharidRate * C_quantity)
		loss = 0
	} else {
		loss = (Cgood.KharidRate * C_quantity) - (Pgood.KharidRate * P_quantity)
		profit = 0
	}

	insertProd := types.Production{
		ChangeGoodID:     Change_goodID,
		ProducedGoodID:   Prod_goodID,
		ChangeGoodName:   change_goodName,
		ChangeQuantity:   C_quantity,
		ProducedGoodName: prod_goodName,
		ChangeGoodUnit:   change_goodUnit,
		ProducedGoodUnit: prod_goodUnit,
		ProducedQuantity: P_quantity,
		Profit:           profit,
		Loss:             loss,
		UserID:           user.ID,
	}

	_, err := h.srv.InsertProduction(h.h.ctx, insertProd, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	// update both goods quantity value
	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), user.ID, Cgood.ID, types.UpdateGood{
		Name:       Cgood.Name,
		Unit:       Cgood.Unit,
		KharidRate: Cgood.KharidRate,
		BikriRate:  Cgood.BikriRate,
		Quantity:   Cgood.Quantity - C_quantity,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	_, err = h.h.srv.GoodsService.UpdateGood(r.Context(), user.ID, Pgood.ID, types.UpdateGood{
		Name:       Pgood.Name,
		Unit:       Pgood.Unit,
		KharidRate: Pgood.KharidRate,
		BikriRate:  Pgood.BikriRate,
		Quantity:   Pgood.Quantity + P_quantity,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.GeneralToastError("Error with goods service.").Render(r.Context(), w)
		return
	}

	goods, _ := h.h.srv.GoodsService.GetAllGoods(r.Context(), user.ID)
	components.GoodsTable(goods, true).Render(r.Context(), w)
}

func (h *ProductionHandler) DeleteProduction(w http.ResponseWriter, r *http.Request) {
	prod_id := r.URL.Query().Get("id")
	user := h.h.UserHandler.GetUser(w, r)

	prodID, _ := primitive.ObjectIDFromHex(prod_id)

	err := h.srv.DeleteProduction(h.h.ctx, prodID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
}

func (h *ProductionHandler) GetProductionPage(w http.ResponseWriter, r *http.Request) {
	user := h.h.UserHandler.GetUser(w, r)
	prods, err := h.srv.GetAllProductions(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}
	pages.Production(prods).Render(r.Context(), w)
}
