package repository

import (
	"fmt"
	todo "github.com/Hudayberdyyev/Rest_ToDo"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{ db: db }
}

func (r *TodoListPostgres) CreateList(userId int, list todo.TodoList) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("insert into %s (user_id, list_id) values ($1, $2) returning id", usersListsTable)
	row = tx.QueryRow(createUsersListQuery, userId, id)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAllLists(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id=$1",
		todoListsTable, usersListsTable)
	if err := r.db.Select(&lists, query, userId); err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoListPostgres) GetListById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf(`select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id 
										where ul.user_id=$1 and ul.list_id = $2`,todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("delete from %s tl using %s ul where tl.id = ul.list_id and ul.user_id=$1 and ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)
	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	argv := make([]interface{}, 0)
	setValues :=  make([]string, 0)
	argc := 1

	if input.Title != nil {
		argv = append(argv, *input.Title)
		setValues = append(setValues, fmt.Sprintf("title=$%d", argc))
		argc++
	}

	if input.Description != nil {
		argv = append(argv, *input.Description)
		setValues = append(setValues, fmt.Sprintf("description=$%d", argc))
		argc++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("update %s tl set %s from %s ul where tl.id = ul.list_id and ul.list_id=$%d and ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argc, argc+1)
	argv = append(argv, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", argv)

	_, err := r.db.Exec(query, argv...)
	return err
}