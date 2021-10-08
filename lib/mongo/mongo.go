package mongo

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var syncOnce sync.Once
var client *mongo.Client

func InitMongo(uri string, maxPoolSize uint64) {
	syncOnce.Do(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(maxPoolSize))
			if err != nil {
				panic(err)
			}
			client = _client
		})
}

func GetCollect(dbName, tableName string) *mongo.Collection {
	return client.Database(dbName).Collection(tableName)
}
