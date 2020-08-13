package api

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	config *Config
	conn   *mongo.Database
}

func NewDB(config *Config) *DB {
	return &DB{
		config: config,
	}
}

func (db *DB) Open() error {
	ctx := context.Background()
	dburl := os.Getenv("dburl")
	clientOpts := options.Client().ApplyURI(dburl)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return (err)
	}
	db.conn = client.Database(db.config.DatabaseName)

	return nil
}

func (db *DB) Close() {

}
