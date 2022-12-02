package handlers

import (
	"context"
	"net/http"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// BooksHandler provides a struct to wrap the api around
type BooksHandler struct {
	ctx          context.Context
	cfg          *viper.Viper
	mongoDBStore datastore.Books
	redisStore   datastore.Redis
}

func NewBooksHandler(ctx context.Context, cfg *viper.Viper, mongoDBStore datastore.Books, redisStore datastore.Redis) *BooksHandler {
	return &BooksHandler{
		ctx:          ctx,
		cfg:          cfg,
		mongoDBStore: mongoDBStore,
		redisStore:   redisStore,
	}
}

// swagger:operation POST /books books addBook
// Adds a new book
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: book
//     in: body
//     name: Book
//     description: The new book to create
//     schema:
//     "$ref": "#/definitions/Book"
//   - name: Authorization
//     in: header
//     name: Authorization
//     description: Auth header, where the JWT token should be provided
//     type: string
//     example: someJWTToken
//     required: true
//
// responses:
//
//	'201':
//	    description: Successful operation
//	    schema:
//	      type: array
//	      items:
//	           "$ref": "#/definitions/Book"
//	'400':
//	    description: invalid input
//	'500':
//	    description: internal server error
func (handler *BooksHandler) NewBook(ctx *gin.Context) {
	var book *model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	if err := handler.mongoDBStore.AddBook(handler.ctx, book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	handler.redisStore.DeleteEntry(handler.ctx, "books")

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
//
//	'200':
//	    description: Successful operation
//	    schema:
//	      type: array
//	      items:
//	           "$ref": "#/definitions/Book"
//	'500':
//	    description: internal server error
func (handler *BooksHandler) ListBooks(ctx *gin.Context) {
	// first query Redis/cache
	books, err := handler.redisStore.GetBooks(handler.ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// if not available on Redis DB, then query MongoDB
	if books == nil {
		books, err = handler.mongoDBStore.ListBooks(handler.ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		err := handler.redisStore.SetBooks(handler.ctx, books)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, books)
}

// swagger:operation GET /books/search books searchBooks
// Filters list of books by tag
// ---
// parameters:
//   - name: tag
//     in: query
//     description: tag to filter on
//     required: false
//     type: string
//
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	    schema:
//	      type: array
//	      items:
//	           "$ref": "#/definitions/Book"
func (handler *BooksHandler) SearchBooks(ctx *gin.Context) {
	tag := ctx.Query("tag")

	books, err := handler.mongoDBStore.SearchBooks(handler.ctx, tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

// swagger:operation PUT /books/{id} books updateBook
// Update an existing book
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the book
//     required: true
//     type: string
//
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: book
//     in: body
//     name: Book
//     description: The new book to create
//     schema:
//     "$ref": "#/definitions/Book"
//   - name: Authorization
//     in: header
//     name: Authorization
//     description: Auth header, where the JWT token should be provided
//     type: string
//     example: someJWTToken
//     required: true
//
// responses:
//
//	'200':
//	    description: Successful operation
//	    schema:
//	      type: array
//	      items:
//	           "$ref": "#/definitions/Book"
//	'400':
//	    description: Invalid input
//	'404':
//	    description: book not found
//	'500':
//	    description: internal server error
func (handler *BooksHandler) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	modifiedCount, err := handler.mongoDBStore.UpdateBook(handler.ctx, id, book)
	if err != nil {
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "book not found",
		})
		return
	}

	handler.redisStore.DeleteEntry(handler.ctx, id)

	ctx.JSON(http.StatusOK, book)
}

// swagger:operation DELETE /books/{id} books deleteBook
// Delete an existing book
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the book
//     required: true
//     type: string
//   - name: Authorization
//     in: header
//     name: Authorization
//     description: Auth header, where the JWT token should be provided
//     type: string
//     example: someJWTToken
//     required: true
//
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//
//	'200':
//	    message: "Book has been deleted"
//	'404':
//	    error: "book not found"
func (handler *BooksHandler) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedCount, err := handler.mongoDBStore.DeleteBook(handler.ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "book not found",
		})
		return
	}

	handler.redisStore.DeleteEntry(handler.ctx, id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book has been deleted",
	})
}

// swagger:operation GET /books/{id} books getBook
// Returns a book
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the book
//     required: true
//     type: string
//
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//
//	'200':
//	    schema:
//	      items:
//	           "$ref": "#/definitions/Book"
//	'500':
//	    error: error description
func (handler *BooksHandler) GetBook(ctx *gin.Context) {
	id := ctx.Param("id")

	book, err := handler.redisStore.GetBook(handler.ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if book.Name == "" {
		book, err = handler.mongoDBStore.GetBook(handler.ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if book.Name == "" {
			ctx.Status(http.StatusNotFound)
			return
		}

		if err := handler.redisStore.SetBook(handler.ctx, id, book); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, book)
}
