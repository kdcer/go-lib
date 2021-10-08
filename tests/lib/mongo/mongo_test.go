package mongoext

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gogf/gf/frame/g"
	"github.com/kdcer/go-lib/lib/mongoext"
)

func mongo_redis(t *testing.T) {
	mongoext.InitMongo(g.Config().GetString("mongo.uri"), g.Config().GetString("mongo.db"), g.Config().GetUint64("mongo.maxPoolSize"))
	mongoext.GetCollect("table").FindOne(context.Background(), bson.M{})
}
