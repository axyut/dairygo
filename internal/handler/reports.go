package handler

import (
	"net/http"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/pages"
)

type ReportsHandler struct {
	h *Handler
}

func (h *ReportsHandler) GetReportsPage(w http.ResponseWriter, r *http.Request) {
	h.h.UserHandler.GetUser(w, r)
	w.WriteHeader(http.StatusAccepted)
	p := pages.ReportsPage()
	client.Layout(p, "Reports").Render(r.Context(), w)
}
