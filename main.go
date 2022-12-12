//go:generate swagger generate spec

// Documentation of our Books API
//
//		   Simple Gin API
//
//	    Schemes: http
//	    Host: localhost:8080
//	    BasePath: /
//	    Version: 1.0.0
//	    Contact: Test User <some_email@example.com> http://github.com/
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Security:
//	    - basic
//
//	   SecurityDefinitions:
//	   basic:
//	     type: basic
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/handlers"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const pathToTemplates = "web/templates/*"

func main() {
	ctx := context.Background()
	cfg := readConfig()

	mongoDBClient := setupMongoDBClient(ctx, cfg)
	mongoBooksClient := datastore.NewBooksClient(mongoDBClient, cfg)
	mongoUsersClient := datastore.NewUsersClient(mongoDBClient, cfg)
	mongoBooksClient.InitBooks(ctx)
	mongoUsersClient.InitUsers(ctx)

	redisClient := setupRedisClient(ctx, cfg)
	redisBooksClient := datastore.NewRedisClient(redisClient, cfg)

	m := middleware.Client{}
	cs := datastore.CookieStoreClient{Config: cfg}

	router := gin.Default()

	if err := handlers.Setup(ctx, cfg, router, pathToTemplates, mongoBooksClient, mongoUsersClient, redisBooksClient, m, cs).Run(); err != nil {
		return
	}
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
		cfg.GetString("mongodb.dbhost"))

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
	log.Printf("redis status: %q", status)

	return client
}
