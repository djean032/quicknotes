package note

import (
	"time"
)

type Note struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Created time.Time `json:"created"`
	Tags    []string  `json:"tags"`
	Links   []string  `json:"links"`
}
