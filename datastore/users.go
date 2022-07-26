package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Users interface {
	Get(ctx context.Context, username string) (model.User, error)
}

// BooksClient is the client responsible for querying mongodb
type UsersClient struct {
	client *mongo.Client
	cfg    *viper.Viper
	col    *mongo.Collection
}

//Get returns a user by username
func (c *UsersClient) Get(ctx context.Context, username string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

// NewUsersClient create a new UsersClient
func NewUsersClient(client *mongo.Client, cfg *viper.Viper) *UsersClient {
	return &UsersClient{
		client: client,
		cfg:    cfg,
		col:    getCollection(cfg, client, "mongodb.dbcollections.users"),
	}
}

func (c *UsersClient) InitUsers(ctx context.Context) {
	setupIndexes(ctx, c.col, "username")

	if err := loadDefaultUsers(ctx, c.col); err != nil {
		log.Fatal(fmt.Errorf("could not insert static data: %w\n", err))
	}
}

func loadDefaultUsers(ctx context.Context, collection *mongo.Collection) error {
	users := make([]model.User, 0)

	file, err := ioutil.ReadFile("default_data/users.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &users); err != nil {
		return err
	}

	var b []interface{}
	for _, user := range users {
		if err = HashPassword(&user); err != nil {
			return err
		}
		b = append(b, user)
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

	log.Printf("Inserted users: %d\n", len(result.InsertedIDs))

	return nil
}

func HashPassword(user *model.User) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hashedPwd)
	return nil
}
