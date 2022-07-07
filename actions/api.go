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

	listID := c.Param("listID")
	if listID == "" {
		return c.Render(http.StatusBadRequest, r.String(""))
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID)
	tx := c.Value("tx").(*pop.Connection)

	listUUID, err := uuid.FromString(listID)
	if err != nil {
		response := make(map[string]interface{})
		response["error"] = "invalid todo_list identifier"
		response["success"] = false

		c.Logger().WithFields(map[string]interface{}{
			"error":      err,
			"listID":     listID,
			"request_id": c.Value("request_id"),
		}).Warn("invalid todo_list identifier")

		return c.Render(http.StatusBadRequest, r.JSON(response))
	}

	label := &models.TodoListLabel{
		Name:       request.Name,
		TodoListID: listUUID,
	}
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
