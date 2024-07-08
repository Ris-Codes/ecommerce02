package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("SECRET"))

func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session-name")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // Session timeout set to 1 hour (3600 seconds)
		HttpOnly: true,
	}
	return session
}

func SaveSession(r *http.Request, w http.ResponseWriter, session *sessions.Session) {
	session.Save(r, w)
}
