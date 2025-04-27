package handler

import (
	"LazyToDo/internal/models"
	"LazyToDo/internal/repository"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

// TodoRepository defines repository for manipulating to-do items.
type TodoRepository interface {
	CreateToDo(*models.ToDo) (models.ToDo, error)
	GetToDos() ([]models.ToDo, error)
	GetToDo(id int64) (models.ToDo, error)
	UpdateToDo(updatedItem *models.ToDo, id int64) (models.ToDo, error)
	DeleteToDo(id int64) error
}

// TodoHandler handles working with ToDoRepository.
type TodoHandler struct {
	repo TodoRepository
}

func createHandler() TodoHandler {
	return TodoHandler{repo: repository.NewToDoRepo()}
}

// AddToDo processes request for adding to-do items to DB.
func AddToDo(c *gin.Context) {
	body := readRequestBody(c)

	item, err := models.FromJson(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process JSON", "error": err.Error()})
		return
	}

	handler := createHandler()
	item, err = handler.repo.CreateToDo(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed creating To-Do item", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added", "item": item})
}

// GetAllToDos processes request for getting all to-do items from DB.
func GetAllToDos(c *gin.Context) {
	handler := createHandler()
	todos, err := handler.repo.GetToDos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed getting To-Do items", "error": err.Error()})
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

	handler := createHandler()
	item, err := handler.repo.GetToDo(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed getting To-Do item", "id": id, "error": err.Error()})
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

	body := readRequestBody(c)

	item, err := models.FromJson(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process JSON", "error": err.Error()})
		return
	}

	handler := createHandler()
	item, err = handler.repo.UpdateToDo(&item, int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed updating To-Do item", "id": id, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update item", "item": item})
}

// DeleteToDo processes request for deleting single to-do item from DB by given id from params.
func DeleteToDo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error processing request", "error": err.Error()})
		return
	}

	handler := createHandler()
	err = handler.repo.DeleteToDo(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed deleting To-Do item", "id": id, "error": err.Error()})
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
