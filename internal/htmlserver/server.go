package htmlserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

type htmlServer struct {
	router *mux.Router
}

func NewServer() *htmlServer {
	s := &htmlServer{
		router: mux.NewRouter(),
	}

	s.configureRouter()

	return s
}

func (s *htmlServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *htmlServer) configureRouter() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/index.html")
	})

	s.router.PathPrefix("/sources/").
		Handler(http.StripPrefix(
			"/sources/",
			http.FileServer(http.Dir("site/sources/")),
		))

}
