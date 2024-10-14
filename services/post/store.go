package post

import (
	"database/sql"
	"time"

	"github.com/idkwattuput/blogging-platform-api-go/types"
	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetPosts() ([]*types.Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	posts := make([]*types.Post, 0)
	for rows.Next() {
		p, err := scanRowsIntoPost(rows)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (s *Store) GetPostById(id int) (*types.Post, error) {
	row := s.db.QueryRow("SELECT * FROM posts WHERE id=$1", id)

	p := new(types.Post)
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		pq.Array(&p.Tags),
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) CreatePost(post types.PostPayload) (*types.Post, error) {
	row := s.db.QueryRow("INSERT INTO posts (title, content, category, tags) VALUES ($1, $2, $3, $4) RETURNING *", post.Title, post.Content, post.Category, pq.Array(post.Tags))
	p := new(types.Post)
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		pq.Array(&p.Tags),
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) UpdatePost(id int, post types.PostPayload) (*types.Post, error) {
	row := s.db.QueryRow(
		"UPDATE posts SET title = $1, content = $2, category = $3, tags = $4, updatedAt = $5 WHERE id = $6 RETURNING *",
		post.Title, post.Content, post.Category, pq.Array(post.Tags), time.Now(), id)
	p := new(types.Post)
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		pq.Array(&p.Tags),
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) DeletePost(id int) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoPost(rows *sql.Rows) (*types.Post, error) {
	post := new(types.Post)

	err := rows.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Category,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}
