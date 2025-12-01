package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ride4Low/contracts/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TripsCollection     = "trips"
	RideFaresCollection = "ride_fares"
)

type MongoConfig struct {
	URI      string
	Database string
}

func NewMongoDefaultConfig() *MongoConfig {
	return &MongoConfig{
		URI:      env.GetString("MONGODB_URI", ""),
		Database: env.GetString("MONGODB_DATABASE", ""),
	}
}

func NewMongoClient(cfg *MongoConfig) (*mongo.Client, error) {
	if cfg.URI == "" {
		return nil, fmt.Errorf("mongodb URI is required")
	}
	if cfg.Database == "" {
		return nil, fmt.Errorf("mongodb database is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	err = CreateTTLIndex(ctx, GetDatabase(client, cfg.Database))
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully connected to MongoDB at %s", cfg.URI)
	return client, nil
}

func GetDatabase(client *mongo.Client, database string) *mongo.Database {
	return client.Database(database)
}

func CreateTTLIndex(ctx context.Context, db *mongo.Database) error {
	// Create TTL index that expires documents after 2 hours (7200 seconds)
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"created_at": 1, // index on the created_at field
		},
		Options: options.Index().SetExpireAfterSeconds(7200), // 2 hours TTL
	}

	_, err := db.Collection(RideFaresCollection).Indexes().CreateOne(ctx, indexModel)
	return err
}
