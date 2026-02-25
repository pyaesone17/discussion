package repository

import (
	"database/sql"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) Upsert(userID, votableID int64, votableType string, value int) error {
	_, err := r.db.Exec(`
		INSERT INTO votes (user_id, votable_type, votable_id, value)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value), updated_at = CURRENT_TIMESTAMP
	`, userID, votableType, votableID, value)
	return err
}

func (r *VoteRepository) Delete(userID, votableID int64, votableType string) error {
	_, err := r.db.Exec(`
		DELETE FROM votes WHERE user_id = ? AND votable_type = ? AND votable_id = ?
	`, userID, votableType, votableID)
	return err
}

func (r *VoteRepository) GetScore(votableID int64, votableType string) (int64, error) {
	var score sql.NullInt64
	err := r.db.QueryRow(`
		SELECT SUM(value) FROM votes WHERE votable_type = ? AND votable_id = ?
	`, votableType, votableID).Scan(&score)
	if err != nil {
		return 0, err
	}
	return score.Int64, nil
}
