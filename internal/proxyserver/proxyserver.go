package proxyserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func Start(config *Config) error {
	apiTargetURL, _ := url.Parse("http://localhost" + config.BindAddrApi)
	apiProxy := httputil.NewSingleHostReverseProxy(apiTargetURL)

	htmlTargetURL, _ := url.Parse("http://localhost" + config.BindAddrHtml)
	htmlProxy := httputil.NewSingleHostReverseProxy(htmlTargetURL)

	r := mux.NewRouter()
	r.PathPrefix("/api/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		apiProxy.ServeHTTP(w, r)
	}).Methods("GET", "POST", "PUT", "DELETE")

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlProxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(config.BindAddrProxy, r)
}
