package main

import (
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
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
