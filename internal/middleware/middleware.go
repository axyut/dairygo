package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(m ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// excpet on /login and /register
		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			next.ServeHTTP(w, r)
			return
		}

		// Check if user is authenticated
		cookie, err := r.Cookie("cookie")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// attach cookie value to context
		val := decode(cookie.Value)
		ctx := context.WithValue(r.Context(), "user_id", val)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func decode(cookie string) string {
	val, err := base64.RawStdEncoding.DecodeString(cookie)
	if err != nil {
		log.Println("Error decoding cookie", err)
		return ""
	}
	return string(val)
}

func GetUser(ctx context.Context) string {
	user_id := ctx.Value("user_id")
	if user_id == nil {
		return ""
	}

	return fmt.Sprintf("%v", user_id)
}
