package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TodoList is used by pop to map your todo_lists database table to your go code.
type TodoList struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name string `json:"name" db:"name"`

	User   User `belongs_to:"user"`
	UserID uuid.UUID

	TodoEntries TodoEntries    `has_many:"todo_entries" order_by:"priority updated_at asc"`
	Labels      TodoListLabels `has_many:"todo_list_labels" order_by:"name asc"`
}

// String is not required by pop and may be deleted
func (t TodoList) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TodoLists is not required by pop and may be deleted
type TodoLists []TodoList

// String is not required by pop and may be deleted
func (t TodoLists) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TodoList) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TodoList) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TodoList) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
