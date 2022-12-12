package handlers

import (
	"context"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//go:generate moq -out interfaces_books_moq_test.go . books
type books interface {
	AddBook(ctx context.Context, book *model.Book) error
	ListBooks(ctx context.Context) ([]model.Book, error)
	SearchBooks(ctx context.Context, tag string) ([]model.Book, error)
	GetBook(ctx context.Context, id string) (model.Book, error)
	UpdateBook(ctx context.Context, id string, book model.Book) (int, error)
	DeleteBook(ctx context.Context, id string) (int, error)
}

//go:generate moq -out interfaces_redis_moq_test.go . redis
type redis interface {
	GetBooks(ctx context.Context) ([]model.Book, error)
	SetBooks(ctx context.Context, books []model.Book) error
	GetBook(ctx context.Context, id string) (model.Book, error)
	SetBook(ctx context.Context, id string, book model.Book) error
	DeleteEntry(ctx context.Context, id string)
}

//go:generate moq -out interfaces_users_moq_test.go . users
type users interface {
	Get(ctx context.Context, username string, password string) (model.User, error)
}

//go:generate moq -out interfaces_middleware_moq_test.go . Middleware
type Middleware interface {
	AuthJWTMiddleware(cfg *viper.Viper) gin.HandlerFunc
	AuthCookieMiddleware() gin.HandlerFunc
}

//go:generate moq -out interfaces_cookie_store_moq_test.go . CookieStore
type CookieStore interface {
	NewCookieStore() (redisStore.Store, error)
}
