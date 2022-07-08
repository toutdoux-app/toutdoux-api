package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TodoEntryRelation is used by pop to map your todo_entry_relations database table to your go code.
type TodoEntryRelation struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	TodoEntry            TodoEntry `belongs_to:"todo_entry"`
	TodoEntryID          uuid.UUID `json:"todo_entry_id"`
	RelatedToTodoEntry   TodoEntry `belongs_to:"todo_entry"`
	RelatedToTodoEntryID uuid.UUID `json:"related_to_todo_entry_id"`

	RelationType string `db:"relation_type" json:"relation_type"`
}

// String is not required by pop and may be deleted
func (t TodoEntryRelation) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TodoEntryRelations is not required by pop and may be deleted
type TodoEntryRelations []TodoEntryRelation

// String is not required by pop and may be deleted
func (t TodoEntryRelations) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TodoEntryRelation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TodoEntryRelation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TodoEntryRelation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
