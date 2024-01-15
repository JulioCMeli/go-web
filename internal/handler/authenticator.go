package handler

import (
	"encoding/json"
	"net/http"
)

// NewAuthenticator creates a new authenticator
func NewAuthenticator(token string) *Authenticator {
	return &Authenticator{
		token: token,
	}
}

// Authenticator is a middleware that authenticates the request
type Authenticator struct {
	token string
}

func (a *Authenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("auth")

		// logic before
		token := r.Header.Get("token")
		if token != a.token {
			code := http.StatusForbidden
			body := MyResponse{Message: "Invalid authorization token", Data: nil}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		// call next
		next.ServeHTTP(w, r)

		// logic after
		// ...
	})
}
