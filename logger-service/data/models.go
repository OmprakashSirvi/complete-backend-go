package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
/*
 Return a new empty Models object with empty [LogEntry] object
 */
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	logsCollection := client.Database("logs").Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	result, err := logsCollection.InsertOne(ctx, LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("There was some error while inserting entry %v", err)
		return err
	}
	log.Printf("The result from collection insert one : %v", result)

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logsCollection := client.Database("logs").Collection("logs")

	// Find and sort by "created_at"
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cur, err := logsCollection.Find(ctx, bson.D{}, opts)

	if err != nil {
		log.Printf("There is some error while getting logs : %v", err)
		return nil, err
	}

	// Close the cursor
	defer cur.Close(ctx)

	var logs []*LogEntry

	for cur.Next(ctx) {
		var item LogEntry

		// Decode the cursor into the item
		err = cur.Decode(&item)

		if err != nil {
			log.Printf("Error while decoding the cursor to LogEntry type object : %v", err)
			return nil, err
		} else {
			// Append the logs with the item
			logs = append(logs, &item)
		}

	}

	return logs, nil

}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	// Create a new logs collection which refers to the collection (same as tables in SQL)

	logsCollection := client.Database("logs").Collection("logs")

	logId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("Error while converting id to ObjectId : %v", err)
		return nil, err
	}

	var entry LogEntry

	// The error in the returned singleResult object will be returned by the Decode(if any)
	err = logsCollection.FindOne(ctx, bson.M{"_id": logId}).Decode(&entry)

	if err != nil {
		log.Printf("Error while finding the log by id : %v", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	logsCollection := client.Database("logs").Collection("logs")
	err := logsCollection.Drop(ctx)
	if err != nil {
		log.Printf("Error while dropping logs : %v", err)
		return err
	}

	return nil

}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	logsCollection := client.Database("logs").Collection("logs")

	primitive.ObjectIDFromHex(l.ID)

	logId, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		log.Printf("Error while converting id to ObjectId : %v", err)
		return nil, err
	}

	upRes, err := logsCollection.UpdateByID(ctx, bson.M{"_id": logId}, 
	bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: l.Name},
			{Key: "data", Value: l.Data},
		}},
	})

	if err != nil {
		log.Printf("Error while updating the logs : %v", err)
		return nil, err
	}

	return upRes, nil

}
