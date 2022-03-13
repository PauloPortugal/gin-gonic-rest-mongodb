package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var books []model.Book

func init() {
	books = make([]model.Book, 0)

	file, _ := ioutil.ReadFile("books.json")
	_ = json.Unmarshal(file, &books)
}

func NewBook(ctx *gin.Context) {
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	book.ID = xid.New().String()
	book.SubmissionDate = time.Now()
	books = append(books, book)
	ctx.JSON(http.StatusCreated, book)
}

func ListBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, books)
}

func SearchBooks(ctx *gin.Context) {
	tag := ctx.Query("tag")
	listOfBooks := make([]model.Book, 0)

	for _, b := range books {
		for _, t := range b.Tags {
			if strings.Contains(t, tag) {
				listOfBooks = append(listOfBooks, b)
			}
		}
	}

	ctx.JSON(http.StatusOK, listOfBooks)
}

func UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	found := false
	for i, v := range books {
		if v.ID == id {
			found = true
			book.ID = id
			book.SubmissionDate = books[i].SubmissionDate
			books[i] = book
		}
	}

	if found {
		ctx.JSON(http.StatusOK, book)
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "book not found",
	})
}

func DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	found := false
	for i, v := range books {
		if v.ID == id {
			found = true
			books = append(books[:i], books[i+1:]...)
		}
	}
	if found {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Book has been deleted",
		})
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "book not found",
	})
}
