package apiserver

import (
	"net/http"

	"github.com/Vacym/neighbors-force/internal/proxyserver"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func Start(config *proxyserver.Config) error {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}

	s := newServer(sessionStore, logLevel)

	return http.ListenAndServe(config.BindAddrApi, s)
}
