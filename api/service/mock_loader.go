package service

import (
	"encoding/json"
	"fmt"
	"os"

	"todo/api/database"
	"todo/api/model"
)

const MOCKDATA_FOLDER_PATH = "./data/jsonmocks/"

// Load Database up with JSON Mock Data if Empty
func init() {
	_, err := database.Connect()
	if err == nil {
		todos, err := database.GetAllTodos()
		if err != nil {
			fmt.Println(err)
		}
		if len(todos) <= 0 {
			count := LoadMockTodos()
			fmt.Println("Initialized DB with Mock Todos: ", count)
		}

	} else {
		fmt.Println(err)
	}
}

func LoadMockTodos() int {
	todosPath := fmt.Sprintf("%vtodos.json", MOCKDATA_FOLDER_PATH)
	todos := make([]model.Todo, 0)
	file, err := os.ReadFile(todosPath)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(file), &todos)
	if err != nil {
		fmt.Println(err)
	}
	objectsForDB := make([]interface{}, 0)
	for i, _ := range todos {
		objectsForDB = append(objectsForDB, interface{}(
			todos[i],
		))
	}
	insertedIds, err := database.InsertObjectsInCollection(objectsForDB, "Todos")
	if err != nil {
		fmt.Println(err)
	}
	return len(insertedIds)
}
