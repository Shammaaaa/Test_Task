package updater

import (
	"context"
	"github.com/shamil/Test_task/config"
	"github.com/shamil/Test_task/internal/domain"
	"github.com/shamil/Test_task/pkg/log"
)

type NotesRepository interface {
	UserSave(ctx context.Context, users ...domain.User) error
}

type UseCase struct {
	notesRepository NotesRepository
}

func NewUpdaterUseCase(notesRepository NotesRepository) *UseCase {
	return &UseCase{
		notesRepository: notesRepository,
	}
}

func (u *UseCase) Work(ctx context.Context) {
	if err := u.notesRepository.UserSave(ctx, config.Users...); err != nil {
		log.Warningf("failed to save users: %s", err)
	}
}
