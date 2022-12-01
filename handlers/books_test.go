package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/PauloPortugal/gin-gonic-rest-mongodb"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBooksHandler_ListBooks(t *testing.T) {
	router := Setup(context.Background(), readConfig(), nil, nil, nil, nil)

	Convey("Given I have GET /books", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/books", nil)
		router.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}
