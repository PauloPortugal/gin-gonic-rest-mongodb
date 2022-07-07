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

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/datastore"
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
	client := setupMongoDBClient(ctx, cfg)

	store := datastore.New(client, cfg)
	store.Init(ctx)

	router := gin.Default()

	booksHandler := handlers.New(ctx, cfg, store)

	router.POST("/books", booksHandler.NewBook)
	router.GET("/books", booksHandler.ListBooks)
	router.PUT("/books/:id", booksHandler.UpdateBook)
	router.DELETE("/books/:id", booksHandler.DeleteBook)
	router.GET("/books/search", booksHandler.SearchBooks)

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
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func setupMongoDBClient(ctx context.Context, cfg *viper.Viper) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/test?authSource=admin",
		cfg.Get("mongodb.dbuser"),
		cfg.Get("mongodb.dbpassword"),
		cfg.Get("mongodb.dbhost"))

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
