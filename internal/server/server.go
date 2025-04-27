package server

import (
	"LazyToDo/internal/handler"
	"github.com/gin-gonic/gin"
)

func Start(port string) error {
	r := gin.Default()
	handler.Route(r)
	err := r.Run(":" + port)
	if err != nil {
		return err
	}
	return nil
}
