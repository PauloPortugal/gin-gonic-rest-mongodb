package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Redis interface {
	GetBooks(ctx context.Context) ([]model.Book, error)
	SetBooks(ctx context.Context, books []model.Book) error
	GetBook(ctx context.Context, id string) (model.Book, error)
	SetBook(ctx context.Context, id string, book model.Book) error
	DeleteEntry(ctx context.Context, id string)
}

// RedisClient is the client responsible for querying redis
type RedisClient struct {
	client *redis.Client
	cfg    *viper.Viper
}

func NewRedisClient(client *redis.Client, cfg *viper.Viper) *RedisClient {
	return &RedisClient{
		client: client,
		cfg:    cfg,
	}
}

// GetBook wrapper to return a cached book from Redis
func (c *RedisClient) GetBook(ctx context.Context, id string) (model.Book, error) {
	var book model.Book

	result, err := c.client.Get(ctx, id).Result()
	if err == redis.Nil {
		return model.Book{}, nil
	} else if err != nil {
		log.Print(fmt.Errorf("could not get all books from redis: %w", err))
		return model.Book{}, err
	}

	if err = json.Unmarshal([]byte(result), &book); err != nil {
		log.Print(fmt.Errorf("unmarshal error from redis result: %w", err))
		return model.Book{}, err
	}

	return book, nil
}

// GetBooks wrapper to return all books from Redis
func (c *RedisClient) GetBooks(ctx context.Context) ([]model.Book, error) {
	books := make([]model.Book, 0)

	result, err := c.client.Get(ctx, "books").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Print(fmt.Errorf("could not get all books from redis: %w", err))
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &books); err != nil {
		log.Print(fmt.Errorf("unmarshal error from redis result: %w", err))
		return nil, err
	}

	return books, nil
}

// SetBooks wrapper to set the current list of Books on Redis
func (c *RedisClient) SetBooks(ctx context.Context, books []model.Book) error {
	data, err := json.Marshal(books)
	if err != nil {
		log.Print(fmt.Errorf("could not marshal books to redis: %w", err))
		return err
	}

	set := c.client.Set(ctx, "books", string(data), time.Hour)
	log.Print(fmt.Errorf("redis set status: %v", set))
	return nil
}

// SetBook wrapper to set a book on Redis
func (c *RedisClient) SetBook(ctx context.Context, id string, book model.Book) error {
	data, err := json.Marshal(book)
	if err != nil {
		log.Print(fmt.Errorf("could not marshal book to redis: %w", err))
		return err
	}

	set := c.client.Set(ctx, id, string(data), time.Hour)
	log.Print(fmt.Errorf("redis set status: %v", set))
	return nil
}

// DeleteEntry wrapper to delete a book entry from Redis
func (c *RedisClient) DeleteEntry(ctx context.Context, id string) {
	del := c.client.Del(ctx, id)
	log.Print(fmt.Printf("redis delete status for [%s]: %q", id, del))
}
