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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/handlers"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func setupMongoDB(ctx context.Context, cfg *viper.Viper) *mongo.Client {
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

	setupIndexes(ctx, client, cfg)

	err = loadStaticData(ctx, client, cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("could not insert static data: %w\n", err))
		return nil
	}

	return client
}

func loadStaticData(ctx context.Context, client *mongo.Client, cfg *viper.Viper) error {
	books := make([]model.Book, 0)

	file, err := ioutil.ReadFile("books.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &books); err != nil {
		return err
	}

	db := fmt.Sprint(cfg.Get("mongodb.dbname"))
	col := fmt.Sprint(cfg.Get("mongodb.dbcollection"))

	var b []interface{}
	for _, book := range books {
		b = append(b, book)
	}
	result, err := client.Database(db).Collection(col).InsertMany(ctx, b)
	if err != nil {
		if mongoErr, ok := err.(mongo.BulkWriteException); ok {
			if len(mongoErr.WriteErrors) > 0 && mongoErr.WriteErrors[0].Code == 11000 {
				return nil
			}
		}
		return err
	}

	log.Printf("Inserted books: %d\n", len(result.InsertedIDs))

	return nil
}

func setupIndexes(ctx context.Context, client *mongo.Client, cfg *viper.Viper) {
	db := fmt.Sprint(cfg.Get("mongodb.dbname"))
	col := fmt.Sprint(cfg.Get("mongodb.dbcollection"))

	idxOpt := &options.IndexOptions{}
	idxOpt.SetUnique(true)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		},
		Options: idxOpt,
	}

	ind, err := client.Database(db).Collection(col).Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(fmt.Errorf("Indexes().CreateOne() ERROR: %w", err))
	} else {
		// API call returns string of the index name
		log.Printf("CreateOne() index: %s\n", ind)
	}
}
