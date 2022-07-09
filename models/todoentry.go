package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// TodoEntry is used by pop to map your todo_entries database table to your go code.
type TodoEntry struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	TodoList   TodoList  `belongs_to:"todo_lists" json:"-"`
	TodoListID uuid.UUID `db:"todo_list_id" json:"-"`

	Title    string         `json:"title" db:"title"`
	Priority int            `json:"priority" db:"priority"`
	DueDate  *time.Time     `json:"due_date,omitempty" db:"due_date"`
	Done     bool           `json:"done" db:"done"`
	Labels   TodoListLabels `many_to_many:"todo_entry_labels" db:"-" json:"labels,omitempty"`

	//Relations TodoEntryRelations `many_to_many:"todo_entry_relations" fk_id:"todo_entry_id" primary_id:"TodoEntryID" json:"relations,omitempty"`
	Relations TodoEntryRelations `many_to_many:"todo_entry_relations" fk_id:"id" db:"-" json:"relations,omitempty"`
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
	vErr := validate.NewErrors()

	for _, todoListLabel := range t.Labels {
		found, err := tx.Where("id = ? AND todo_list_id = ?", todoListLabel.ID, t.TodoListID).Exists(&TodoListLabel{})
		if err != nil {
			return vErr, err
		}

		if !found {
			vErr.Add("label", fmt.Sprintf("todo_list_label %s not found", todoListLabel.ID))
		}
	}

	for _, todoListEntryRelation := range t.Relations {
		if t.ID != todoListEntryRelation.TodoEntry.ID {
			vErr.Add("relation", "related todo list has invalid todo_entry_id. Must be self ID")
			continue
		}

		exists, err := TodoEntryExists(tx, todoListEntryRelation.RelatedToTodoEntryID)
		if err != nil {
			return vErr, errors.Wrapf(err, "checking if related TodoEntry %s existed", todoListEntryRelation.TodoEntry.ID)
		}
		if !exists {
			vErr.Add("relation", fmt.Sprintf("todo_list_entry %s not found", todoListEntryRelation.TodoEntry.ID))
			continue
		}
	}

	vErr.Append(validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title"},
		&validators.IntIsGreaterThan{Field: t.Priority, Name: "Priority", Compared: 0},
	))

	return vErr, nil
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

func TodoEntryExists(tx *pop.Connection, todoEntryID uuid.UUID) (bool, error) {
	found, err := tx.Where("id = ?", todoEntryID).Exists(&TodoEntry{})
	if err != nil {
		return false, err
	}

	return found, nil
}

func TodoEntryIDToTodoListID(tx *pop.Connection, todoID uuid.UUID) (uuid.UUID, error) {
	var todoEntry TodoEntry
	err := tx.Find(&todoEntry, todoID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return todoEntry.TodoListID, nil
}
