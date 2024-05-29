package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/client/components"
	"github.com/axyut/dairygo/client/pages"
	"github.com/axyut/dairygo/internal/service"
	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	h   *Handler
	srv *service.UserService
}

func (user *UserHandler) NewUser(w http.ResponseWriter, r *http.Request) {

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
			w.WriteHeader(http.StatusBadRequest)
			pages.RegisterError("Invalid Credentials.").Render(r.Context(), w)
			return
		}
	}
	_, err := user.srv.InsertUser(user.h.ctx, username, email, password)
	if err != nil {
		user.h.logger.Error("UserHandler.CreateUser", "Error", err)
		pages.RegisterError("Server Error.").Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func (user *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

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
			pages.LoginError("Empty Fields!").Render(r.Context(), w)
			return
		}
	}
	newuser, err := user.srv.LoginUser(user.h.ctx, email_username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		pages.LoginError("Invalid Credentials!").Render(r.Context(), w)
		return
	}

	cookieValue := base64.StdEncoding.EncodeToString([]byte(newuser.ID.Hex()))
	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{
		Name:     "cookie",
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (user *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) types.User {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return types.User{}
	}
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return types.User{}
	}
	userData, err := user.srv.GetUserByID(user.h.ctx, id)
	if err != nil {
		http.Error(w, "Couldn't fullfill your request.", http.StatusExpectationFailed)
		return types.User{}
	}
	return userData
}

func (user *UserHandler) GetUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		components.GeneralToastError("Please login.").Render(r.Context(), w)
		return
	}
	sendUser := user.GetUser(w, r)
	w.WriteHeader(http.StatusOK)
	components.GetUserReq(sendUser).Render(r.Context(), w)
}

func (user *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie := http.Cookie{
		Name:    "cookie",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func (user *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userData := user.GetUser(w, r)
	w.WriteHeader(http.StatusOK)
	pg := pages.Profile(userData)
	client.Layout(pg, "User Profile").Render(r.Context(), w)
}
