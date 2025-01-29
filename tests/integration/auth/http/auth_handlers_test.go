package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	publicDto "symphony_chat/internal/application/dto"
	"symphony_chat/internal/application/middleware"
	"symphony_chat/internal/domain/users"
	authhttp "symphony_chat/tests/integration/auth/http"
	"symphony_chat/tests/integration/setup"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		expectedHttpCode int
		expectedErrCode string
		beforeTestAction func(t *testing.T, testDB *setup.TestDB)
	}{
		{
			name: "Success registration",
			credentials: publicDto.LoginCredentials {
				Login: "Bomj.Obichnyi@gmail.com",
				Password: "eptakaktakto228",
			},
			expectedHttpCode: http.StatusOK,
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
			expectedHttpCode: http.StatusConflict,
			expectedErrCode: users.ErrLoginAlreadyExists.Code,
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
		{
			name: "Invalid login format (short Login)",
			credentials: publicDto.LoginCredentials{
				Login: "gmail",
				Password: "kolomin.andrey2005",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrCode: "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid login format (long Login)",
			credentials: publicDto.LoginCredentials {
				Login: "cnqibqvqgbnqognbqbqiomvevnevevne.gmail.com",
				Password: "kolomin.andrey2005",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrCode: "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid login format (disallowed symbols)",
			credentials: publicDto.LoginCredentials {
				Login: "cbevvbev%^gmail.com",
				Password: "kolomin.andrey2005",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrCode: "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid password format (short Password)",
			credentials: publicDto.LoginCredentials {
				Login: "Kolomin.Andrey@gmail.com",
				Password: "12345",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrCode: "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid password format (long Password)",
			credentials: publicDto.LoginCredentials {
				Login: "Kolomin.Andrey@gmail.com",
				Password: "kolomin.andrey2.andrey2005",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedErrCode: "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
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

			require.Equal(t, tc.expectedHttpCode, w.Code)


			if tc.expectedErrCode != "" {
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Contains(t, response["code"], tc.expectedErrCode)
			} else {
				var response map[string]publicDto.JWTTokenDTO
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

func TestLogInHandler(t *testing.T) {
	testDB, err := setup.NewTestDB()
	require.NoError(t, err)

	defer func() {
		err := testDB.Close()
		require.NoError(t, err)
	}()

	router := authhttp.SetupRouter(t, testDB)

	testCases := []struct {
		name              string
		credentials       publicDto.LoginCredentials
		expectedHttpCode  int
		expectedErrCode   string
		beforeTestAction  func(t *testing.T, testDB *setup.TestDB)
	}{
		{
			name: "Success login",
			credentials: publicDto.LoginCredentials {
				Login:    "Andrei.Karpukh2000@gmail.com",
				Password: "andrei_kriper2004boi",
			},
			expectedHttpCode: http.StatusOK,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				bodyJson, err := json.Marshal(publicDto.LoginCredentials {
					Login:    "Andrei.Karpukh2000@gmail.com",
					Password: "andrei_kriper2004boi",
				})

				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(bodyJson))
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)
				require.Equal(t, res.Code, http.StatusOK)
			},
		},

		{
			name: "Wrong password",
			credentials: publicDto.LoginCredentials {
				Login:    "Andrei.Karpukh2000@gmail.com",
				Password: "andrei_kriper2004boi1",
			},
			expectedHttpCode: http.StatusUnauthorized,
			expectedErrCode:  "WRONG_PASSWORD",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				bodyJson, err := json.Marshal(publicDto.LoginCredentials {
					Login:    "Andrei.Karpukh2000@gmail.com",
					Password: "andrei_kriper2004boi",
				})

				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(bodyJson))
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)
				require.Equal(t, res.Code, http.StatusOK)
			},
		},

		{
			name: "User not found",
			credentials: publicDto.LoginCredentials {
				Login:    "Andrei.Karpukh2000@gmail.com",
				Password: "andrei_kriper2004boi1",
			},
			expectedHttpCode: http.StatusNotFound,
			expectedErrCode:  "AUTH_USER_WITH_THIS_LOGIN_NOT_FOUND",
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.beforeTestAction(t, testDB)

			bodyJson, err := json.Marshal(tc.credentials)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(bodyJson))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			require.Equal(t, tc.expectedHttpCode, res.Code)

			if tc.expectedErrCode != "" {
				var response map[string]string

				err := json.Unmarshal(res.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, tc.expectedErrCode, response["code"])
				require.Empty(t, response["access_token"])
				require.Empty(t, response["refresh_token"])

				var foundRefreshTokenInCookies bool

				cookies := res.Result().Cookies()

				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						foundRefreshTokenInCookies = true
						break
					}
				}

				require.False(t, foundRefreshTokenInCookies)

			} else {
				var response map[string]publicDto.JWTTokenDTO
				err := json.Unmarshal(res.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response["access_token"])
				require.NotEmpty(t, response["refresh_token"])

				var foundRefreshTokenInCookies bool

				cookies := res.Result().Cookies()

				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						foundRefreshTokenInCookies = true
						break
					}
				}

				require.True(t, foundRefreshTokenInCookies)
			}
		})
	}
}

func TestLogOutHandler(t *testing.T) {
	testDB, err := setup.NewTestDB()
	require.NoError(t, err)

	defer func() {
		err := testDB.Close()
		require.NoError(t, err)
	}()

	router := authhttp.SetupRouter(t, testDB)

	testCases := []struct {
		name   				string
		expectedHttpCode 	int
		expectedErrCode 	string
		beforeTestAction 	func(t *testing.T, testDB *setup.TestDB) (*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name: "Success Logout",
			expectedHttpCode: http.StatusOK,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) (*httptest.ResponseRecorder, *http.Request) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				jsonBody, err := json.Marshal(publicDto.LoginCredentials {
					Login:    "Andrei.Karpukh2000@gmail.com",
					Password: "fhigbgiwgwwhnwihwgwb",
				})

				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)
				require.Equal(t, res.Code, http.StatusOK)

				var foundRefreshTokenInCookies bool

				cookies := res.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						foundRefreshTokenInCookies = true
						break
					}
				}

				require.True(t, foundRefreshTokenInCookies)

				var responseBody map[string]publicDto.JWTTokenDTO

				err = json.Unmarshal(res.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.NotEmpty(t, responseBody["access_token"])
				require.NotEmpty(t, responseBody["refresh_token"])

				req = httptest.NewRequest("POST", "/auth/logout", nil)

				req.Header.Set("Authorization", "Bearer "+ responseBody["access_token"].Token)

				res = httptest.NewRecorder()

				return res, req
			},
		},

		{
			name: "Empty Authorization header",
			expectedHttpCode: http.StatusUnauthorized,
			expectedErrCode: middleware.ErrAccessTokenWasNotProvided.Code,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) (*httptest.ResponseRecorder, *http.Request) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)


				req := httptest.NewRequest("POST", "/auth/logout", nil)
				res := httptest.NewRecorder()

				return res, req
			},
		},
		{
			name: "Invalid Authorization header",
			expectedHttpCode: http.StatusUnauthorized,
			expectedErrCode: middleware.ErrInvalidAccessTokenFormat.Code,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) (*httptest.ResponseRecorder, *http.Request) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/logout", nil)
				res := httptest.NewRecorder()

				req.Header.Set("Authorization", "invalid authorization header")

				return res, req
			},
		},

		{
			name: "Expired access token",
			expectedHttpCode: http.StatusOK,
			beforeTestAction: func(t *testing.T, testDB *setup.TestDB) (*httptest.ResponseRecorder, *http.Request) {
				err := testDB.TruncateAllTables()
				require.NoError(t, err)

				jsonBody, err := json.Marshal(publicDto.LoginCredentials {
					Login:    "Andrei.Karpukh2000@gmail.com",
					Password: "fhigbgiwgwwhnwihwgwb",
				})

				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)
				require.Equal(t, res.Code, http.StatusOK)

				var foundRefreshTokenInCookies bool

				cookies := res.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						foundRefreshTokenInCookies = true
						break
					}
				}

				require.True(t, foundRefreshTokenInCookies)

				var responseBody map[string]publicDto.JWTTokenDTO

				err = json.Unmarshal(res.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.NotEmpty(t, responseBody["access_token"])
				require.NotEmpty(t, responseBody["refresh_token"])

				req = httptest.NewRequest("POST", "/auth/logout", nil)

				duration := os.Getenv("ACCESS_TTL_IN_MINUTES")
				require.NotEmpty(t, duration)
				durationInt, err := strconv.Atoi(duration)
				require.NoError(t, err)

				time.Sleep(time.Duration(durationInt) *time.Minute)

				req.Header.Set("Authorization", "Bearer "+ responseBody["access_token"].Token)
				
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: responseBody["refresh_token"].Token,
				})

				res = httptest.NewRecorder()

				return res, req

			},
		},
	}
	

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, req := tc.beforeTestAction(t, testDB)

			router.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedHttpCode, res.Code)

			if tc.expectedErrCode != "" {
				var response map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedErrCode, response["code"])
			} else {
				var responseBody map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				assert.Equal(t, "LOGOUT_SUCCESSFUL", responseBody["code"])
				assert.Equal(t, "clear_tokens", responseBody["message"])

				if tc.name == "Expired access token" {
					require.NotEmpty(t, res.Header().Get("New-Access-Token"))
				}

				var IsCookieContainsInvalidatedRefreshToken bool

				cookies := res.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" && cookie.Value == "" {
						IsCookieContainsInvalidatedRefreshToken = true
						break
					}
				}

				assert.True(t, IsCookieContainsInvalidatedRefreshToken)
			}

		})
	}
}