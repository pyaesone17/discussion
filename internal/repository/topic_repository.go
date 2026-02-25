package repository

import (
	"database/sql"
	"discussion-forum/internal/models"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(topic *models.CreateTopicRequest) (*models.Topic, error) {
	result, err := r.db.Exec(
		"INSERT INTO topics (title, user_id) VALUES (?, ?)",
		topic.Title, topic.UserID,
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

func (r *TopicRepository) GetByID(id int64) (*models.Topic, error) {
	topic := &models.Topic{}
	err := r.db.QueryRow(`
		SELECT t.id, t.title, t.user_id, u.username, t.created_at, t.updated_at,
			COALESCE((SELECT SUM(value) FROM votes WHERE votable_type='topic' AND votable_id=t.id), 0) AS score
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`, id).Scan(&topic.ID, &topic.Title, &topic.UserID, &topic.Username, &topic.CreatedAt, &topic.UpdatedAt, &topic.Score)

	if err != nil {
		return nil, err
	}

	return topic, nil
}

func (r *TopicRepository) GetAll() ([]*models.Topic, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.user_id, u.username, t.created_at, t.updated_at,
			COALESCE((SELECT SUM(value) FROM votes WHERE votable_type='topic' AND votable_id=t.id), 0) AS score
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := []*models.Topic{}
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Title, &topic.UserID, &topic.Username, &topic.CreatedAt, &topic.UpdatedAt, &topic.Score); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func (r *TopicRepository) Update(id int64, req *models.UpdateTopicRequest) (*models.Topic, error) {
	if req.Title != "" {
		_, err := r.db.Exec("UPDATE topics SET title = ? WHERE id = ?", req.Title, id)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

func (r *TopicRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM topics WHERE id = ?", id)
	return err
}

func (r *TopicRepository) Search(query string) ([]*models.Topic, error) {
	if query == "" {
		return []*models.Topic{}, nil
	}

	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.user_id, u.username, t.created_at, t.updated_at,
			COALESCE((SELECT SUM(value) FROM votes WHERE votable_type='topic' AND votable_id=t.id), 0) AS score
		FROM topics t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE MATCH(t.title) AGAINST(? IN NATURAL LANGUAGE MODE)
		ORDER BY MATCH(t.title) AGAINST(? IN NATURAL LANGUAGE MODE) DESC
	`, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := []*models.Topic{}
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Title, &topic.UserID, &topic.Username, &topic.CreatedAt, &topic.UpdatedAt, &topic.Score); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
