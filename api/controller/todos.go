package controller

import (
	"fmt"
	"net/http"

	"todo/api/database"

	"github.com/gin-gonic/gin"
)

func UpdateTodo(c *gin.Context) {
	type UpdateTodoRequest struct {
		Done bool `json:"done"`
	}
	var jsonRequest UpdateTodoRequest
	if err := c.ShouldBindJSON(&jsonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(jsonRequest)
	todoId := c.Param("id")
	err := database.UpdateTodo(todoId, jsonRequest.Done)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo could not be updated"})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Successfully updated"})
		return
	}
}

func GetAllTodos(c *gin.Context) {
	todos, err := database.GetAllTodos()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todos not found"})
	}
	c.IndentedJSON(http.StatusOK, todos)
}

func GetTodoById(c *gin.Context) {
	todoId := c.Param("id")
	todo, err := database.GetTodoById(todoId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
	}
	c.IndentedJSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	type CreateTodoRequest struct {
		Title string `json:"title"`
	}
	var jsonRequest CreateTodoRequest
	if err := c.ShouldBindJSON(&jsonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := database.CreateTodo(jsonRequest.Title)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo could not be created"})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, todo)
		return
	}
}
