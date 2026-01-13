package repository

import (
	"database/sql"
	"discussion-forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.CreateUserRequest) (*models.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		user.Username, user.Email,
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

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	rows, err := r.db.Query(
		"SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	if req.Username != "" {
		_, err := r.db.Exec("UPDATE users SET username = ? WHERE id = ?", req.Username, id)
		if err != nil {
			return nil, err
		}
	}

	if req.Email != "" {
		_, err := r.db.Exec("UPDATE users SET email = ? WHERE id = ?", req.Email, id)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
