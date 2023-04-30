package apiserver

import (
	"net/http"

	"github.com/Vacym/neighbors-force/internal/proxyserver"
	"github.com/gorilla/sessions"
)

func Start(config *proxyserver.Config) error {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(sessionStore)

	return http.ListenAndServe(config.BindAddrApi, s)
}
