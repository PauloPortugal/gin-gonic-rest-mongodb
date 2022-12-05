package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthHandler_SignIn(t *testing.T) {
	var bodyMsg map[string]string

	Convey("Given POST /signin", t, func() {
		w := httptest.NewRecorder()

		data, _ := json.Marshal(map[string]string{"username": "username", "password": "password"})
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(data))

		Convey("When MongoDB has the user record", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{GetFunc: func(ctx context.Context, username string, password string) (model.User, error) {
					return model.User{Username: username, Password: password}, nil
				}},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And it is the expected book", func() {
				So(bodyMsg["message"], ShouldEqual, "User signed in")
			})
		})

		Convey("When login details are incorrect", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{GetFunc: func(ctx context.Context, username string, password string) (model.User, error) {
					return model.User{}, fmt.Errorf("user not found")
				}},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 401 Unauthorized is returned", func() {
				So(w.Code, ShouldEqual, http.StatusUnauthorized)
			})

			Convey("And response body has the expected error message", func() {
				So(bodyMsg["error"], ShouldEqual, "Invalid username or password")
			})
		})

		Convey("When login payload is incorrect", func() {
			data, _ := json.Marshal("incorrect payload")
			req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(data))
			router := createRouter(
				&booksMock{},
				&usersMock{GetFunc: func(ctx context.Context, username string, password string) (model.User, error) {
					return model.User{}, fmt.Errorf("user not found")
				}},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 400 Bad Request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func TestAuthHandler_SignOut(t *testing.T) {
	var bodyMsg map[string]string

	Convey("Given POST /signout", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/signout", nil)

		Convey("When signout is successful", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And it is the expected book", func() {
				So(bodyMsg["message"], ShouldEqual, "User signed out")
			})
		})
	})
}
