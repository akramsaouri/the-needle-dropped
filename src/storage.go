package main

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// Storage ...
type Storage struct {
	database *mongo.Database
}

// Connect to mongodb source
func (storage *Storage) Connect(uri, databaseName string) error {
	client, err := mongo.Connect(context.Background(), uri)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}
	storage.database = client.Database(databaseName)
	return err
}

// InsertDoc ...
func (storage *Storage) InsertDoc(collectionName string, doc interface{}) error {
	collection := storage.database.Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), doc)
	return err
}

// ViewDoc ...
func (storage *Storage) ViewDoc(collectionName string, filter interface{}) (video Video, err error) {
	collection := storage.database.Collection(collectionName)
	ctx := context.Background()
	err = collection.FindOneAndDelete(ctx, filter).Decode(&video)
	if err != nil {
		return video, err
	}
	return video, nil
}

// ListDocs ...
func (storage *Storage) ListDocs(collectionName string) ([]Video, error) {
	videos := make([]Video, 0)
	collection := storage.database.Collection(collectionName)
	ctx := context.Background()
	cur, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var video Video
		err = cur.Decode(&video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, err
}
