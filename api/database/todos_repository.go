package database

import (
	"context"
	"errors"
	"time"

	"todo/api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TODOS_COLLECTION_KEY = "Todos"

func GetAllTodos() ([]model.Todo, error) {
	todos := make([]model.Todo, 0)
	cursor, err := GetAllObjectsInCollection(TODOS_COLLECTION_KEY, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var todo model.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func CreateTodo(title string) (model.Todo, error) {
	var todo model.Todo
	todo.Create(title)
	data := []interface{}{todo}
	insertedIds, err := InsertObjectsInCollection(data, TODOS_COLLECTION_KEY)
	if err != nil || len(insertedIds) <= 0 {
		return todo, err
	}
	insertedId := insertedIds[0]
	todo.ID = insertedId
	return todo, nil
}

func GetTodoById(todoID string) (model.Todo, error) {
	var todo model.Todo
	objectId, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return todo, err
	}
	db, err := Connect()
	if err != nil {
		return todo, err
	}
	col := db.GetDatabase().Collection(TODOS_COLLECTION_KEY)

	result := col.FindOne(context.TODO(), bson.M{
		"_id": objectId,
	})
	if err != nil {
		return todo, err
	}
	if result == nil {
		return todo, errors.New("Todo not found")
	}

	err = db.Disconnect()
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func UpdateTodo(todoID string, done bool) error {
	objectId, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return err
	}
	todo, err := GetTodoById(todoID)
	if err != nil {
		return err
	}
	todo.Done = done
	db, err := Connect()
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.M{"done": done, "updated": time.Now()}}}
	col := db.GetDatabase().Collection(TODOS_COLLECTION_KEY)
	col.UpdateOne(
		context.TODO(),
		filter,
		update,
	)

	err = db.Disconnect()
	if err != nil {
		return err
	}
	return nil
}
