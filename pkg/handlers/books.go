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

// swagger:operation POST /books books addBook
// Adds a new book
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: book
//   in: body
//   name: Book
//   description: The new book to create
//   schema:
//         "$ref": "#/definitions/Book"
// responses:
//     '201':
//         description: Successful operation
//         schema:
//           type: array
//           items:
//                "$ref": "#/definitions/Book"
//     '400':
//         description: Invalid input
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

// swagger:operation GET /books books listBooks
// Returns list of books
// ---
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//         schema:
//           type: array
//           items:
//                "$ref": "#/definitions/Book"
func ListBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, books)
}

// swagger:operation GET /books/search books searchBooks
// Filters list of books by tag
// ---
// parameters:
// - name: tag
//   in: query
//   description: tag to filter on
//   required: false
//   type: string
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//         schema:
//           type: array
//           items:
//                "$ref": "#/definitions/Book"
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

// swagger:operation PUT /books/{id} books updateBook
// Update an existing book
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the book
//   required: true
//   type: string
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: book
//   in: body
//   name: Book
//   description: The new book to create
//   schema:
//         "$ref": "#/definitions/Book"
// responses:
//     '200':
//         description: Successful operation
//         schema:
//           type: array
//           items:
//                "$ref": "#/definitions/Book"
//     '400':
//         description: Invalid input
//     '404':
//         description: book not found
func UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
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

// swagger:operation DELETE /books/{id} books deleteBook
// Delete an existing book
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the book
//   required: true
//   type: string
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//     '200':
//         message: "Book has been deleted"
//     '404':
//         error: "book not found"
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
