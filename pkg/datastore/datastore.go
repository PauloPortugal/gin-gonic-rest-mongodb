package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore interface {
	AddBook(ctx context.Context, book *model.Book) error
	ListBooks(ctx context.Context) ([]model.Book, error)
	SearchBooks(ctx context.Context, tag string) ([]model.Book, error)
	UpdateBook(ctx context.Context, id string, book model.Book) (int, error)
	DeleteBook(ctx context.Context, id string) (int, error)
}

// Client is the client responsible for querying mongodb
type Client struct {
	client *mongo.Client
	cfg    *viper.Viper
	col    *mongo.Collection
}

func New(client *mongo.Client, cfg *viper.Viper) *Client {
	return &Client{
		client: client,
		cfg:    cfg,
		col:    getBooksCollection(cfg, client),
	}
}

func (c *Client) Init(ctx context.Context) {
	booksCollection := getBooksCollection(c.cfg, c.client)

	setupIndexes(ctx, booksCollection, c.cfg)

	if err := loadStaticData(ctx, booksCollection); err != nil {
		log.Fatal(fmt.Errorf("could not insert static data: %w\n", err))
	}
}

// AddBook wrapper to add a book to the MongoDB collection
func (c *Client) AddBook(ctx context.Context, book *model.Book) error {
	book.ID = primitive.NewObjectID()
	book.SubmissionDate = time.Now()
	_, err := c.col.InsertOne(ctx, book)
	if err != nil {
		log.Print(fmt.Errorf("could not add new book: %w", err))
		return err
	}
	return nil
}

// ListBooks wrapper to return all books from the MongoDB collection
func (c *Client) ListBooks(ctx context.Context) ([]model.Book, error) {
	books := make([]model.Book, 0)
	cur, err := c.col.Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all books: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &books); err != nil {
		log.Print(fmt.Errorf("could marshall the books results: %w", err))
		return nil, err
	}

	return books, nil
}

// SearchBooks wrapper to return all books based on a 'tag' from the MongoDB collection
func (c *Client) SearchBooks(ctx context.Context, tag string) ([]model.Book, error) {
	books := make([]model.Book, 0)
	cur, err := c.col.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		log.Print(fmt.Errorf("could not search book using tag [%s]: %w", tag, err))
		return nil, err
	}

	if err := cur.All(ctx, &books); err != nil {
		log.Print(fmt.Errorf("could marshall the books results: %w", err))
		return nil, err
	}

	return books, nil
}

func (c *Client) UpdateBook(ctx context.Context, id string, book model.Book) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{{ //nolint:govet
		"$set", bson.D{
			{"name", book.Name},           //nolint:govet
			{"author", book.Author},       //nolint:govet
			{"publisher", book.Publisher}, //nolint:govet
			{"tags", book.Tags},           //nolint:govet
			{"review", book.Review},       //nolint:govet
		},
	}})
	if err != nil {
		log.Print(fmt.Errorf("could not update book with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

// DeleteBook wrapper to delete a book from the MongoDB collection
func (c *Client) DeleteBook(ctx context.Context, id string) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("could marshall the books results: %w", err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}

func getBooksCollection(cfg *viper.Viper, client *mongo.Client) *mongo.Collection {
	db := fmt.Sprint(cfg.Get("mongodb.dbname"))
	col := fmt.Sprint(cfg.Get("mongodb.dbcollection"))

	return client.Database(db).Collection(col)
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
