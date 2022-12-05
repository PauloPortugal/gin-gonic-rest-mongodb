package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const pathToTemplates = "../web/templates/*"

func createRouter(booksClientMock *booksMock, usersClientMock *usersMock, redisClientMock *redisMock, m Middleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := Setup(context.Background(), viper.GetViper(), gin.New(), pathToTemplates,
		booksClientMock, usersClientMock, redisClientMock, m)
	return router
}

func getBook() model.Book {
	return model.Book{
		ID:        primitive.ObjectID{},
		Name:      "Moondust",
		Author:    "Andrew Smith",
		Publisher: "Bloomsbury Publishing PLC",
		PublishedAt: model.PublishedDate{
			Month: "July",
			Year:  "2009",
		},
		Tags:           []string{"space exploration", "astronauts", "nasa"},
		ImagePath:      "/assets/images/Moondust.jpg",
		Review:         4.6,
		SubmissionDate: time.Now(),
	}
}

func mockMiddleware() *MiddlewareMock {
	return &MiddlewareMock{AuthCookieMiddlewareFunc: func() gin.HandlerFunc {
		return func(ctx *gin.Context) { ctx.Next() }
	}}
}

func TestBooksHandler_ListBooks(t *testing.T) {
	var books []model.Book
	var errorRes map[string]string

	Convey("Given GET /books", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/books", nil)

		Convey("When redis DB has only one Book in cache", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{
					GetBooksFunc: func(ctx context.Context) ([]model.Book, error) {
						return []model.Book{getBook()}, nil
					},
				},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &books)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And only one book is returned", func() {
				So(len(books), ShouldEqual, 1)
			})

			Convey("And it is the expected book", func() {
				So(books[0].Name, ShouldEqual, getBook().Name)
			})
		})

		Convey("When redis DB returns an error", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{
					GetBooksFunc: func(ctx context.Context) ([]model.Book, error) {
						return nil, fmt.Errorf("some RedisDB error")
					},
				},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

			Convey("Then 500 Internal Server Error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(errorRes["error"], ShouldEqual, "some RedisDB error")
			})
		})

		Convey("When redis DB does not have any Books in cache", func() {
			rmock := &redisMock{
				GetBooksFunc: func(ctx context.Context) ([]model.Book, error) {
					return nil, nil
				},
			}
			Convey("And a Book exists on MongoDB", func() {
				bmock := &booksMock{ListBooksFunc: func(ctx context.Context) ([]model.Book, error) {
					return []model.Book{getBook()}, nil
				}}
				Convey("And Redis SetBooks is successful", func() {
					rmock.SetBooksFunc = func(ctx context.Context, books []model.Book) error { return nil }
					router := createRouter(
						bmock,
						&usersMock{},
						rmock,
						mockMiddleware(),
					)

					router.ServeHTTP(w, req)
					_ = json.Unmarshal(w.Body.Bytes(), &books)

					Convey("Then 200 OK is returned", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})

					Convey("And only one book is returned", func() {
						So(len(books), ShouldEqual, 1)
					})

					Convey("And it is the expected book", func() {
						So(books[0].Name, ShouldEqual, getBook().Name)
					})
				})

				Convey("And Redis SetBooks is unsuccessful", func() {
					rmock.SetBooksFunc = func(ctx context.Context, books []model.Book) error { return fmt.Errorf("some RedisDB SetBooks error") }
					router := createRouter(
						bmock,
						&usersMock{},
						rmock,
						mockMiddleware(),
					)

					router.ServeHTTP(w, req)
					_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

					Convey("Then 500 Internal Server Error is returned", func() {
						So(w.Code, ShouldEqual, http.StatusInternalServerError)
					})

					Convey("And response body has the expected error message", func() {
						So(errorRes["error"], ShouldEqual, "some RedisDB SetBooks error")
					})
				})
			})

			Convey("And MongoDB returns an error", func() {
				router := createRouter(
					&booksMock{ListBooksFunc: func(ctx context.Context) ([]model.Book, error) {
						return nil, fmt.Errorf("some MongoDB error")
					}},
					&usersMock{},
					rmock,
					mockMiddleware(),
				)

				router.ServeHTTP(w, req)
				_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

				Convey("Then 500 Internal Server Error is returned", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})

				Convey("And response body has the expected error message", func() {
					So(errorRes["error"], ShouldEqual, "some MongoDB error")
				})
			})
		})
	})
}

func TestBooksHandler_SearchBooks(t *testing.T) {
	var books []model.Book
	var errorRes map[string]string

	Convey("Given GET /books/search?tag=nasa", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/books/search?tag=nasa", nil)

		Convey("When MongoDB DB has only one Book", func() {
			router := createRouter(
				&booksMock{
					SearchBooksFunc: func(ctx context.Context, tag string) ([]model.Book, error) {
						return []model.Book{getBook()}, nil
					},
				},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &books)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And only one book is returned", func() {
				So(len(books), ShouldEqual, 1)
			})

			Convey("And it is the expected book", func() {
				So(books[0].Name, ShouldEqual, getBook().Name)
			})
		})

		Convey("And MongoDB returns an error", func() {
			router := createRouter(
				&booksMock{SearchBooksFunc: func(ctx context.Context, tag string) ([]model.Book, error) {
					return nil, fmt.Errorf("some MongoDB error")
				}},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

			Convey("Then 500 Internal Server Error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(errorRes["error"], ShouldEqual, "some MongoDB error")
			})
		})
	})
}

func TestBooksHandler_GetBook(t *testing.T) {
	var book model.Book
	var errorRes map[string]string

	Convey("Given GET /books/{id}", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/books/some-id", nil)

		Convey("When redis DB has the book in cache", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{
					GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
						return getBook(), nil
					},
				},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &book)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And it is the expected book", func() {
				So(book.Name, ShouldEqual, getBook().Name)
			})
		})

		Convey("When redis DB returns an error", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{
					GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
						return model.Book{}, fmt.Errorf("some RedisDB error")
					},
				},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

			Convey("Then 500 Internal Server Error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(errorRes["error"], ShouldEqual, "some RedisDB error")
			})
		})

		Convey("When redis DB does not have the book in cache", func() {
			rmock := &redisMock{
				GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
					return model.Book{}, nil
				},
			}

			Convey("And the book exists on MongoDB", func() {
				bmock := &booksMock{GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
					return getBook(), nil
				}}

				Convey("And Redis SetBook is successful", func() {
					rmock.SetBookFunc = func(ctx context.Context, id string, book model.Book) error { return nil }
					router := createRouter(
						bmock,
						&usersMock{},
						rmock,
						mockMiddleware(),
					)

					router.ServeHTTP(w, req)
					_ = json.Unmarshal(w.Body.Bytes(), &book)

					Convey("Then 200 OK is returned", func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})

					Convey("And it is the expected book", func() {
						So(book.Name, ShouldEqual, getBook().Name)
					})
				})

				Convey("And Redis SetBooks is unsuccessful", func() {
					rmock.SetBookFunc = func(ctx context.Context, id string, book model.Book) error {
						return fmt.Errorf("some RedisDB SetBook error")
					}
					router := createRouter(
						bmock,
						&usersMock{},
						rmock,
						mockMiddleware(),
					)

					router.ServeHTTP(w, req)
					_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

					Convey("Then 500 Internal Server Error is returned", func() {
						So(w.Code, ShouldEqual, http.StatusInternalServerError)
					})

					Convey("And response body has the expected error message", func() {
						So(errorRes["error"], ShouldEqual, "some RedisDB SetBook error")
					})
				})
			})

			Convey("And MongoDB returns an error", func() {
				router := createRouter(
					&booksMock{GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
						return model.Book{}, fmt.Errorf("some MongoDB error")
					}},
					&usersMock{},
					rmock,
					mockMiddleware(),
				)

				router.ServeHTTP(w, req)
				_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

				Convey("Then 500 Internal Server Error is returned", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})

				Convey("And response body has the expected error message", func() {
					So(errorRes["error"], ShouldEqual, "some MongoDB error")
				})
			})

			Convey("And MongoDB does not have the book", func() {
				router := createRouter(
					&booksMock{GetBookFunc: func(ctx context.Context, id string) (model.Book, error) {
						return model.Book{}, nil
					}},
					&usersMock{},
					rmock,
					mockMiddleware(),
				)

				router.ServeHTTP(w, req)
				_ = json.Unmarshal(w.Body.Bytes(), &errorRes)

				Convey("Then 404 Not Found is returned", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
				})
			})
		})
	})
}

func TestBooksHandler_DeleteBook(t *testing.T) {
	var bodyMsg map[string]string

	Convey("Given DELETE /books/{id}", t, func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/books/some-id", nil)

		Convey("When MongoDB has book and successfully deletes the book", func() {
			router := createRouter(
				&booksMock{DeleteBookFunc: func(ctx context.Context, id string) (int, error) { return 1, nil }},
				&usersMock{},
				&redisMock{DeleteEntryFunc: func(ctx context.Context, id string) {}},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And a confirmation message that the book was deleted is returned", func() {
				So(bodyMsg["message"], ShouldEqual, "Book has been deleted")
			})
		})

		Convey("When MongoDB does not have the book", func() {
			router := createRouter(
				&booksMock{DeleteBookFunc: func(ctx context.Context, id string) (int, error) { return 0, nil }},
				&usersMock{},
				&redisMock{DeleteEntryFunc: func(ctx context.Context, id string) {}},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 404 Not Found is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})

			Convey("And a confirmation message that the book was not found is returned", func() {
				So(bodyMsg["error"], ShouldEqual, "book not found")
			})
		})

		Convey("When MongoDB returns an error", func() {
			router := createRouter(
				&booksMock{DeleteBookFunc: func(ctx context.Context, id string) (int, error) { return 0, fmt.Errorf("some MongoDB error") }},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 404 Not Found is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(bodyMsg["error"], ShouldEqual, "some MongoDB error")
			})
		})
	})
}

func TestBooksHandler_UpdateBook(t *testing.T) {
	var bodyMsg map[string]string
	var book model.Book

	Convey("Given PUT /books/{id}", t, func() {
		w := httptest.NewRecorder()
		data, _ := json.Marshal(getBook())

		req, _ := http.NewRequest(http.MethodPut, "/books/some-id", bytes.NewReader(data))

		Convey("When MongoDB successfully updates the book", func() {
			router := createRouter(
				&booksMock{UpdateBookFunc: func(ctx context.Context, id string, book model.Book) (int, error) {
					return 1, nil
				}},
				&usersMock{},
				&redisMock{DeleteEntryFunc: func(ctx context.Context, id string) {}},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &book)

			Convey("Then 200 OK is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("And it is the expected book", func() {
				So(book.Name, ShouldEqual, getBook().Name)
			})
		})

		Convey("When request payload is incorrect", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)
			data, _ := json.Marshal("incorrect payload")
			req, _ := http.NewRequest(http.MethodPut, "/books/some-id", bytes.NewReader(data))
			router.ServeHTTP(w, req)

			Convey("Then 400 Bad Request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When MongoDB returns an error", func() {
			router := createRouter(
				&booksMock{UpdateBookFunc: func(ctx context.Context, id string, book model.Book) (int, error) {
					return 0, fmt.Errorf("some MongoDB error")
				}},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 500 Internal Server Error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(bodyMsg["error"], ShouldEqual, "some MongoDB error")
			})
		})

		Convey("When MongoDB does not have the book entry", func() {
			router := createRouter(
				&booksMock{UpdateBookFunc: func(ctx context.Context, id string, book model.Book) (int, error) {
					return 0, nil
				}},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 404 Not Found is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})

			Convey("And response body has the expected error message", func() {
				So(bodyMsg["error"], ShouldEqual, "book not found")
			})
		})
	})
}

func TestBooksHandler_NewBook(t *testing.T) {
	var bodyMsg map[string]string
	var book model.Book

	Convey("Given POST /books", t, func() {
		data, _ := json.Marshal(getBook())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewReader(data))

		Convey("When MongoDB successfully creates book", func() {
			router := createRouter(
				&booksMock{AddBookFunc: func(ctx context.Context, book *model.Book) error {
					return nil
				}},
				&usersMock{},
				&redisMock{DeleteEntryFunc: func(ctx context.Context, id string) {}},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &book)

			Convey("Then 201 Created is returned", func() {
				So(w.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("And it is the expected book", func() {
				So(book.Name, ShouldEqual, getBook().Name)
			})
		})

		Convey("When request payload is incorrect", func() {
			router := createRouter(
				&booksMock{},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)
			data, _ := json.Marshal("incorrect payload")
			req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewReader(data))
			router.ServeHTTP(w, req)

			Convey("Then 400 Bad Request is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When MongoDB returns an error", func() {
			router := createRouter(
				&booksMock{AddBookFunc: func(ctx context.Context, book *model.Book) error {
					return fmt.Errorf("some MongoDB error")
				}},
				&usersMock{},
				&redisMock{},
				mockMiddleware(),
			)

			router.ServeHTTP(w, req)
			_ = json.Unmarshal(w.Body.Bytes(), &bodyMsg)

			Convey("Then 500 Internal Server Error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And response body has the expected error message", func() {
				So(bodyMsg["error"], ShouldEqual, "some MongoDB error")
			})
		})
	})
}
