package repository

import "github.com/jackc/pgx"

type Authorization interface {
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
	return &Repository{}
}
