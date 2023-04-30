package proxyserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func Start() error {
	apiTargetURL, _ := url.Parse("http://localhost:8081")
	apiProxy := httputil.NewSingleHostReverseProxy(apiTargetURL)

	htmlTargetURL, _ := url.Parse("http://localhost:8082")
	htmlProxy := httputil.NewSingleHostReverseProxy(htmlTargetURL)

	r := mux.NewRouter()
	r.PathPrefix("/api/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		apiProxy.ServeHTTP(w, r)
	}).Methods("GET", "POST", "PUT", "DELETE")

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlProxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(":8080", r)
}
