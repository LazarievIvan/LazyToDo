package handler

import (
	_ "LazyToDo/cmd/todo/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Route endpoints to handler.
func Route(r *gin.Engine) {
	r.POST("/add", AddToDo)
	r.GET("/todos", GetAllToDos)
	r.GET("/todos/:id", GetSingleToDo)
	r.PUT("/todos/:id", UpdateToDo)
	r.DELETE("/todos/:id", DeleteToDo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/static/swagger.yaml")))
	r.Static("/static", "/app/cmd/todo/docs")
}
