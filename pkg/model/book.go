package model

import "time"

type Book struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Author         string        `json:"author"`
	Publisher      string        `json:"publisher"`
	PublishedAt    PublishedDate `json:"published_at"`
	Tags           []string      `json:"tags"`
	Review         float32       `json:"review"`
	SubmissionDate time.Time     `json:"submission_date"`
}

type PublishedDate struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}
