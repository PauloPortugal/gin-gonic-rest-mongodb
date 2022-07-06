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

	client := setupMongoDBClient(ctx, cfg)
	booksCollection := setBooksCollection(cfg, client)

	setupIndexes(ctx, booksCollection, cfg)

	if err := loadStaticData(ctx, booksCollection); err != nil {
		log.Fatal(fmt.Errorf("could not insert static data: %w\n", err))
	}

	router := gin.Default()
	booksHandler := handlers.New(ctx, cfg, booksCollection)

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

func setBooksCollection(cfg *viper.Viper, client *mongo.Client) *mongo.Collection {
	db := fmt.Sprint(cfg.Get("mongodb.dbname"))
	col := fmt.Sprint(cfg.Get("mongodb.dbcollection"))

	return client.Database(db).Collection(col)
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

func loadStaticData(ctx context.Context, collection *mongo.Collection) error {
	books := make([]model.Book, 0)

	file, err := ioutil.ReadFile("books.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &books); err != nil {
		return err
	}

	var b []interface{}
	for _, book := range books {
		b = append(b, book)
	}
	result, err := collection.InsertMany(ctx, b)
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

func setupIndexes(ctx context.Context, collection *mongo.Collection, cfg *viper.Viper) {
	idxOpt := &options.IndexOptions{}
	idxOpt.SetUnique(true)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		},
		Options: idxOpt,
	}

	ind, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(fmt.Errorf("Indexes().CreateOne() ERROR: %w", err))
	} else {
		// BooksHandler call returns string of the index name
		log.Printf("CreateOne() index: %s\n", ind)
	}
}
