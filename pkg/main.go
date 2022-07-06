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

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()

	cfg := readConfig()

	setupMongoDB(ctx, cfg)

	router := gin.Default()

	router.POST("/books", handlers.NewBook)
	router.GET("/books", handlers.ListBooks)
	router.PUT("/books/:id", handlers.UpdateBook)
	router.DELETE("/books/:id", handlers.DeleteBook)
	router.GET("/books/search", handlers.SearchBooks)

	err := router.Run()
	if err != nil {
		return
	}
}

func readConfig() *viper.Viper {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func setupMongoDB(ctx context.Context, cfg *viper.Viper) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin",
		cfg.Get("mongodb.dbuser"),
		cfg.Get("mongodb.dbpassword"),
		cfg.Get("mongodb.dbhost"),
		cfg.Get("mongodb.dbname"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(fmt.Errorf("could not connect to databse: %w", err))
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(fmt.Errorf("could not ping databse: %w", err))
	}
	log.Println("Connected to MongoDB")

	return client
}
