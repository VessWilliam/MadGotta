package repository

import (
	"gotta/internal/domains"
	"gotta/internal/models"

	"github.com/jmoiron/sqlx"
)

var _ domains.ITodoRepo = (*todoRepository)(nil)

type todoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) domains.ITodoRepo {
	return &todoRepository{db: db}
}

// Add implements [domains.ITodoRepo].
func (t *todoRepository) Add(text string) error {
	_, err := t.db.Exec("INSERT INTO todos (text, completed) VALUES (?, ?)", text, false)
	return err
}

// All implements [domains.ITodoRepo].
func (t *todoRepository) All() ([]models.Todo, error) {
	var todos []models.Todo
	err := t.db.Select(&todos, "SELECT id, text, completed FROM todos ORDER BY id")
	return todos, err
}

// Delete implements [domains.ITodoRepo].
func (t *todoRepository) Delete(id int) error {
	_, err := t.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// Toggle implements [domains.ITodoRepo].
func (t *todoRepository) Toggle(id int) error {
	_, err := t.db.Exec("UPDATE todos SET completed = NOT completed WHERE id = ?", id)
	return err
}
