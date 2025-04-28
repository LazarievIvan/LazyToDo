package handler

import (
	"LazyToDo/internal/models"
	"LazyToDo/internal/repository"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

// TodoRepository defines repository for manipulating to-do items.
type TodoRepository interface {
	CreateToDo(*models.ToDo) (models.ToDo, error)
	GetToDos(bag *models.ParamsBag) ([]models.ToDo, error)
	GetToDo(id int64) (models.ToDo, error)
	UpdateToDo(updatedItem *models.ToDo, id int64) (models.ToDo, error)
	DeleteToDo(id int64) error
}

// TodoHandler handles working with ToDoRepository.
type TodoHandler struct {
	repo TodoRepository
}

var createHandler = func() TodoHandler {
	return TodoHandler{repo: repository.NewToDoRepo()}
}

// AddToDo processes request for adding to-do items to DB.
func AddToDo(c *gin.Context) {
	body := readRequestBody(c)

	item, err := models.FromJson(body)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Failed to process JSON",
				"error":   err.Error(),
			})
		return
	}

	handler := createHandler()
	item, err = handler.repo.CreateToDo(&item)
	if err != nil {
		var dbError *models.DBError
		if errors.As(err, &dbError) {
			c.JSON(dbError.Code(), gin.H{"message": dbError.Error(), "error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed creating To-Do item", "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added", "item": item})
}

// GetAllToDos processes request for getting all to-do items from DB.
func GetAllToDos(c *gin.Context) {
	handler := createHandler()

	params := aggregateParams(c)
	todos, err := handler.repo.GetToDos(params)

	if err != nil {
		var dbError *models.DBError
		if errors.As(err, &dbError) {
			c.JSON(dbError.Code(), gin.H{"message": dbError.Error(), "error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed getting To-Do items", "error": err.Error()})
		}
		return
	}
	if len(todos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No To-Do items found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Got them all", "items": todos})
}

// GetSingleToDo processes request for getting single to-do item by given id from params.
func GetSingleToDo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error processing request", "error": err.Error()})
		return
	}
	if id < 1 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid request",
				"error":   fmt.Sprintf("Invalid id: %s", c.Param("id")),
			})
		return
	}

	handler := createHandler()
	item, err := handler.repo.GetToDo(int64(id))
	if err != nil {
		var dbError *models.DBError
		if errors.As(err, &dbError) {
			c.JSON(dbError.Code(), gin.H{"message": dbError.Error(), "error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed getting To-Do item", "id": id, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Retrieved item", "item": item})
}

// UpdateToDo processes request for updating single to-do item with given id from params.
func UpdateToDo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error processing request", "error": err.Error()})
		return
	}
	if id < 1 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid request",
				"error":   fmt.Sprintf("Invalid id: %s", c.Param("id")),
			})
		return
	}

	body := readRequestBody(c)

	item, err := models.FromJson(body)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Failed to process JSON",
				"error":   err.Error(),
			})
		return
	}

	handler := createHandler()
	item, err = handler.repo.UpdateToDo(&item, int64(id))
	if err != nil {
		var dbError *models.DBError
		if errors.As(err, &dbError) {
			c.JSON(dbError.Code(), gin.H{"message": dbError.Error(), "error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed updating To-Do item", "id": id, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated item", "item": item})
}

// DeleteToDo processes request for deleting single to-do item from DB by given id from params.
func DeleteToDo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Error processing request",
				"error":   err.Error(),
			})
		return
	}
	if id < 1 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid request",
				"error":   fmt.Sprintf("Invalid id: %s", c.Param("id")),
			})
		return
	}

	handler := createHandler()
	err = handler.repo.DeleteToDo(int64(id))

	if err != nil {
		var dbError *models.DBError
		if errors.As(err, &dbError) {
			c.JSON(dbError.Code(), gin.H{"message": dbError.Error(), "error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed deleting To-Do item", "id": id, "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted", "ID": id})
}

func readRequestBody(c *gin.Context) []byte {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to read/close the request body: %v", err)
		}
	}(c.Request.Body)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
	}
	return body
}

func aggregateParams(c *gin.Context) *models.ParamsBag {
	return &models.ParamsBag{
		Sort:   extractSortingParams(c),
		Filter: extractFilterParams(c),
		Paging: extractPaginationParams(c),
	}
}

func extractSortingParams(c *gin.Context) models.SortParams {
	orderBy := c.Query("orderBy")
	ascOrdering := c.Query("ASC")
	asc := true
	switch ascOrdering {
	case "true":
		asc = true
		break
	case "false":
		asc = false
		break
	default:
		asc = true
		break
	}
	return models.SortParams{Field: orderBy, ASC: asc}
}

func extractFilterParams(c *gin.Context) models.FilterParams {
	statusFilter := c.Query("status")
	var filters []models.Filter
	if statusFilter != "" {
		filters = append(filters, models.Filter{Field: "status", Value: statusFilter})
	}
	return models.FilterParams{Filters: filters}
}

func extractPaginationParams(c *gin.Context) models.PaginationParams {
	limitParam := c.Query("limit")
	pageParam := c.Query("page")
	// If no limit - don't apply pagination.
	if len(limitParam) == 0 {
		return models.PaginationParams{}
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		fmt.Println(err)
		limit = 0
	}
	// If no page specified, apply limit without offset.
	if len(pageParam) == 0 {
		return models.PaginationParams{Limit: limit}
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		fmt.Println(err)
		page = 0
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return models.PaginationParams{Limit: limit, Offset: offset}
}
