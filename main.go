//go:generate swagger generate spec

// Documentation of our Books API
//
//	   Simple Gin API
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 1.0.0
//     Contact: Test User <some_email@example.com> http://github.com/
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/handlers"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/middleware"
	webHandlers "github.com/PauloPortugal/gin-gonic-rest-mongodb/web/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()

	cfg := readConfig()
	router := gin.Default()

	mongoDBClient := setupMongoDBClient(ctx, cfg)
	mongoBooksClient := datastore.NewBooksClient(mongoDBClient, cfg)
	mongoUsersClient := datastore.NewUsersClient(mongoDBClient, cfg)

	redisClient := setupRedisClient(ctx, cfg)
	redisBooksClient := datastore.NewRedisClient(redisClient, cfg)

	mongoBooksClient.InitBooks(ctx)
	mongoUsersClient.InitUsers(ctx)

	booksHandler := handlers.NewBooksHandler(ctx, cfg, mongoBooksClient, redisBooksClient)
	authHandler := handlers.NewAuthHandler(ctx, cfg, mongoUsersClient, redisBooksClient)

	// API public endpoints
	cookieStore, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte(cfg.GetString("redis.sessionSecret")))
	router.Use(sessions.Sessions("books_api_token", cookieStore))
	router.POST("/signin", authHandler.SignIn)
	router.GET("/books", booksHandler.ListBooks)
	router.GET("/books/:id", booksHandler.GetBook)
	router.GET("/books/search", booksHandler.SearchBooks)

	// API private endpoints
	authorised := router.Group("/")
	authorised.Use(middleware.AuthCookieMiddleware())
	authorised.POST("/signout", authHandler.SignOut)
	authorised.POST("/books", booksHandler.NewBook)
	authorised.PUT("/books/:id", booksHandler.UpdateBook)
	authorised.DELETE("/books/:id", booksHandler.DeleteBook)

	// web endpoints
	router.LoadHTMLGlob("web/templates/*")
	router.StaticFile("404.html", "./web/static/404.html")
	webHandler := webHandlers.NewWebHandler(ctx, cfg, mongoBooksClient, redisBooksClient)
	router.Static("/assets", "web/assets")
	router.GET("/", webHandler.IndexPage)
	router.GET("/web/book/:id", webHandler.BookPage)

	corsSetup(router)
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

	err := router.Run()
	if err != nil {
		return
	}
}

func corsSetup(router *gin.Engine) {

}

func readConfig() *viper.Viper {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func setupMongoDBClient(ctx context.Context, cfg *viper.Viper) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/test?authSource=admin",
		cfg.GetString("mongodb.dbuser"),
		cfg.GetString("mongodb.dbpassword"),
		cfg.GetString("MONGODB.DBHOST"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to databse: %w", err))
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(fmt.Errorf("could not ping databse: %w", err))
	}

	log.Println("Connected to MongoDB")

	return client
}

func setupRedisClient(ctx context.Context, cfg *viper.Viper) *redis.Client {
	redisOpt := &redis.Options{
		Addr:     fmt.Sprint(cfg.Get("redis.host")),
		Password: "",
		DB:       0,
	}
	client := redis.NewClient(redisOpt)

	status := client.Ping(ctx)
	log.Print(fmt.Sprintf("redis status: %q", status))

	return client
}
