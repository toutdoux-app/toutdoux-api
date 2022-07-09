package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotchance/pie/v2"
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

type APITodoListEntriesCreateRelatedTodoRequest struct {
	RelatedTodoID string `json:"related_todo"`
	RelationType  string `json:"relation_type"`
}

type APITodoListEntriesCreateRequest struct {
	Priority  int                                          `json:"priority"`
	Title     string                                       `json:"title"`
	Labels    []string                                     `json:"labels"`
	Relations []APITodoListEntriesCreateRelatedTodoRequest `json:"relations"`
}

func (a APITodoListEntriesCreateRequest) Validate() error {
	if a.Priority < 0 {
		return fmt.Errorf("priority MUST be positive integer")
	}

	if a.Title == "" {
		return fmt.Errorf("title can't be empty")
	}

	return nil
}

func (a APITodoListEntriesCreateRequest) GetLabels() []string {
	return pie.Unique(a.Labels)
}

func APITodoListEntriesCreate(c buffalo.Context) error {
	var request APITodoListEntriesCreateRequest
	if err := c.Bind(&request); err != nil {
		response := make(map[string]interface{})
		response["success"] = false
		response["error"] = fmt.Sprintf("fail to process arguments: %s", err)
		return c.Render(http.StatusBadRequest, r.JSON(response))
	}

	if err := request.Validate(); err != nil {
		response := make(map[string]interface{})
		response["success"] = false
		response["error"] = err.Error()
		return c.Render(http.StatusBadRequest, r.JSON(response))
	}

	userID := c.Session().Get("current_user_id").(uuid.UUID)
	todoList := c.Value("todo_list").(*models.TodoList)
	tx := c.Value("tx").(*pop.Connection)

	// TodoList labels will be created "on the fly" if it does not exist
	var todoListLabels models.TodoListLabels
	for _, label := range request.GetLabels() {
		var shouldCreateLabel bool
		todoListLabel := &models.TodoListLabel{}

		err := tx.Where("name = ? AND todo_list_id = ?", label, todoList.ID).First(todoListLabel)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shouldCreateLabel = true
			} else {
				response := make(map[string]interface{})
				response["success"] = false
				response["error"] = fmt.Sprintf("fail to query for label %q", label)
				return c.Render(http.StatusInternalServerError, r.JSON(response))
			}
		}

		if shouldCreateLabel {
			todoListLabel = &models.TodoListLabel{
				Name:       label,
				TodoListID: todoList.ID,
			}
			todoListLabel.SetUserID(userID)

			vErr, err := tx.ValidateAndCreate(todoListLabel)
			if err != nil {
				response := make(map[string]interface{})
				response["success"] = false
				response["error"] = fmt.Sprintf("fail to create label %q", label)
				return c.Render(http.StatusInternalServerError, r.JSON(response))
			}

			if len(vErr.Errors) > 0 {
				response := make(map[string]interface{})
				response["error"] = fmt.Sprintf("invalid values when creating label %q", label)
				response["detail"] = vErr.Errors
				response["success"] = false

				return c.Render(http.StatusBadRequest, r.JSON(response))
			}
		}

		todoListLabels = append(todoListLabels, *todoListLabel)
	}

	todoListEntry := &models.TodoEntry{
		Title:      request.Title,
		TodoListID: todoList.ID,
		Labels:     todoListLabels,
		Done:       false,
		Priority:   request.Priority,
	}

	if err := tx.Create(todoListEntry); err != nil {
		c.Logger().WithFields(map[string]interface{}{
			"error":      err,
			"request_id": c.Value("request_id"),
			"user_id":    userID.String(),
		}).Error("fail to create todo list entry")

		response := make(map[string]interface{})
		response["success"] = false
		response["error"] = "fail to create todo list entry"
		return c.Render(http.StatusInternalServerError, r.JSON(response))
	}

	// now that the TodoEntry is created, we have access to the generated ID

	var relations models.TodoEntryRelations
	for _, relation := range request.Relations {
		mRelation := models.TodoEntryRelation{
			TodoEntryID:  todoListEntry.ID,
			RelationType: relation.RelationType,
		}

		relatedTodoID, err := uuid.FromString(relation.RelatedTodoID)
		if err != nil {
			response := make(map[string]interface{})
			response["success"] = false
			response["error"] = fmt.Sprintf("invalid related_todo_id %s", relation.RelatedTodoID)
			return c.Render(http.StatusBadRequest, r.JSON(response))
		}

		mRelation.RelatedToTodoEntryID = relatedTodoID

		relations = append(relations, mRelation)
	}

	vErr, err := tx.ValidateAndCreate(relations)
	if err != nil {
		c.Logger().WithFields(map[string]interface{}{
			"error":      err,
			"request_id": c.Value("request_id"),
			"user_id":    userID.String(),
		}).Error("fail to create todo list entry relations")

		response := make(map[string]interface{})
		response["success"] = false
		response["error"] = "fail to create todo list entry relations"
		return c.Render(http.StatusInternalServerError, r.JSON(response))
	}

	if len(vErr.Errors) > 0 {
		response := make(map[string]interface{})
		response["error"] = "invalid values when creating todo entry relations"
		response["detail"] = vErr.Errors
		response["success"] = false

		return c.Render(http.StatusBadRequest, r.JSON(response))
	}

	// if err := tx.EagerPreload().Find(&models.TodoEntry{}, todoListEntry.ID); err != nil {
	// 	c.Logger().WithFields(map[string]interface{}{
	// 		"error":      err,
	// 		"request_id": c.Value("request_id"),
	// 		"user_id":    userID.String(),
	// 	}).Error("fail to load todo list entry relations")

	// 	response := make(map[string]interface{})
	// 	response["success"] = false
	// 	response["error"] = "fail to load todo list entry relations"
	// 	return c.Render(http.StatusInternalServerError, r.JSON(response))
	// }

	var tdList models.TodoEntry
	err = tx.EagerPreload().Find(&tdList, todoListEntry.ID)
	if err != nil {
		c.Logger().WithFields(map[string]interface{}{
			"error":      err,
			"request_id": c.Value("request_id"),
			"user_id":    userID.String(),
		}).Error("fail to load created todo list entry")

		response := make(map[string]interface{})
		response["success"] = false
		response["error"] = "fail to load created todo list entry"
		return c.Render(http.StatusInternalServerError, r.JSON(response))
	}

	return c.Render(http.StatusOK, r.JSON(tdList))
}
