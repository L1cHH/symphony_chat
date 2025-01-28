package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	publicDto "symphony_chat/internal/application/dto"
	"symphony_chat/internal/domain/users"
	authhttp "symphony_chat/tests/integration/auth/http"
	"symphony_chat/tests/integration/setup"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSignUpHandler(t *testing.T) {
	db, err := setup.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	router := authhttp.SetupRouter(t, db)

	testCases := []struct {
		name string
		credentials publicDto.LoginCredentials
		expectedCode int
		expectedErr string
		beforeTestAction func(t *testing.T, testDB *setup.TestDB)
	}{
		{
			name: "Success registration",
			credentials: publicDto.LoginCredentials {
				Login: "Bomj.Obichnyi@gmail.com",
				Password: "eptakaktakto228",
			},
			expectedCode: http.StatusOK,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
		{
			name: "Duplicate User",
			credentials: publicDto.LoginCredentials {
				Login: "Alexey.Gnida2003.gmail.com",
				Password: "hesusOneLover335",
			},
			expectedCode: http.StatusConflict,
			expectedErr: users.ErrLoginAlreadyExists.Code,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				body := publicDto.LoginCredentials {
					Login: "Alexey.Gnida2003.gmail.com",
					Password: "hesusOneLover335",
				}

				jsonBody, err := json.Marshal(body)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTestAction != nil {
				tc.beforeTestAction(t, db)
			}

			jsonBody, err := json.Marshal(tc.credentials)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)


			if tc.expectedErr != "" {
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Contains(t, response["code"], tc.expectedErr)
			} else {
				var response map[string]publicDto.JWTAccessTokenDTO
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response["access_token"])
				require.NotEmpty(t, response["refresh_token"])

				cookies := w.Result().Cookies()
				var foundRefreshToken bool 
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						foundRefreshToken = true
						break
					}
				}

				require.True(t, foundRefreshToken, "refresh token not found in cookies")
			}
		})
	}
}