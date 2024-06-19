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
	err := user.srv.GetUserByEmail(user.h.ctx, email)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		pages.RegisterError("Email already exists.").Render(r.Context(), w)
		return
	}
	_, err = user.srv.InsertUser(user.h.ctx, username, email, password)
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
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("UserID Error. Re-Login.").Render(r.Context(), w)
		return types.User{}
	}
	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("UserID Error. Re-Login.").Render(r.Context(), w)
		return types.User{}
	}
	userData, err := user.srv.GetUserByID(user.h.ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		components.GeneralToastError("UserID Error. Re-Login.").Render(r.Context(), w)
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
	components.GetUserReq(sendUser, "nav").Render(r.Context(), w)
}

func (user *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:    "cookie",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/login")
}

func (user *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userData := user.GetUser(w, r)
	w.WriteHeader(http.StatusOK)
	pg := pages.Profile(userData)
	client.Layout(pg, "User Profile").Render(r.Context(), w)
}

func (user *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userData := user.GetUser(w, r)

	err := user.srv.DeleteUser(user.h.ctx, userData.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't fullfill your request.").Render(r.Context(), w)
		return
	}

	user.LogoutUser(w, r)
}

func (user *UserHandler) UpdateUserDefaults(w http.ResponseWriter, r *http.Request) {
	userData := user.GetUser(w, r)
	good_id := r.FormValue("def_sell_good_id")
	check := r.FormValue("def_sell_payment")

	if good_id == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Send required fields").Render(r.Context(), w)
		return
	}

	var payment bool
	if check == "on" {
		payment = true
	} else if check == "off" {
		payment = false
	}

	goodID, err := primitive.ObjectIDFromHex(good_id)

	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Invalid Good ID").Render(r.Context(), w)
		return
	}
	good, err := user.h.srv.GoodsService.GetGoodByID(r.Context(), userData.ID, goodID)

	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't get the good from ID.").Render(r.Context(), w)
		return
	}

	if userData.Default == nil {
		userData.Default = make(map[types.UserDefault]string)
	}
	userData.Default[types.SellGood] = good.ID.Hex()
	userData.Default[types.SellPayment] = fmt.Sprintf("%t", payment)

	_, err = user.srv.UpdateUser(user.h.ctx, userData)

	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		components.GeneralToastError("Couldn't update User's Default Values.").Render(r.Context(), w)
		return
	}
	goods, _ := user.h.srv.GoodsService.GetAllGoods(r.Context(), userData.ID)
	w.WriteHeader(http.StatusOK)
	components.TransUnit(goods, good.Name+" set as default.", userData).Render(r.Context(), w)
}
