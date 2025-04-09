package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

type userMiddleware func(w http.ResponseWriter, r *http.Request, user database.User)

type client struct {
	lastReq   time.Time
	tokens    int
	maxTokens int
	mu        sync.Mutex
}

type limiter struct {
	clients map[string]*client
	mu      sync.Mutex
}

func newLimiter() *limiter {
	return &limiter{
		clients: make(map[string]*client),
	}
}

func (l *limiter) limiter(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.mu.Lock()
		c, ok := l.clients[r.RemoteAddr]
		if !ok {
			c = &client{
				lastReq:   time.Now(),
				tokens:    10,
				maxTokens: 10,
			}
			l.clients[r.RemoteAddr] = c
		}
		l.mu.Unlock()

		c.mu.Lock()
		defer c.mu.Unlock()

		elapsed := time.Since(c.lastReq)

		refillRate := time.Second * 2
		tokensToAdd := int(elapsed / refillRate)

		if tokensToAdd > 0 {
			c.tokens = min(c.maxTokens, c.tokens+tokensToAdd)
		}

		// Check if the client has tokens to make a request
		if c.tokens > 0 {
			c.tokens--
			c.lastReq = time.Now()
			handle.ServeHTTP(w, r)
		} else {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}

func gracefulerror(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "internal error"

				log.Printf(msg, err, string(debug.Stack()))
				http.Error(w, "interncal server error", http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

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
// 		session, err := gothic.Store.Get(r, os.Getenv("SESSION_NAME"))
// 		if err != nil {
// 			ResponseWithError(w, err.Error(), 401)
// 			return
// 		}

// 		authHeader := r.Header.Get("Authorization")
// 		header_format := strings.Split(authHeader, " ")
// 		if len(header_format) < 2 {
// 			ResponseWithError(w, "invalid authorization header", 401)
// 			return
// 		}

// 		userId, ok := session.Values[header_format[1]]
// 		if !ok {
// 			ResponseWithError(w, "invalid user id", 401)
// 			return
// 		}

// 		r = r.WithContext(
// 			context.WithValue(r.Context(), "user_id", userId))
// 		next.ServeHTTP(w, r)
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

func (s *my_server) checkUserWorkspaceRoleMiddleware(handler userMiddleware) userMiddleware {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		w_id := r.PathValue("workspace_id")
		workspace_id, err := uuid.Parse(w_id)
		if err != nil {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		workspace, err := s.db.GetWorkspaceByID(r.Context(), workspace_id)
		if err != nil {
			if err == sql.ErrNoRows {
				ResponseWithError(w, "workspace not found", http.StatusNotFound)
			} else {
				ResponseWithError(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		workspace_member, err := s.db.GetWorkspaceUserRole(r.Context(), database.GetWorkspaceUserRoleParams{
			UserID:      user.ID,
			WorkspaceID: workspace.ID,
		})

		if err != nil {
			if err == sql.ErrNoRows {
				ResponseWithError(w, "user not found in workspace", http.StatusNotFound)
			} else {
				ResponseWithError(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user_role", workspace_member.Role))
		handler(w, r, user)
	}
}
