package datastore

import (
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/spf13/viper"
)

// CookieStoreClient is the client responsible for managing cookies
type CookieStoreClient struct {
	Config *viper.Viper
}

func (c CookieStoreClient) NewCookieStore() (redisStore.Store, error) {
	return redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte(c.Config.GetString("redis.sessionSecret")))
}
