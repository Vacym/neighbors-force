package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/game/create", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

// func TestServer_handleSessionsCreate(t *testing.T) {
// 	s := newServer(sessions.NewCookieStore([]byte("secret")))

// 	testCases := []struct {
// 		name         string
// 		payload      any
// 		expectedCode int
// 	}{
// 		{
// 			name: "valid",
// 			payload: map[string]string{
// 				"email":    u.Email,
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusOK,
// 		},
// 		{
// 			name:         "invalid payload",
// 			payload:      "invalid",
// 			expectedCode: http.StatusBadRequest,
// 		},
// 		{
// 			name: "invalid email",
// 			payload: map[string]string{
// 				"email":    "test@mail.org",
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusUnauthorized,
// 		},
// 		{
// 			name: "invalid password",
// 			payload: map[string]string{
// 				"email":    u.Email,
// 				"password": "invalid",
// 			},
// 			expectedCode: http.StatusUnauthorized,
// 		},
// 		{
// 			name: "invalid both",
// 			payload: map[string]string{
// 				"email":    "test@mail.org",
// 				"password": "invalid",
// 			},
// 			expectedCode: http.StatusUnauthorized,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rec := httptest.NewRecorder()

// 			b := &bytes.Buffer{}
// 			json.NewEncoder(b).Encode(tc.payload)

// 			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)

// 			s.ServeHTTP(rec, req)
// 			assert.Equal(t, tc.expectedCode, rec.Code)
// 		})
// 	}
// }
