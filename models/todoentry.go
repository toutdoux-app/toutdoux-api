package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TodoEntry is used by pop to map your todo_entries database table to your go code.
type TodoEntry struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	TodoList   TodoList  `belongs_to:"todo_list" json:"-"`
	TodoListID uuid.UUID `db:"todo_list_id" json:"-"`

	Title    string         `json:"title" db:"title"`
	Priority int            `json:"priority" db:"priority"`
	Done     bool           `json:"done" db:"done"`
	Labels   TodoListLabels `many_to_many:"todo_entry_labels" json:"labels"`
}

// String is not required by pop and may be deleted
func (t TodoEntry) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TodoEntries is not required by pop and may be deleted
type TodoEntries []TodoEntry

// String is not required by pop and may be deleted
func (t TodoEntries) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TodoEntry) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TodoEntry) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TodoEntry) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
