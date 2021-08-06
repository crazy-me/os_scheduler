package mongo

import (
	"context"
	"github.com/crazy-me/os_scheduler/work/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	Cli *Mongo
)

type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
}

// InitMongo 初始化Mongo
func InitMongo() (err error) {
	var (
		client *mongo.Client
		db     *mongo.Database
	)
	if client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(conf.C.Mongo.Endpoints).
		SetConnectTimeout(time.Duration(conf.C.Mongo.Timeout)*time.Second)); err != nil {
		return
	}
	db = client.Database(conf.C.Mongo.Db)
	Cli = &Mongo{
		client: client,
		db:     db,
	}
	return
}

// InsertOne 写入单条文档
func (mongodb *Mongo) InsertOne(table string, document interface{}) (insertId string, err error) {
	var (
		insertOneResult *mongo.InsertOneResult
	)
	if insertOneResult, err = mongodb.db.Collection(table).InsertOne(context.TODO(), document); err != nil {
		return
	}
	id := insertOneResult.InsertedID.(primitive.ObjectID)
	insertId = id.Hex()
	return
}
