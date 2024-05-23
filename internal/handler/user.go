package handler

import (
	"net/http"

	"github.com/axyut/dairygo/internal/service"
)

type UserHandler struct {
	h   *Handler
	srv *service.UserService
}

func (user *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user.h.logger.Info("UserHandler.CreateUser", "Method", r.Method)

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username, email, password := r.URL.Query().Get("username"), r.URL.Query().Get("email"), r.URL.Query().Get("password")
	if username == "" || email == "" || password == "" {
		username = r.FormValue("username")
		email = r.FormValue("email")
		password = r.FormValue("password")
		if username == "" || email == "" || password == "" {
			// send error template
			http.Error(w, "Empty Fields!", http.StatusBadRequest)
			return
		}
	}
	err := user.srv.NewUser(user.h.ctx, username, email, password)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return
	}
	// send success template
	w.Write([]byte("User Created!"))
}
