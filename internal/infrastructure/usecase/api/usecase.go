package api

import (
	"context"
	"github.com/shamil/Test_task/internal/domain"
)

type NotesRepository interface {
	GetByUser(ctx context.Context, userID int) ([]domain.Note, error)
	Create(ctx context.Context, note domain.Note) (int64, error)
	Authenticate(ctx context.Context, username, password string) (int, error)
}

type UseCase struct {
	NotesRepository NotesRepository
}

func NewApiUseCase(notesrepository NotesRepository) *UseCase {
	return &UseCase{NotesRepository: notesrepository}
}

func (u *UseCase) GetByUser(ctx context.Context, userID int) ([]domain.Note, error) {
	return u.NotesRepository.GetByUser(ctx, userID)
}

func (u *UseCase) Create(ctx context.Context, note domain.Note) (int64, error) {
	return u.NotesRepository.Create(ctx, note)
}

func (u *UseCase) Authenticate(ctx context.Context, username, password string) (int, error) {
	return u.NotesRepository.Authenticate(ctx, username, password)
}
