package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gameCreateValidPayload = map[string]int{
	"rows":        3,
	"cols":        3,
	"num_players": 2,
	"player_id":   0,
}

func TestServer_handleGameCreate(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name:         "valid",
			payload:      gameCreateValidPayload,
			expectedCode: http.StatusCreated,
		},
		{
			name: "no player_id",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 2,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "no rows",
			payload: map[string]int{
				"cols":        3,
				"num_players": 2,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "no cols",
			payload: map[string]int{
				"rows":        3,
				"num_players": 2,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "no num_players",
			payload: map[string]int{
				"rows": 3,
				"cols": 3,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "few rows",
			payload: map[string]int{
				"rows":        1,
				"cols":        3,
				"num_players": 2,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "few cols",
			payload: map[string]int{
				"rows":        3,
				"cols":        1,
				"num_players": 2,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "few num_players",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 1,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "many num_players",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 5,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "negative player_id",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 2,
				"player_id":   -1,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "user_id is excessive",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 2,
				"player_id":   4,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/game/create", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

var makeAttackValidPayload = map[string]game.Coords{
	"from": {Row: 0, Col: 0},
	"to":   {Row: 0, Col: 1},
}

func TestServer_handleMakeAttack(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name:         "valid",
			payload:      makeAttackValidPayload,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "negative coords",
			payload: map[string]game.Coords{
				"from": {Row: 0, Col: 0},
				"to":   {Row: 0, Col: -1},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create game and save cookies
			createGameRec := httptest.NewRecorder()
			gameBuf := &bytes.Buffer{}
			json.NewEncoder(gameBuf).Encode(gameCreateValidPayload)
			createGameReq, _ := http.NewRequest(http.MethodPost, "/test/create_full", gameBuf)

			s.ServeTestHTTP(createGameRec, createGameReq)

			cookies := createGameRec.Result().Cookies()

			require.Equal(t, http.StatusCreated, createGameRec.Code)

			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/game/attack", b)

			// Add cookies to request
			for _, c := range cookies {
				req.AddCookie(c)
			}

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// test without creating game
	t.Run("game is not exist", func(t *testing.T) {
		rec := httptest.NewRecorder()

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(testCases[0].payload)

		req, _ := http.NewRequest(http.MethodPost, "/game/attack", b)

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

func TestServer_handleEndAttack(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "valid",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create game and save cookies
			createGameRec := httptest.NewRecorder()
			gameBuf := &bytes.Buffer{}
			json.NewEncoder(gameBuf).Encode(gameCreateValidPayload)
			createGameReq, _ := http.NewRequest(http.MethodPost, "/test/create_full", gameBuf)

			s.ServeTestHTTP(createGameRec, createGameReq)

			cookies := createGameRec.Result().Cookies()

			require.Equal(t, http.StatusCreated, createGameRec.Code)

			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/game/end_attack", nil)

			// Add cookies to request
			for _, c := range cookies {
				req.AddCookie(c)
			}

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// test without creating game
	t.Run("game is not exist", func(t *testing.T) {
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/game/end_attack", nil)

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

func TestServer_handleUpgrade(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"cell":   game.Coords{Row: 0, Col: 0},
				"levels": 1,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid coords",
			payload: map[string]interface{}{
				"cell":   game.Coords{Row: -1, Col: 0},
				"levels": 1,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create game and save cookies
			createGameRec := httptest.NewRecorder()
			gameBuf := &bytes.Buffer{}
			json.NewEncoder(gameBuf).Encode(gameCreateValidPayload)
			createGameReq, _ := http.NewRequest(http.MethodPost, "/test/create_full", gameBuf)

			s.ServeTestHTTP(createGameRec, createGameReq)

			cookies := createGameRec.Result().Cookies()

			require.Equal(t, http.StatusCreated, createGameRec.Code)

			// Call /game/end_attack to make the endpoint valid
			endAttackRec := httptest.NewRecorder()
			endAttackReq, _ := http.NewRequest(http.MethodPost, "/game/end_attack", nil)

			// Add cookies to request
			for _, c := range cookies {
				endAttackReq.AddCookie(c)
			}

			s.ServeHTTP(endAttackRec, endAttackReq)
			require.Equal(t, http.StatusOK, endAttackRec.Code)

			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/game/upgrade", b)

			// Add cookies to request
			for _, c := range cookies {
				req.AddCookie(c)
			}

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// test without creating game
	t.Run("game is not exist", func(t *testing.T) {
		rec := httptest.NewRecorder()

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(testCases[0].payload)

		req, _ := http.NewRequest(http.MethodPost, "/game/upgrade", b)

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

func TestServer_handleEndTurn(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "valid with end_attack",
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid without end_attack",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create game and save cookies
			createGameRec := httptest.NewRecorder()
			gameBuf := &bytes.Buffer{}
			json.NewEncoder(gameBuf).Encode(gameCreateValidPayload)
			createGameReq, _ := http.NewRequest(http.MethodPost, "/test/create_full", gameBuf)

			s.ServeTestHTTP(createGameRec, createGameReq)

			cookies := createGameRec.Result().Cookies()

			require.Equal(t, http.StatusCreated, createGameRec.Code)

			if tc.name == "valid with end_attack" {
				// Call /game/end_attack to make the endpoint valid
				endAttackRec := httptest.NewRecorder()
				endAttackReq, _ := http.NewRequest(http.MethodPost, "/game/end_attack", nil)

				// Add cookies to request
				for _, c := range cookies {
					endAttackReq.AddCookie(c)
				}

				s.ServeHTTP(endAttackRec, endAttackReq)
				require.Equal(t, http.StatusOK, endAttackRec.Code)
			}

			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/game/end_turn", nil)

			// Add cookies to request
			for _, c := range cookies {
				req.AddCookie(c)
			}

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// test without creating game
	t.Run("game is not exist", func(t *testing.T) {
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/game/end_turn", nil)

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

func TestServer_handleGetMap(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "valid with created game",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid without created game",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "valid with created game" {
				// Create game and save cookies
				createGameRec := httptest.NewRecorder()
				gameBuf := &bytes.Buffer{}
				json.NewEncoder(gameBuf).Encode(gameCreateValidPayload)
				createGameReq, _ := http.NewRequest(http.MethodPost, "/test/create_full", gameBuf)

				s.ServeTestHTTP(createGameRec, createGameReq)

				cookies := createGameRec.Result().Cookies()

				require.Equal(t, http.StatusCreated, createGameRec.Code)

				rec := httptest.NewRecorder()

				req, _ := http.NewRequest(http.MethodGet, "/game/get_map", nil)

				// Add cookies to request
				for _, c := range cookies {
					req.AddCookie(c)
				}

				s.ServeHTTP(rec, req)
				assert.Equal(t, tc.expectedCode, rec.Code)
			} else {
				// Test without creating game
				rec := httptest.NewRecorder()

				req, _ := http.NewRequest(http.MethodGet, "/game/get_map", nil)

				s.ServeHTTP(rec, req)
				assert.Equal(t, tc.expectedCode, rec.Code)
			}
		})
	}
}
