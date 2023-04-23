package apiserver

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func Start(config *Config) error {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(sessionStore)

	return http.ListenAndServe(config.BindAddr, s)
}
