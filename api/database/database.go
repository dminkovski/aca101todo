package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"todo/api/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	Client *mongo.Client
}

var DATABASE_URI = constants.DB_CONNECTION_URI

func init() {
	value, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		DATABASE_URI = value
	} else {
		fmt.Println("Using Localhost as Database since no DATABASE_URI env provided")
	}
}

func Connect() (*DB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DATABASE_URI))
	if err != nil {
		fmt.Println("Database Connect failed")
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return &DB{
		Client: client,
	}, nil
}

func (db *DB) Disconnect() error {
	client := db.Client
	if client != nil {
		return client.Disconnect(context.TODO())
	}
	return errors.New("no client to disconnect")
}

func (db *DB) GetDatabase() *mongo.Database {
	return db.Client.Database("db")
}

func (db *DB) GetClient() *mongo.Client {
	return db.Client
}

func InsertObjectsInCollection(data []interface{}, collectionKey string) ([]primitive.ObjectID, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	col := db.GetDatabase().Collection(collectionKey)

	result, err := col.InsertMany(context.TODO(), data)
	if err != nil {
		return nil, err
	}

	indexModel := []mongo.IndexModel{{
		Keys: bson.D{
			{"title", "text"},
		},
	},
	}

	_, err = col.Indexes().CreateMany(context.TODO(), indexModel)
	if err != nil {
		return nil, err
	}
	err = db.Disconnect()
	if err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, 0)
	for _, id := range result.InsertedIDs {
		val, ok := id.(primitive.ObjectID)
		if ok {
			ids = append(ids, primitive.ObjectID(val))
		}
	}
	return ids, nil
}

func GetAllObjectsInCollection(collectionKey string, filter interface{}, opts *options.FindOptions) (*mongo.Cursor, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	col := db.GetDatabase().Collection(collectionKey)
	cursor, err := col.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	err = db.Disconnect()
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func GetObjectIDFromString(id string) (primitive.ObjectID, error) {
	var nilObject primitive.ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nilObject, err
	}
	return objectId, nil
}
