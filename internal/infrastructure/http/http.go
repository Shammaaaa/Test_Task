package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/shamil/Test_task/internal/domain"
	"github.com/shamil/Test_task/internal/infrastructure/auth"
	"github.com/shamil/Test_task/pkg/speller"
	"net/http"
)

type NotesUsecase interface {
	GetByUser(ctx context.Context, userID int) ([]domain.Note, error)
	Create(ctx context.Context, note domain.Note) (int64, error)
	Authenticate(ctx context.Context, username, password string) (int, error)
}

type HandlerImpl struct {
	notesUsecase NotesUsecase
}

func New(useCase NotesUsecase) *HandlerImpl {
	return &HandlerImpl{notesUsecase: useCase}
}

func (h *HandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	var note domain.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenStr := r.Header.Get("Authorization")
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	note.UserID = int64(claims.UserID)

	titleSuggestions, err := speller.CheckSpelling(note.Title)
	if err != nil {
		http.Error(w, "Failed to check spelling for title", http.StatusInternalServerError)
		return
	}

	bodySuggestions, err := speller.CheckSpelling(note.Body)
	if err != nil {
		http.Error(w, "Failed to check spelling for body", http.StatusInternalServerError)
		return
	}

	// Если есть ошибки, возвращаем их пользователю
	if len(titleSuggestions) > 0 || len(bodySuggestions) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"title_suggestions": titleSuggestions,
			"body_suggestions":  bodySuggestions,
		})
		return
	}

	// Если ошибок нет, сохраняем заметку
	id, err := h.notesUsecase.Create(r.Context(), note)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *HandlerImpl) GetByUser(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	notes, err := h.notesUsecase.GetByUser(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *HandlerImpl) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.notesUsecase.Authenticate(r.Context(), user.Username, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(userID)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *HandlerImpl) MountRoutes(r chi.Router) {
	r.Post("/create", h.Create)
	r.Get("/notes", h.GetByUser)
	r.Post("/login", h.Authenticate)
}
