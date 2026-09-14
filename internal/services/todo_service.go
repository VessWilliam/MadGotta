package services

import (
	"errors"
	"gotta/internal/domains"
	"gotta/internal/models"
	"strings"
)

type TodoService struct {
	repo domains.ITodoRepo
}

func NewTodoService(repo domains.ITodoRepo) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) TodoList() ([]models.Todo, error) {
	return s.repo.All()
}

func (s *TodoService) Add(text string) ([]models.Todo, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return nil, errors.New("todo text cannot be empty")
	}

	if len(text) > 200 {
		return nil, errors.New("todo text too long")
	}

	if err := s.repo.Add(text); err != nil {
		return nil, err
	}

	return s.repo.All()
}

func (s *TodoService) Toggle(id int) ([]models.Todo, error) {
	if err := s.repo.Toggle(id); err != nil {
		return nil, err
	}

	return s.repo.All()
}

func (s *TodoService) Delete(id int) ([]models.Todo, error) {
	if err := s.repo.Delete(id); err != nil {
		return nil, err
	}

	return s.repo.All()
}
