package main

import (
	"net/http"

	"todo/api/constants"
	"todo/api/controller"
	_ "todo/api/database"
	_ "todo/api/model"
	_ "todo/api/service"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "I am alive"})
	})

	router.GET("/todos", controller.GetAllTodos)
	router.GET("/todos/:id", controller.GetTodoById)
	router.POST("/todos/:id", controller.UpdateTodo)
	router.POST("/todos", controller.CreateTodo)

	router.Static("/assets/img", "./assets/img/")

	router.Run(constants.APPLICATION_PORT)
}
