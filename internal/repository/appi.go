package repository

import (
	"context"
	"errors"
	"github.com/shamil/Test_task/internal/domain"
)

func (r *Repository) Create(ctx context.Context, note domain.Note) (int64, error) {
	query := `INSERT INTO notes (user_id, title, body) 
	VALUES ($1, $2, $3)`

	result, err := r.db.ExecContext(ctx, query,
		note.UserID, note.Title, note.Body)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (r *Repository) GetByUser(ctx context.Context, userID int) ([]domain.Note, error) {
	query := `SELECT id, title, body 
	FROM notes 
	WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Body); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *Repository) Authenticate(ctx context.Context, username, password string) (int, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return 0, err
	}

	if user.Password != password {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}
