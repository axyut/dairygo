package handler

import (
	"net/http"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
)

type UserHandler struct {
	h   *Handler
	srv *service.UserService
}

func (user *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user.h.logger.Info("UserHandler.CreateUser", "Method", r.Method)

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		registerTemp := pages.RegisterPage()
		client.Layout(registerTemp, "Register DairyGo").Render(user.h.ctx, w)
		return
	}

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
			// http.Error(w, "Empty Fields!", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			pages.RegisterError("Invalid Credentials.").Render(r.Context(), w)
			return
		}
	}
	err := user.srv.NewUser(user.h.ctx, username, email, password)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		// w.WriteHeader(http.StatusExpectationFailed)
		pages.RegisterError("Server Error.").Render(r.Context(), w)
		return
	}
	// send success template
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func (user *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	user.h.logger.Info("UserHandler.LoginUser", "Method", r.Method)

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		loginTemp := pages.Login()
		client.Layout(loginTemp, "Login DairyGo").Render(user.h.ctx, w)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email_username, password := r.URL.Query().Get("email_username"), r.URL.Query().Get("password")

	if email_username == "" || password == "" {
		email_username = r.FormValue("email_username")
		password = r.FormValue("password")

		if email_username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Empty Fields!", http.StatusBadRequest)
			logErr := pages.LoginError("Empty Fields!")
			logErr.Render(user.h.ctx, w)
			return
		}
	}
	err := user.srv.LoginUser(user.h.ctx, email_username, password)
	if err != nil {
		// http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		w.WriteHeader(http.StatusUnauthorized)
		// pages.Toast("Invalid Credentials!")
		logErr := pages.LoginError("Invalid Credentials!")
		logErr.Render(user.h.ctx, w)

		return
	}
	// 	cookieValue := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", sessionID, userID)))

	// expiration := time.Now().Add(365 * 24 * time.Hour)
	// cookie := http.Cookie{
	// 	Name:     h.sessionCookieName,
	// 	Value:    cookieValue,
	// 	Expires:  expiration,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// }
	// http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (user *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user.h.logger.Info("UserHandler.GetUser", "Method", r.Method)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// userData := user.srv.GetUser(user.h.ctx)
	// if userData == nil {
	// 	http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(userData)
}
