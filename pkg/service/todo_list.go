package service

import (
	todo "github.com/Hudayberdyyev/Rest_ToDo"
	"github.com/Hudayberdyyev/Rest_ToDo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{ repo: repo }
}

func (s *TodoListService) CreateList(userId int, list todo.TodoList) (int, error) {
	return s.repo.CreateList(userId, list)
}

func (s *TodoListService) GetAllLists(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAllLists(userId)
}

func (s *TodoListService) GetListById(userId,listId int) (todo.TodoList, error) {
	return s.repo.GetListById(userId, listId)
}

func (s *TodoListService) Delete(userId, listId int) error {
	if _, err := s.repo.GetListById(userId, listId); err != nil {
		return err
	}
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := s.repo.GetListById(userId, listId); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}