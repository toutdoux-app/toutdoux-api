package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// TodoListLabel is used by pop to map your todo_list_labels database table to your go code.
type TodoListLabel struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name string `json:"name" db:"name"`

	TodoList   TodoList  `belongs_to:"todo_list"`
	TodoListID uuid.UUID `json:"todo_list_id" db:"todo_list_id"`
}

// String is not required by pop and may be deleted
func (t TodoListLabel) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TodoListLabels is not required by pop and may be deleted
type TodoListLabels []TodoListLabel

// String is not required by pop and may be deleted
func (t TodoListLabels) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TodoListLabel) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "name"},
		// check to see if the todoList really exists:
		&validators.FuncValidator{
			Field:   t.TodoListID.String(),
			Name:    "todo_list_id",
			Message: "%s list does not exist",
			Fn: func() bool {
				var b bool
				todoList := &TodoList{}
				b, err := tx.Where("todo_list_id = ?", t.TodoListID).Exists(todoList)
				if err != nil {
					return false
				}
				return b
			},
		},
		// check to see if the name is already taken:
		&validators.FuncValidator{
			Field:   t.Name,
			Name:    "name",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("name = ? AND todo_list_id = ?", t.Name, t.TodoListID)
				b, err = q.Exists(t)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TodoListLabel) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TodoListLabel) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
