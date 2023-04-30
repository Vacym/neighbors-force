package htmlserver

import (
	"net/http"

	"github.com/Vacym/neighbors-force/internal/proxyserver"
)

func Start(config *proxyserver.Config) error {
	s := NewServer()

	return http.ListenAndServe(config.BindAddrHtml, s)
}
