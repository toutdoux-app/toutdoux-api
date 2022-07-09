package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

var initialDataInjectors = []func(*pop.Connection) error{
	injectInitialTodoEntryRelationTypes,
}

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	for _, initialDataInjector := range initialDataInjectors {
		DB.Transaction(func(tx *pop.Connection) error {
			return initialDataInjector(tx)
		})
	}
}
