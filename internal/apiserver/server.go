package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type key int

const (
	sessionName     = "testAuth"
	ctxKeyUser  key = iota
)

type apiServer struct {
	router       *mux.Router
	sessionStore sessions.Store
	activeUsers  map[string]*User
}

func newServer(sessionStore sessions.Store) *apiServer {
	s := &apiServer{
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
		activeUsers:  make(map[string]*User),
	}

	s.configureRouter()

	return s
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *apiServer) configureRouter() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.Use(s.UserMiddleware)
	apiRouter.HandleFunc("/game/create", s.createGame()).Methods("POST")

}

func (s *apiServer) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userID, ok := session.Values["user_id"].(string)
		if !ok {
			userID = uuid.New().String()
			session.Values["user_id"] = userID
			session.Save(r, w)
		}

		user, ok := s.activeUsers[userID]
		if !ok {
			fmt.Println("create new user")
			user = NewUser()
			s.activeUsers[userID] = user
		}

		ctx := context.WithValue(r.Context(), ctxKeyUser, user)
		fmt.Println("add user to context")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *apiServer) createGame() http.HandlerFunc {
	type request struct {
		Rows       int `json:"rows"`
		Cols       int `json:"cols"`
		NumPlayers int `json:"num_players"`
	}

	fmt.Println("logging")

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		g, err := game.NewGame(req.Rows, req.Cols, req.NumPlayers)

		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*User)
		user.Game = g

		s.respond(w, r, http.StatusCreated, g.ToMap())
	}
}

func (s *apiServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{
		"error:": err.Error(),
	})
}

func (s *apiServer) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
