package domains

import "gotta/internal/models"

type ITodoRepo interface {
	All() ([]models.Todo, error)
	Add(text string) error
	Toggle(id int) error
	Delete(id int) error
}
