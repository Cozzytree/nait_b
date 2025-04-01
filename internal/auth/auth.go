package auth

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func InitAuth() {
	gob.Register(time.Time{})
	storeKey := []byte(os.Getenv("COOKIE_STORE_KEY"))
	if len(storeKey) == 0 {
		panic("no store key")
	}

	var GOOGLE_CLIENT = os.Getenv("GOOGLE_CLIENT_ID")
	var GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_SECRET")
	var GOOGLE_CALLBACK = os.Getenv("GOOGLE_CALLBACK_URL")

	var GITHUB_CLIENT = os.Getenv("GITHUB_CLIENT_ID")
	var GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_SECRET")
	var GITHUB_CALLBACK = os.Getenv("GITHUB_CALLBACK_URL")

	if GOOGLE_CLIENT == "" || GOOGLE_CLIENT_SECRET == "" || GOOGLE_CALLBACK == "" {
		panic("invalid auth settings")
	}

	if GITHUB_CLIENT == "" || GITHUB_CLIENT_SECRET == "" || GITHUB_CALLBACK == "" {
		panic("invalid auth settings")
	}

	// Create cookie store and configure options
	store := sessions.NewCookieStore([]byte(storeKey))
	store.Options = &sessions.Options{
		HttpOnly: true,           // Prevents access to cookies via JavaScript
		MaxAge:   int(time.Hour), // Set session expiration (7 days)
		Secure:   true,           // Set to true in production for HTTPS
		Path:     "/",            // Cookie is available throughout the site
	}

	gothic.Store = store
	goth.UseProviders(
		google.New(GOOGLE_CLIENT, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK),
		github.New(GITHUB_CLIENT, GITHUB_CLIENT_SECRET, GITHUB_CALLBACK),
	)
}
