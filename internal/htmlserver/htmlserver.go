package htmlserver

import (
	"net/http"
)

func Start(bindAddr string) error {
	s := NewServer()

	return http.ListenAndServe(bindAddr, s)
}
