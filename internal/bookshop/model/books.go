package model

import "time"

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Authors       []string  `json:"authors"`
	PublishedDate string    `json:"published_date"`
	PageCount     PageCount `json:"page_count"`
	Categories    []string  `json:"categories"`
	Version       int32     `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
