package handler

import "github.com/gin-gonic/gin"

// Route endpoints to handler.
func Route(r *gin.Engine) {
	r.POST("/add", AddToDo)
	r.GET("/todos", GetAllToDos)
	r.GET("/todos/:id", GetSingleToDo)
	r.PUT("/todos/:id", UpdateToDo)
	r.DELETE("/todos/:id", DeleteToDo)
}
