package hs

import (
	"context"
	"net/http"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	indexTemplate = "index.tmpl"
	bookTemplate  = "book.tmpl"
	notFoundPage  = "./web/static/404.html"
	errorPage     = "./web/static/500.html"
)

// Booksh provides a struct to wrap the api around
type WebHandler struct {
	c            context.Context
	cfg          *viper.Viper
	mongoDBStore datastore.Books
	redisStore   datastore.Redis
}

func NewWebHandler(c context.Context, cfg *viper.Viper, booksCli datastore.Books, redisCli datastore.Redis) *WebHandler {
	return &WebHandler{
		c:            c,
		cfg:          cfg,
		mongoDBStore: booksCli,
		redisStore:   redisCli,
	}
}

func (h *WebHandler) IndexPage(c *gin.Context) {
	books, err := h.mongoDBStore.ListBooks(h.c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, indexTemplate, gin.H{
		"books": books,
	})
}

func (h *WebHandler) BookPage(c *gin.Context) {
	id := c.Param("id")

	book, err := h.redisStore.GetBook(h.c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if book.Name == "" {
		book, err = h.mongoDBStore.GetBook(h.c, id)
		if err != nil {
			c.File(errorPage)
			return
		}

		if book.Name == "" {
			c.Status(http.StatusNotFound)
			c.File(notFoundPage)
			return
		}

		if err := h.redisStore.SetBook(h.c, id, book); err != nil {
			c.File(errorPage)
			return
		}
	}

	c.HTML(http.StatusOK, bookTemplate, gin.H{
		"book": book,
	})
}
