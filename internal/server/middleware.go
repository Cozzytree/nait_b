package server

import (
	"context"
	"net/http"
	"os"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/markbates/goth/gothic"
)

type userMiddleware func(w http.ResponseWriter, r *http.Request, user database.User)

func getUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := gothic.Store.Get(r, os.Getenv("SESSION_NAME"))
		if err != nil {
			ResponseWithError(w, err.Error(), 401)
			return
		}
		userId, ok := session.Values["user_id"]
		if !ok {
			ResponseWithError(w, "invalid user id", 401)
			return
		}

		r = r.WithContext(
			context.WithValue(r.Context(), "user_id", userId))
		next.ServeHTTP(w, r)
	})
}

// func getUserSession(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		auth := r.Header.Get("Authorization")
// 		bearer := strings.Split(auth, " ")
// 		if len(bearer) < 2 {
// 			ResponseWithError(w, "invalid authorization header", 401)
// 			return
// 		}

// 		if bearer[0] != "nait" {
// 			ResponseWithError(w, "invalid authorization header", 401)
// 			return
// 		}

// 		token := bearer[1]
// 		if token == "" {
// 			ResponseWithError(w, "invalid token", 401)
// 			return
// 		}

// 		r = r.WithContext(context.WithValue(r.Context(), "user_id", token))
// 	})
// }

func (s *my_server) getUserMiddleware(handler userMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("user_id")

		user_id, ok := id.(string)
		if !ok {
			ResponseWithError(w, "invalid user id", 404)
			return
		}

		user, err := s.db.GetUser(r.Context(), user_id)
		if err != nil {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handler(w, r, user)
	}
}
