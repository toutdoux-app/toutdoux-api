package actions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/riton/toutdoux/models"
)

func SetTodoListMiddleware(next buffalo.Handler, listIDParamName string) buffalo.Handler {
	return func(c buffalo.Context) error {
		listID := c.Param(listIDParamName)
		if listID == "" {
			response := make(map[string]interface{})
			response["error"] = "empty list id"
			response["success"] = false
			return c.Render(http.StatusBadRequest, r.JSON(response))
		}

		userID := c.Session().Get("current_user_id").(uuid.UUID)
		tx := c.Value("tx").(*pop.Connection)

		todoList := &models.TodoList{}
		err := tx.Where("id = ? AND user_id = ?", strings.ToLower(strings.TrimSpace(listID)), userID).First(todoList)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response := make(map[string]interface{})
				response["error"] = "no such list id"
				response["success"] = false
				return c.Render(http.StatusBadRequest, r.JSON(response))
			}
			response := make(map[string]interface{})
			response["error"] = "fail to query todo list by ID"
			response["success"] = false
			return c.Render(http.StatusInternalServerError, r.JSON(response))
		}

		c.Set("todo_list", todoList)

		return next(c)
	}
}
