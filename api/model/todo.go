package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `json:"title"`
	Updated time.Time          `json:"updated"`
	Created time.Time          `json:"created"`
	Done    bool               `json:"done"`
}

func (todo Todo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":      todo.ID,
		"title":   todo.Title,
		"created": todo.Created,
		"done":    todo.Done,
	})
}

func (todo *Todo) Create(title string) {
	todo.Title = title
	todo.Created = time.Now()
	todo.Updated = time.Now()
}
