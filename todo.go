package main

// TodoItem represents a todo item
type TodoItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority" validate:"min:1;max:10"`
}