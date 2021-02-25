package db

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	baseMongoURI = "mongodb://mongo:27017"
	document     = "market_data"
)

func getMongoURI() string {
	return baseMongoURI
}

// Config : db config
type Config struct {
	*mongo.Database
}

var dbConfig *Config

var syncOnce sync.Once

// GetDB : get db
func GetDB() *Config {
	syncOnce.Do(func() {
		dbConfig = loadDBConfig()
	})
	return dbConfig
}

// loadDBConfig : load db config
func loadDBConfig() *Config {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(getMongoURI()))
	if err != nil {
		panic(err.Error())
	}

	coins := client.Database(document)
	return &Config{coins}
}
