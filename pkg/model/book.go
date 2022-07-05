package model

import "time"

// Book represents a book
//
// swagger:model
type Book struct {
	// the id for this book
	// required: true
	// min: 1
	ID string `json:"id"`

	// the name of the book
	// required: true
	// min length: 3
	Name string `json:"name"`

	// the book author's name
	// required: true
	// min length: 3
	Author string `json:"author"`

	// The publisher's name
	// required: true
	// min length: 3
	Publisher string `json:"publisher"`

	// The date the book was published
	// required: true
	PublishedAt PublishedDate `json:"published_at"`

	// the associated tags with this book
	Tags []string `json:"tags"`

	// Score review
	Review float32 `json:"review"`

	//The day the book was added
	SubmissionDate time.Time `json:"submission_date"`
}

// PublishedDate Represents the month and year of when the book was published
//
// swagger:model
type PublishedDate struct {
	// the month of the year
	// required: true
	// length: 2
	Month string `json:"month"`

	// the month of the year
	// required: true
	// min length: 1
	Year string `json:"year"`
}
