package database

import (
	"database/sql"
	todotypes "demo/todolist/types"
	"errors"
)

type TodoDatabase struct {
	db *sql.DB
}

type dbTodo struct {
	id     int
	task   string
	note   string
	status string
}

func CreateTodoDatabase(db *sql.DB) (TodoDatabase, error) {
	if db == nil {
		return TodoDatabase{}, errors.New("Missing parameter 'db' when calling CreateTodoDatabase")
	}

	return TodoDatabase{
		db: db,
	}, nil
}

func (dbRow dbTodo) toTodo() todotypes.Todo {
	return todotypes.Todo{
		Id:     dbRow.id,
		Task:   dbRow.task,
		Note:   dbRow.note,
		Status: todotypes.TodoStatusFromString(dbRow.status),
	}
}

func (db TodoDatabase) GetTodoById(id int) (todotypes.Todo, error) {
	row := db.db.QueryRow("SELECT todo_status.status AS status, task, note FROM todo LEFT JOIN todo_status on todo.status_id=todo_status.id WHERE todo.id = ?;", id)

	dbRow := dbTodo{id: id}

	err := row.Scan(&dbRow.status, &dbRow.task, &dbRow.note)

	if err != nil {
		return todotypes.Todo{}, err
	}

	return dbRow.toTodo(), err
}

func (db TodoDatabase) GetTodos() ([]todotypes.Todo, error) {
	todos := []todotypes.Todo{}

	rows, err := db.db.Query("SELECT todo.id AS id, todo_status.status AS status, task, note FROM todo LEFT JOIN todo_status on todo.status_id=todo_status.id;")
	if err != nil {
		return todos, err
	}

	var foundError error = nil

	for rows.Next() {
		dbRow := dbTodo{}

		err := rows.Scan(&dbRow.id, &dbRow.status, &dbRow.task, &dbRow.note)
		if err != nil {
			foundError = err
		}

		todos = append(todos, dbRow.toTodo())
	}

	return todos, foundError
}

func (db TodoDatabase) SetTodoStatus(id int, status todotypes.TodoStatus) (todotypes.Todo, error) {
	row := db.db.QueryRow("UPDATE todo set status_id = (SELECT id FROM todo_status WHERE status = ?) WHERE todo.id = ? RETURNING todo.id, todo.task, todo.note;", status.ToUpperStr(), id)

	dbRow := dbTodo{status: status.ToUpperStr()}
	err := row.Scan(&dbRow.id, &dbRow.task, &dbRow.note)
	if err != nil {
		return todotypes.Todo{}, err
	}

	return dbRow.toTodo(), nil
}

func (db TodoDatabase) DeleteTodo(id int) error {
    tx, err := db.db.Begin()
    if err != nil {
        if tx != nil {
            tx.Rollback()
        }
        return err
    }

    res, err := db.db.Exec("DELETE FROM todo WHERE id = ?;", id)
    if err != nil {
        tx.Rollback()
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if  err != nil {
        if tx.Rollback() != nil {
            return err
        }
    }

    if rowsAffected > 1 {
        _ = tx.Rollback()
        return errors.New("Too many rows affected")
    }

    return tx.Commit();
}

func (db TodoDatabase) AddTodo(task string, note string) (todotypes.Todo, error) {
	row := db.db.QueryRow("INSERT INTO todo (id, task, note, status_id) VALUES (null, ?, ?, (SELECT id FROM todo_status WHERE status = ? LIMIT 1)) RETURNING todo.id, task, note;", task, note, todotypes.Pending.ToUpperStr())

	dbRow := dbTodo{status: todotypes.Pending.ToUpperStr()}
	err := row.Scan(&dbRow.id, &dbRow.task, &dbRow.note)
	if err != nil {
		return todotypes.Todo{}, err
	}

	return dbRow.toTodo(), err
}
