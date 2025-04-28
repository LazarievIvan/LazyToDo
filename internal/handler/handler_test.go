package handler

import (
	"LazyToDo/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const DummyId = 1

// mockRepo implements interface TodoRepository
type mockRepo struct {
	Error       error
	ReturnValue models.ToDo
}

func (m *mockRepo) CreateToDo(item *models.ToDo) (models.ToDo, error) {
	if m.Error != nil {
		return models.ToDo{}, m.Error
	}
	return *item, nil
}

func (m *mockRepo) GetToDos(bag *models.ParamsBag) ([]models.ToDo, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	if m.ReturnValue == (models.ToDo{}) {
		return []models.ToDo{}, m.Error
	}
	return []models.ToDo{m.ReturnValue}, nil
}

func (m *mockRepo) GetToDo(id int64) (models.ToDo, error) {
	if m.Error != nil {
		return models.ToDo{}, m.Error
	}
	return m.ReturnValue, nil
}

func (m *mockRepo) UpdateToDo(item *models.ToDo, id int64) (models.ToDo, error) {
	if m.Error != nil {
		return models.ToDo{}, m.Error
	}
	return m.ReturnValue, nil
}

func (m *mockRepo) DeleteToDo(id int64) error {
	if m.Error != nil {
		return m.Error
	}
	return nil
}

// TestAddToDo covers all possible cases of adding to-do with respective return statuses.
func TestAddToDo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		requestBody        string
		mockError          error
		expectedStatusCode int
	}{
		{
			name:               "CreateToDo returns BadRequest",
			requestBody:        `{"invalid json"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "CreateToDo returns InternalServerError",
			requestBody:        `{"description": "Description", "status": "TO DO"}`,
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "CreateToDo returns OK",
			requestBody:        `{"description": "Description", "status": "TO DO"}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/add", strings.NewReader(test.requestBody))

			createHandlerMethod := createHandler
			createHandler = func() TodoHandler {
				return TodoHandler{repo: &mockRepo{Error: test.mockError}}
			}

			t.Cleanup(func() {
				createHandler = createHandlerMethod
			})

			AddToDo(c)
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

// TestGetAllToDos covers all possible cases of getting all to-do with respective return statuses.
func TestGetAllToDos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		mockError          error
		returnValue        models.ToDo
		expectedStatusCode int
	}{
		{
			name:               "GetToDos return InternalServerError",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "GetToDos return NotFound",
			mockError:          nil,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "GetToDos return OK",
			mockError:          nil,
			returnValue:        models.ToDo{ID: DummyId},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/todos", nil)

			createHandlerMethod := createHandler
			createHandler = func() TodoHandler {
				return TodoHandler{repo: &mockRepo{Error: test.mockError, ReturnValue: test.returnValue}}
			}

			t.Cleanup(func() {
				createHandler = createHandlerMethod
			})

			GetAllToDos(c)
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

// TestGetSingleToDo covers all possible cases of getting single to-do with respective return statuses.
func TestGetSingleToDo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name               string
		requestParam       string
		mockError          error
		returnValue        models.ToDo
		expectedStatusCode int
	}{
		{
			name:               "GetSingleToDo returns BadRequest for invalid ID param",
			requestParam:       "string",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "GetSingleToDo returns BadRequest for invalid ID",
			requestParam:       "0",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "GetSingleToDo returns NotFound",
			requestParam:       strconv.Itoa(DummyId),
			mockError:          models.NewDBError("Not Found", http.StatusNotFound, errors.New("something went wrong")),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "GetSingleToDo returns OK",
			requestParam:       strconv.Itoa(DummyId),
			mockError:          nil,
			returnValue:        models.ToDo{ID: DummyId},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/todos/"+test.requestParam, nil)
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: test.requestParam},
			}

			createHandlerMethod := createHandler
			createHandler = func() TodoHandler {
				return TodoHandler{repo: &mockRepo{Error: test.mockError, ReturnValue: test.returnValue}}
			}

			t.Cleanup(func() {
				createHandler = createHandlerMethod
			})

			GetSingleToDo(c)
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

// TestUpdateToDo covers all possible cases of updating single to-do with respective return statuses.
func TestUpdateToDo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name               string
		requestParam       string
		requestBody        string
		mockError          error
		returnValue        models.ToDo
		expectedStatusCode int
	}{
		{
			name:               "UpdateToDo returns BadRequest with invalid ID type",
			requestParam:       "string",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "UpdateToDo returns BadRequest with invalid ID",
			requestParam:       "0",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "UpdateToDo returns BadRequest with wrong body",
			requestParam:       strconv.Itoa(DummyId),
			requestBody:        `{"invalid json"}`,
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "UpdateToDo returns NotFound",
			requestParam:       strconv.Itoa(DummyId),
			requestBody:        `{"description": "Description", "status": "TO DO"}`,
			mockError:          models.NewDBError("Not Found", http.StatusNotFound, errors.New("something went wrong")),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "UpdateToDo returns InternalServerError",
			requestParam:       strconv.Itoa(DummyId),
			requestBody:        `{"description": "Description", "status": "TO DO"}`,
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "UpdateToDo returns OK",
			requestParam:       strconv.Itoa(DummyId),
			requestBody:        `{"description": "Description", "status": "TO DO"}`,
			mockError:          nil,
			returnValue:        models.ToDo{ID: DummyId},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPut, "/todos/"+test.requestParam, strings.NewReader(test.requestBody))
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: test.requestParam},
			}

			createHandlerMethod := createHandler
			createHandler = func() TodoHandler {
				return TodoHandler{repo: &mockRepo{Error: test.mockError, ReturnValue: test.returnValue}}
			}
			t.Cleanup(func() {
				createHandler = createHandlerMethod
			})

			UpdateToDo(c)
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

// TestDeleteToDo covers all possible cases of deleting single to-do with respective return statuses.
func TestDeleteToDo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name               string
		requestParam       string
		mockError          error
		expectedStatusCode int
	}{
		{
			name:               "DeleteToDo returns BadRequest with invalid ID type",
			requestParam:       "string",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "DeleteToDo returns BadRequest with invalid ID",
			requestParam:       "0",
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "DeleteToDo returns NotFound",
			requestParam:       strconv.Itoa(DummyId),
			mockError:          models.NewDBError("Not Found", http.StatusNotFound, errors.New("something went wrong")),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "DeleteToDo returns InternalServerError",
			requestParam:       strconv.Itoa(DummyId),
			mockError:          errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "DeleteToDo returns OK",
			requestParam:       strconv.Itoa(DummyId),
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodDelete, "/todos/"+test.requestParam, nil)
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: test.requestParam},
			}

			createHandlerMethod := createHandler
			createHandler = func() TodoHandler {
				return TodoHandler{repo: &mockRepo{Error: test.mockError}}
			}
			t.Cleanup(func() {
				createHandler = createHandlerMethod
			})

			DeleteToDo(c)
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
