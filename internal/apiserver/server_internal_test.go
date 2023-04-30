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

func TestServer_handleGameCreate(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]int{
				"rows":        3,
				"cols":        3,
				"num_players": 2,
				"player_id":   0,
			},
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

func TestServer_handleMakeAttack(t *testing.T) {
	s := newServer(sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]game.Coords{
				"from": {Row: 0, Col: 0},
				"to":   {Row: 0, Col: 1},
			},
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

	createGamePayload := map[string]int{
		"rows":        3,
		"cols":        3,
		"num_players": 2,
		"player_id":   0,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create game and save cookies
			createGameRec := httptest.NewRecorder()
			gameBuf := &bytes.Buffer{}
			json.NewEncoder(gameBuf).Encode(createGamePayload)
			createGameReq, _ := http.NewRequest(http.MethodPost, "/game/create", gameBuf)

			s.ServeHTTP(createGameRec, createGameReq)

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
