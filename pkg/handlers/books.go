package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BooksHandler provides a struct to wrap the api around
type BooksHandler struct {
	cfg *viper.Viper
	col *mongo.Collection
	ctx context.Context
}

func New(ctx context.Context, cfg *viper.Viper, col *mongo.Collection) *BooksHandler {
	return &BooksHandler{
		ctx: ctx,
		cfg: cfg,
		col: col,
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
//     '500':
//         description: internal server error
func (handler *BooksHandler) NewBook(ctx *gin.Context) {
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	book.ID = primitive.NewObjectID()
	book.SubmissionDate = time.Now()
	_, err := handler.col.InsertOne(handler.ctx, book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
//     '500':
//         description: internal server error
func (handler *BooksHandler) ListBooks(ctx *gin.Context) {
	books := make([]model.Book, 0)
	cur, err := handler.col.Find(handler.ctx, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err = cur.All(handler.ctx, &books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, &books)
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
func (handler *BooksHandler) SearchBooks(ctx *gin.Context) {
	tag := ctx.Query("tag")
	books := make([]model.Book, 0)
	cur, err := handler.col.Find(handler.ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err = cur.All(handler.ctx, &books); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, books)
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
//     '500':
//         description: internal server error
func (handler *BooksHandler) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.col.UpdateOne(handler.ctx, bson.M{"_id": objID}, bson.D{{ //nolint:govet
		"$set", bson.D{
			{"name", book.Name},           //nolint:govet
			{"author", book.Author},       //nolint:govet
			{"publisher", book.Publisher}, //nolint:govet
			{"tags", book.Tags},           //nolint:govet
			{"review", book.Review},       //nolint:govet
		},
	}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if res.MatchedCount == 1 {
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
func (handler *BooksHandler) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.col.DeleteOne(handler.ctx, bson.M{"_id": objID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if res.DeletedCount == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Book has been deleted",
		})
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "book not found",
	})
}
