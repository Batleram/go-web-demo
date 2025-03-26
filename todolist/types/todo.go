package types

import (
	"strings"
)

type TodoStatus int64

const (
	Unspecified TodoStatus = 0
	Pending     TodoStatus = 1
	InProgress  TodoStatus = 2
	Completed   TodoStatus = 3
)

func TodoStatusFromInt(source int) TodoStatus {
	switch source {
	case 1:
		return Pending
	case 2:
		return InProgress
	case 3:
		return Completed
	default:
		return Unspecified
	}
}

func TodoStatusFromString(source string) TodoStatus {
	switch strings.ToUpper(source) {
	case "PENDING":
		return Pending
	case "INPROGRESS":
		return InProgress
	case "COMPLETED":
		return Completed
	default:
		return Unspecified
	}
}

func (todoStatus TodoStatus) ToStr() string {
	switch todoStatus {
	case Pending:
		return "Pending"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	default:
		return "Unspecified"
	}
}

func (todoStatus TodoStatus) ToUpperStr() string {
	return strings.ToUpper(todoStatus.ToStr())
}

type Todo struct {
	Id     int
	Task   string
	Note   string
	Status TodoStatus
}

func (todo Todo) IsComplete() bool {
	return todo.Status == Completed
}

func (todo Todo) IsInProgress() bool {
	return todo.Status == InProgress
}
