package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func handleBeginAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *my_server) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Redirect(w, r,
			fmt.Sprintf("http://localhost:5173/error?v=%v", err.Error()),
			http.StatusPermanentRedirect)
		return
	}

	session, err := gothic.Store.New(r, os.Getenv("SESSION_NAME"))
	if err != nil {
		log.Println("session error ", err.Error())
	}

	session.Values["user_id"] = user.UserID
	session.Values["accessToken"] = user.AccessToken
	session.Values["refreshToken"] = user.RefreshToken

	err = session.Save(r, w)
	if err != nil {
		log.Println("error while saving session", err.Error())
		return
	}

	fmt.Println(" ---email----\n", user.Email)
	fmt.Println(" ---user_id----\n", user.UserID)
	fmt.Println(" ---avatar----\n", user.AvatarURL)
	fmt.Println(" ---provider----\n", user.Provider)
	fmt.Println(" ---nickname----\n", user.NickName)

	username := strings.Join([]string{user.FirstName, user.LastName}, " ")
	if len(user.NickName) >= 1 {
		username = user.NickName
	}

	if username == "" {
		username = user.Name
	}

	// create the user
	created_user, err := s.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:     username,
		Email:    user.Email,
		AuthID:   user.UserID,
		Provider: user.Provider,
		Avatar: sql.NullString{
			Valid:  user.AvatarURL != "",
			String: user.AvatarURL,
		},
	})
	if err != nil {
		http.Redirect(w, r,
			strings.Join([]string{os.Getenv("REDIRECT_URL"), "error", fmt.Sprintf("?%v", err.Error())}, "/"),
			http.StatusPermanentRedirect)
		return
	}

	// create user initial workspace
	createdWorkspace, err := s.db.CreateNewWorkspace(r.Context(), database.CreateNewWorkspaceParams{
		Name:   strings.Join([]string{fmt.Sprintf("%v's", username), "workspace"}, " "),
		UserID: created_user,
	})
	if err != nil {
		err = s.db.DeleteUser(r.Context(), created_user)
		if err != nil {
			http.Redirect(w, r,
				fmt.Sprintf("%v/%v", os.Getenv("REDIRECT_URL"), fmt.Sprintf("/error?e=%v", err.Error())),
				http.StatusPermanentRedirect)
		}
	}

	http.Redirect(w, r,
		fmt.Sprintf("%v/%v", os.Getenv("REDIRECT_URL"), createdWorkspace),
		http.StatusPermanentRedirect)
}
