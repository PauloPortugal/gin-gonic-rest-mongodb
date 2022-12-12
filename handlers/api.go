package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	webHandlers "github.com/PauloPortugal/gin-gonic-rest-mongodb/web/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Setup(ctx context.Context, cfg *viper.Viper, router *gin.Engine, pathToTemplates string, mongoBooksClient datastore.Books,
	mongoUsersClient datastore.Users, redisBooksClient datastore.Redis, m Middleware, cs CookieStore) *gin.Engine {

	booksHandler := NewBooksHandler(ctx, cfg, mongoBooksClient, redisBooksClient)
	authHandler := NewAuthHandler(ctx, cfg, mongoUsersClient, redisBooksClient)

	// API public endpoints
	cookieStore, err := cs.NewCookieStore()
	if err != nil {
		panic(fmt.Errorf("error: %w", err))
	}
	router.Use(sessions.Sessions("books_api_token", cookieStore))
	router.POST("/signin", authHandler.SignIn)
	router.GET("/books", booksHandler.ListBooks)
	router.GET("/books/:id", booksHandler.GetBook)
	router.GET("/books/search", booksHandler.SearchBooks)

	// API private endpoints
	authorised := router.Group("/")
	authorised.Use(m.AuthCookieMiddleware())
	authorised.POST("/signout", authHandler.SignOut)
	authorised.POST("/books", booksHandler.NewBook)
	authorised.PUT("/books/:id", booksHandler.UpdateBook)
	authorised.DELETE("/books/:id", booksHandler.DeleteBook)

	// web endpoints
	router.LoadHTMLGlob(pathToTemplates)
	router.StaticFile("404.html", "./web/static/404.html")
	webHandler := webHandlers.NewWebHandler(ctx, cfg, mongoBooksClient, redisBooksClient)
	router.Static("/assets", "web/assets")
	router.GET("/", webHandler.IndexPage)
	router.GET("/web/book/:id", webHandler.BookPage)

	corsSetup(router)

	return router
}

func corsSetup(router *gin.Engine) {
	// allow swagger UI requests
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
	}))
}
