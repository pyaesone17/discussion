package repository

import (
	"database/sql"
	"discussion-forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.CreatePostRequest) (*models.Post, error) {
	result, err := r.db.Exec(
		"INSERT INTO posts (topic_id, user_id, content) VALUES (?, ?, ?)",
		post.TopicID, post.UserID, post.Content,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *PostRepository) GetByID(id int64) (*models.Post, error) {
	post := &models.Post{}
	err := r.db.QueryRow(`
		SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.created_at, p.updated_at,
			COALESCE((SELECT SUM(value) FROM votes WHERE votable_type='post' AND votable_id=p.id), 0) AS score
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&post.ID, &post.TopicID, &post.UserID, &post.Username, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Score)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) GetByTopicID(topicID int64) ([]*models.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.created_at, p.updated_at,
			COALESCE((SELECT SUM(value) FROM votes WHERE votable_type='post' AND votable_id=p.id), 0) AS score
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.topic_id = ?
		ORDER BY p.created_at ASC
	`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(&post.ID, &post.TopicID, &post.UserID, &post.Username, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Score); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) Update(id int64, req *models.UpdatePostRequest) (*models.Post, error) {
	if req.Content != "" {
		_, err := r.db.Exec("UPDATE posts SET content = ? WHERE id = ?", req.Content, id)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

func (r *PostRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
