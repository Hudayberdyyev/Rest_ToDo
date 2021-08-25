package repository

import (
	"fmt"
	todo "github.com/Hudayberdyyev/Rest_ToDo"
	"github.com/jackc/pgx"
)

type TodoListPostgres struct {
	db *pgx.Conn
}

func NewTodoListPostgres(db *pgx.Conn) *TodoListPostgres {
	return &TodoListPostgres{ db: db }
}

func (s *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	createListQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoListsTable)
	if err = tx.QueryRow(createListQuery, list.Title, list.Description).Scan(&listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("insert into %s (user_id, list_id) values ($1, $2)", usersListsTable)
	if _, err = tx.Exec(createUsersListsQuery, userId, listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (s *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id=$1",
		todoListsTable, usersListsTable)
	rows, err := s.db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var list todo.TodoList
		err = rows.Scan(&list.Id, &list.Title, &list.Description)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, err
}

func (s *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf(`select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id 
								where ul.user_id=$1 and ul.list_id=$2`, todoListsTable, usersListsTable)
	err := s.db.QueryRow(query, userId, listId).Scan(&list.Id, &list.Title, &list.Description)
	return list, err
}
