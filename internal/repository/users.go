package repository

import (
	"context"
	"github.com/shamil/Test_task/internal/domain"
	"github.com/shamil/Test_task/internal/infrastructure/database"
)

func (r *Repository) UserSave(ctx context.Context, users ...domain.User) error {
	const query = `INSERT INTO users (username, password) VALUES ($1, $2)
                   ON CONFLICT(username) DO UPDATE SET password = excluded.password`

	err := database.WithTransaction(ctx, r.db, func(transaction database.Transaction) error {
		for _, user := range users {
			_, err := transaction.Exec(query, user.Username, user.Password)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
