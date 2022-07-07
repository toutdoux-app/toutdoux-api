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
	todoList := c.Value("todo_list").(*models.TodoList)
	return c.Render(http.StatusOK, r.JSON(todoList))
}

type APIListCreateRequest struct {
	Name string `json:"name" form:"name"`
}

func APIListCreate(c buffalo.Context) error {
	request := &APIListCreateRequest{}
	if err := c.Bind(request); err != nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf(""))
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID)

	// See https://andrew-sledge.gitbooks.io/the-unofficial-pop-book/content/common-patterns/creating-new-records.html
	todoList := &models.TodoList{
		Name:   strings.TrimSpace(strings.ToLower(request.Name)),
		UserID: userID,
	}

	tx := c.Value("tx").(*pop.Connection)

	_, err := tx.ValidateAndCreate(todoList)
	if err != nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf(""))
	}

	return c.Redirect(http.StatusFound, "/api/list/"+todoList.ID.String())
}

func APIListsAll(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	userID := c.Session().Get("current_user_id").(uuid.UUID)

	// https://andrew-sledge.gitbooks.io/the-unofficial-pop-book/content/common-patterns/querying-for-several-records.html
	lists := &models.TodoLists{}
	err := tx.Select("id", "name", "created_at", "updated_at", "user_id").Where("user_id = ?", userID.String()).All(lists)
	if err != nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf(""))
	}

	return c.Render(http.StatusOK, r.JSON(lists))
}

type APIListLabelCreateRequest struct {
	Name string `json:"name" form:"name"`
}

func APIListLabelCreate(c buffalo.Context) error {
	request := &APIListLabelCreateRequest{}
	if err := c.Bind(request); err != nil {
		return c.Error(http.StatusBadRequest, fmt.Errorf(""))
	}

	todoList := c.Value("todo_list").(*models.TodoList)
	tx := c.Value("tx").(*pop.Connection)

	label := &models.TodoListLabel{
		Name:       request.Name,
		TodoListID: todoList.ID,
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID)
	label.SetUserID(userID)

	vErr, err := tx.ValidateAndCreate(label)
	if err != nil {
		c.Logger().Error(err)
		return c.Error(http.StatusInternalServerError, fmt.Errorf(""))
	}

	if len(vErr.Errors) > 0 {
		response := make(map[string]interface{})
		response["error"] = "invalid values"
		response["detail"] = vErr.Errors
		response["success"] = false

		return c.Render(http.StatusBadRequest, r.JSON(response))
	}

	response := make(map[string]interface{})
	response["success"] = true
	response["label"] = label

	return c.Render(http.StatusOK, r.JSON(response))
}

func APITodoListEntriesList(c buffalo.Context) error {
	todoList := c.Value("todo_list").(*models.TodoList)
	tx := c.Value("tx").(*pop.Connection)

	q := tx.Where("todo_list_id = ?", todoList.ID)
	if doneParam := c.Param("done"); doneParam != "" {
		if doneParam == "true" || doneParam == "false" {
			q.Where("done = " + doneParam)
		}
	}

	todoEntries := &models.TodoEntries{}
	err := q.All(todoEntries)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			response := make(map[string]interface{})
			response["success"] = false
			return c.Render(http.StatusInternalServerError, r.JSON(response))
		}
	}

	return c.Render(http.StatusOK, r.JSON(todoEntries))
}
