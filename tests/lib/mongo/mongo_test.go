package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/kdcer/go-lib/lib/mongo"

	"github.com/gogf/gf/frame/g"
)

func mongo_redis(t *testing.T) {
	mongo.InitMongo(g.Config().GetString("mongo.uri"), g.Config().GetUint64("mongo.maxPoolSize"))
	mongo.GetCollect("db", "table").FindOne(context.Background(), bson.M{})
}
