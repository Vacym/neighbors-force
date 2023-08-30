package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Vacym/neighbors-force/internal/bot"
	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

// Error definition for incorrect player ID.
var (
	errIncorrectPlayerId = errors.New("incorrect player_id")
)

// Key type for context value.
type key int

const (
	sessionName     = "testAuth"
	ctxKeyUser  key = iota
)

// apiServer handles API requests.
type apiServer struct {
	router       *mux.Router
	testRouter   *mux.Router // Separate router for test handlers
	sessionStore sessions.Store
	activeUsers  map[string]*User
	logger       *logrus.Logger
}

// newServer creates a new instance of apiServer.
func newServer(sessionStore sessions.Store, logLevel logrus.Level) *apiServer {
	s := &apiServer{
		router:       mux.NewRouter(),
		testRouter:   mux.NewRouter(),
		sessionStore: sessionStore,
		activeUsers:  make(map[string]*User),
		logger:       logrus.New(),
	}

	s.logger.SetLevel(logLevel)
	s.configureRouter()

	s.logger.Info("API server initialized")

	return s
}

// ServeHTTP implements the http.Handler interface.
func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// ServeTestHTTP serves requests for testing.
func (s *apiServer) ServeTestHTTP(w http.ResponseWriter, r *http.Request) {
	s.testRouter.ServeHTTP(w, r)
}

// configureRouter configures the API routes.
func (s *apiServer) configureRouter() {
	s.router.Use(s.UserMiddleware)
	s.router.HandleFunc("/game/create", s.handleCreateGame()).Methods("POST")
	s.router.HandleFunc("/game/attack", s.handleMakeAttack()).Methods("POST")
	s.router.HandleFunc("/game/end_attack", s.handleEndAttack()).Methods("POST")
	s.router.HandleFunc("/game/upgrade", s.handleMakeUpgrade()).Methods("POST")
	s.router.HandleFunc("/game/end_turn", s.handleEndTurn()).Methods("POST")
	s.router.HandleFunc("/game/get_map", s.handleGetMap()).Methods("GET")

	// Add a test handler, used only in tests.
	s.testRouter.Use(s.UserMiddleware)
	s.testRouter.HandleFunc("/test/create_full", s.CreateFullGame()).Methods("POST")
}

// UserMiddleware is a middleware that handles user-related tasks.
func (s *apiServer) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.logger.WithError(err).Error("Error getting session")
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
			s.logger.Info("Creating new user")
			user = NewUser()
			s.activeUsers[userID] = user
		}

		ctx := context.WithValue(r.Context(), ctxKeyUser, user)
		s.logger.WithField("user_id", userID).Debug("Adding user to context")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleCreateGame handles the creation of a new game.
func (s *apiServer) handleCreateGame() http.HandlerFunc {
	type request struct {
		Rows       int `json:"rows"`
		Cols       int `json:"cols"`
		NumPlayers int `json:"num_players"`
		PlayerId   int `json:"player_id"`
		botLevel1  int `json:"bot_level_1"`
		botLevel2  int `json:"bot_level_2"`
		botLevel3  int `json:"bot_level_3"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.WithError(err).Error("Error decoding request")
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.PlayerId < 0 || req.PlayerId >= req.NumPlayers {
			s.logger.WithField("player_id", req.PlayerId).Error("Incorrect player ID")
			s.error(w, r, http.StatusUnprocessableEntity, errIncorrectPlayerId)
			return
		}

		g, err := game.NewGame(req.Rows, req.Cols, req.NumPlayers, 0)

		if err != nil {
			s.logger.WithError(err).Error("Error creating new game")
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*User)
		user.createGame(g, req.PlayerId)

		s.logger.WithFields(logrus.Fields{
			"cols": g.Board.Cols(),
			"rows": g.Board.Rows(),
		}).Info("New game")
		s.respond(w, r, http.StatusCreated, g.ToMap())
	}
}

// handleMakeAttack handles the attack action.
func (s *apiServer) handleMakeAttack() http.HandlerFunc {
	type request struct {
		From struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"from"`
		To struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"to"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.WithError(err).Error("Error decoding request")
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*User)
		err := user.attack(game.Coords(req.From), game.Coords(req.To))

		if err != nil {
			s.logger.WithError(err).Error("Error handling attack")
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.logger.WithFields(logrus.Fields{
			"from": req.From,
			"to":   req.To,
		}).Info("Attack executed")
		s.respond(w, r, http.StatusOK, user.GameBox.Game.ToMap())
	}
}

// handleEndAttack handles ending the attack phase.
func (s *apiServer) handleEndAttack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*User)
		err := user.endAttack()

		if err != nil {
			s.logger.WithError(err).Error("Error ending attack")
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.logger.Info("Attack phase ended")
		s.respond(w, r, http.StatusOK, user.GameBox.Game.ToMap())
	}
}

// handleMakeUpgrade handles the cell upgrade action.
func (s *apiServer) handleMakeUpgrade() http.HandlerFunc {
	type request struct {
		Cell   game.Coords `json:"cell"`
		Levels int         `json:"levels"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.WithError(err).Error("Error decoding request")
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*User)
		err := user.makeUpgrade(req.Cell, req.Levels)

		if err != nil {
			s.logger.WithError(err).Error("Error making upgrade")
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.logger.WithField("cell", req.Cell).Info("Upgrade executed")
		s.respond(w, r, http.StatusOK, user.GameBox.Game.ToMap())
	}
}

// handleEndTurn handles ending the current turn.
func (s *apiServer) handleEndTurn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*User)
		err := user.endTurn()

		if err != nil {
			s.logger.WithError(err).Error("Error ending turn")
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		doAllBotsTurns(user.GameBox.Game, user.GameBox.UserId)

		s.logger.Info("Turn ended")
		s.respond(w, r, http.StatusOK, user.GameBox.Game.ToMap())
	}
}

// handleGetMap handles retrieving the game map.
func (s *apiServer) handleGetMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*User)

		if user.GameBox.Game == nil {
			s.logger.Error("Game does not exist")
			s.error(w, r, http.StatusUnprocessableEntity, errGameIsNotExist)
			return
		}

		s.logger.Info("Retrieved game map")
		s.respond(w, r, http.StatusOK, user.GameBox.Game.ToMap())
	}
}

// CreateFullGame is an endpoint imitation for tests.
func (s *apiServer) CreateFullGame() http.HandlerFunc {
	type request struct {
		Rows       int `json:"rows"`
		Cols       int `json:"cols"`
		NumPlayers int `json:"num_players"`
		PlayerId   int `json:"player_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.PlayerId < 0 || req.PlayerId >= req.NumPlayers {
			s.error(w, r, http.StatusUnprocessableEntity, errIncorrectPlayerId)
			return
		}

		g, err := game.NewCompleteBoardGame(req.Rows, req.Cols, req.NumPlayers)

		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*User)
		user.createGame(g, req.PlayerId)

		s.respond(w, r, http.StatusCreated, g.ToMap())
	}
}

// doAllBotsTurns performs the turns for all AI players.
func doAllBotsTurns(g *game.Game, playerId int) {
	for g.Turn() != playerId && g.Players[playerId].CellsCount() != 0 {
		bot.DoTurn(g, g.Players[g.Turn()])
		bot.DoUpgrade(g, g.Players[g.Turn()])
	}
}

// error responds with an error message.
func (s *apiServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{
		"error": err.Error(),
	})
}

// respond writes a response to the client.
func (s *apiServer) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
