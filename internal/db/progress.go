package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProgressRecord struct {
	Blockchain   string                 `bson:"blockchain"`
	Address      string                 `bson:"address"`
	Status       string                 `bson:"status"` // pending, done, failed
	LastUpdated  time.Time              `bson:"last_updated"`
	Retries      int                    `bson:"retries"`
	ResponseData map[string]interface{} `bson:"response_data,omitempty"`
}

type ProgressDB struct {
	collection *mongo.Collection
}

func NewProgressDB(uri, dbName, collectionName string) *ProgressDB {
    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.Background(), clientOpts)
    if err != nil {
        log.Fatalf("MongoDB connection failed: %v\nUsing URI: %s", err, uri)
    }
    return &ProgressDB{
        collection: client.Database(dbName).Collection(collectionName),
    }
}

func (p *ProgressDB) IsProcessed(ctx context.Context, blockchain, address string) (bool, error) {
	filter := bson.M{"blockchain": blockchain, "address": address, "status": "done"}
	count, err := p.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (p *ProgressDB) MarkAsDone(ctx context.Context, blockchain, address string, response map[string]interface{}) error {
	filter := bson.M{"blockchain": blockchain, "address": address}
	update := bson.M{
		"$set": bson.M{
			"status":        "done",
			"last_updated":  time.Now(),
			"response_data": response,
			"retries":       0,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := p.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (p *ProgressDB) MarkAsFailed(ctx context.Context, blockchain, address string, errMsg string) error {
	filter := bson.M{"blockchain": blockchain, "address": address}
	update := bson.M{
		"$set": bson.M{
			"status":       "failed",
			"last_updated": time.Now(),
		},
		"$inc": bson.M{"retries": 1},
	}
	opts := options.Update().SetUpsert(true)
	_, err := p.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (p *ProgressDB) GetPendingTasks(ctx context.Context, limit int64) ([]ProgressRecord, error) {
	filter := bson.M{"status": bson.M{"$ne": "done"}}
	opts := options.Find().SetLimit(limit)
	cursor, err := p.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var records []ProgressRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}
