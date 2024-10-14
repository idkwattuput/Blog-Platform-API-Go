package types

import (
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostStore interface {
	GetPosts() ([]*Post, error)
	GetPostById(id int) (*Post, error)
	CreatePost(PostPayload) (*Post, error)
	UpdatePost(id int, post PostPayload) (*Post, error)
	DeletePost(id int) error
}

type PostPayload struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}
