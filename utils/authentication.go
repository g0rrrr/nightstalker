package utils

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"nightstalker/models"
	"net/http"
)

var Store *sessions.CookieStore

func GetCookieStore(r *http.Request) *sessions.CookieStore {
	if Store == nil {
        cookie_key := "encrypt_your_cookies"
		Store = sessions.NewCookieStore([]byte(cookie_key))
	}

	return Store
}

func GetCurrentUser(r *http.Request) *models.User {
	cached := context.Get(r, "user")
	if cached != nil {
		return cached.(*models.User)
	}

	session, _ := GetCookieStore(r).Get(r, "sirsid")

	if session.Values["username"] == nil || session.Values["password"] == nil {
		return nil
	}
	err, current_user := models.AuthenticateUser(session.Values["username"].(string), session.Values["password"].(string))

	if err != nil {
		return nil
	}

	context.Set(r, "user", current_user)
	return current_user
}
