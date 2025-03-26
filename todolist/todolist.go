package todolist

import (
	ct "demo/core/templates"
	tdb "demo/todolist/database"
	t "demo/todolist/templates"
	"demo/todolist/types"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"

	"context"
	"strconv"
)

var todoDatabase *tdb.TodoDatabase = nil

func RegisterRoutes(db *tdb.TodoDatabase) (*fiber.App, error) {
	if db == nil {
		return nil, errors.New("Missing paramater db")
	}

	todoDatabase = db

	router := fiber.New()
	router.Get("/app", appHandler)
	router.Patch("/todo", todoPatchHandler)
	router.Delete("/todo", todoDeleteHandler)
	router.Post("/todo", todoPostHandler)

	return router, nil
}

func appHandler(c *fiber.Ctx) error {
	c.Set("Content-type", "text/html")

	todos, err := todoDatabase.GetTodos()
	if err != nil {
		fmt.Println("Internal error fetching todos: ", err.Error())
	}

	state := t.AppState{
		AddTodoFormState: t.AddTodoFormState{
			ErrorString: "",
		},
		Todos: todos}

	app := t.App(state)

	if c.Query("source", "") == "nav" {
		return app.Render(context.Background(), c.Response().BodyWriter())
	}

	return ct.Index(app).Render(context.Background(), c.Response().BodyWriter())

}

func todoPatchHandler(c *fiber.Ctx) error {
	todo_id, err := strconv.Atoi(c.Query("todo", ""))
	if err != nil || todo_id < 1 {
		c.WriteString("Please specify a todo")
		return c.SendStatus(400)
	}

	action := c.Query("action", "")
	if action == "" {
		c.WriteString("Please specify an action")
		return c.SendStatus(400)
	}

	todo, err := todoDatabase.GetTodoById(todo_id)
	if err != nil {
		c.WriteString("Could not find such todo")
		return c.SendStatus(400)
	}

	switch strings.ToLower(action) {
	case "complete":
		return completeTodo(todo, c)
	case "start":
		return startTodo(todo, c)
	case "pause":
		return pauseTodo(todo, c)
	default:
		c.WriteString("Invalid action")
		return c.SendStatus(400)
	}
}

func completeTodo(todo types.Todo, c *fiber.Ctx) error {
	todo, err := todoDatabase.SetTodoStatus(todo.Id, types.Completed)
	if err != nil {
        fmt.Println("Error completing todo: " + err.Error())
		return c.SendStatus(500)
	}

	c.Set("Content-type", "text/html")

	return t.Todo(todo).Render(context.Background(), c.Response().BodyWriter())
}

func pauseTodo(todo types.Todo, c *fiber.Ctx) error {
	validStates := []types.TodoStatus{types.InProgress, types.Pending}
	if !slices.Contains(validStates, todo.Status) {
		c.WriteString("Todo not in progress")
		return c.SendStatus(400)
	}

	todo, err := todoDatabase.SetTodoStatus(todo.Id, types.Pending)
	if err != nil {
        fmt.Println("Error pausing todo: " + err.Error())
		return c.SendStatus(500)
	}

	return t.Todo(todo).Render(context.Background(), c.Response().BodyWriter())
}

func startTodo(todo types.Todo, c *fiber.Ctx) error {
	validStates := []types.TodoStatus{types.InProgress, types.Pending}
	if !slices.Contains(validStates, todo.Status) {
		c.WriteString("Todo not pending")
		return c.SendStatus(400)
	}

	todo, err := todoDatabase.SetTodoStatus(todo.Id, types.InProgress)
	if err != nil {
        fmt.Println("Error starting todo: " + err.Error())
		return c.SendStatus(500)
	}

	return t.Todo(todo).Render(context.Background(), c.Response().BodyWriter())
}

func todoDeleteHandler(c *fiber.Ctx) error {
	todo_id, err := strconv.Atoi(c.Query("todo", ""))
	if err != nil || todo_id < 1 {
		c.WriteString("Please specify a todo")
		return c.SendStatus(400)
	}

	todo, err := todoDatabase.GetTodoById(todo_id)
	if err != nil {
		c.WriteString("Could not find such todo")
		return c.SendStatus(400)
	}

    err = todoDatabase.DeleteTodo(todo.Id);
	if err != nil {
        fmt.Println("Error deleting todo: " + err.Error())
		return c.SendStatus(500)
	}

    return c.Status(200).SendString("")
}

func todoPostHandler(c *fiber.Ctx) error {
	task := c.FormValue("task", "")
	if task == "" {
		return t.AddTodoForm(t.AddTodoFormState{ErrorString: "Please specify a task to do"}).Render(context.Background(), c.Response().BodyWriter())
	}

	todo, err := todoDatabase.AddTodo(task, "")
	if err != nil {
		fmt.Println("Error adding todo: " + err.Error())
		return t.AddTodoForm(t.AddTodoFormState{ErrorString: "Internal error"}).Render(context.Background(), c.Response().BodyWriter())
	}

	ctx := context.Background()
	w := c.Response().BodyWriter()

	t.OOBTodo(todo).Render(ctx, w)
	return t.AddTodoForm(t.AddTodoFormState{}).Render(ctx, w)
}
