package repository

import (
	todo "github.com/Hudayberdyyev/Rest_ToDo"
	"github.com/jackc/pgx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
