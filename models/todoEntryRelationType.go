package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

var initialRelationTypes = map[string]string{
	"related to": "related to",
	"blocked by": "blocking",
	"blocking":   "blocked by",
	"precedes":   "follows",
	"follows":    "precedes",
}

func injectInitialTodoEntryRelationTypes(tx *pop.Connection) error {
	for name, reverseName := range initialRelationTypes {
		found, err := tx.Where("name = ?", name).Exists(&TodoEntryRelationType{})
		if err != nil {
			log.Fatalf("fail to search for todo_entry_relation_type %s: %s", name, err)
		}

		if found {
			continue
		}

		rt := TodoEntryRelationType{
			Name:        name,
			ReverseName: reverseName,
		}

		if err := tx.Create(&rt); err != nil {
			log.Fatalf("fail to create todo_entry_relation_type %s: %s", name, err)
		}
	}
	return nil
}

// TodoEntryRelationType is used by pop to map your todo_entry_relation_types database table to your go code.
type TodoEntryRelationType struct {
	ID        uuid.UUID `json:"-" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// Name and ReverseName are things such as:
	// Name=blocked by
	// ReverseName=blocking
	Name        string `json:"name" db:"name"`
	ReverseName string `json:"reverse_name" db:"reverse_name"`

	TodoEntryRelations TodoEntryRelations `has_many:"todo_entry_relations" json:"-"`
}

// String is not required by pop and may be deleted
func (t TodoEntryRelationType) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TodoEntryRelationTypes is not required by pop and may be deleted
type TodoEntryRelationTypes []TodoEntryRelationType

// String is not required by pop and may be deleted
func (t TodoEntryRelationTypes) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TodoEntryRelationType) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	vErr := validate.Validate(
		&validators.FuncValidator{
			Field:   t.Name,
			Name:    "Name",
			Message: "%s is already taken",
			Fn: func() bool {
				var found bool
				found, err = tx.Where("name = ?", t.Name).Exists(&TodoEntryRelationType{})
				if err != nil {
					return false
				}

				return !found
			},
		},
	)
	return vErr, err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TodoEntryRelationType) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TodoEntryRelationType) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
