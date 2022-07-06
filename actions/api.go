package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/riton/toutdoux/models"
)

// HomeHandler is a default handler to serve up
// a home page.
func APIHealthHandler(c buffalo.Context) error {
	healthStatus := make(map[string]interface{})
	healthStatus["success"] = true

	return c.Render(http.StatusOK, r.JSON(healthStatus))
}

func APIGetListByID(c buffalo.Context) error {
	listID := c.Param("listID")
	if listID == "" {
		return c.Render(http.StatusBadRequest, r.String(""))
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID)

	todoList := &models.TodoList{}
	tx := c.Value("tx").(*pop.Connection)

	err := tx.Where(
		"id = ? AND user_id = ?",
		strings.ToLower(strings.TrimSpace(listID)),
		userID.String(),
	).First(todoList)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// couldn't find document
			response := make(map[string]interface{})
			response["error"] = fmt.Sprintf("no such todoList %s", listID)
			response["success"] = false

			return c.Render(http.StatusOK, r.JSON(response))
		}
		return errors.WithStack(err)
	}

	return nil
}
