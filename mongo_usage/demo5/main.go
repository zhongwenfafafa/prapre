package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type StartTimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

type DeleteCond struct {
	BeforeCond StartTimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
		delResult  *mongo.DeleteResult
	)

	// 1. 建立连接
	// 创建客户端
	if client, err = mongo.NewClient(
		options.Client().ApplyURI("mongodb://127.0.0.1:27017"),
		options.Client().SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}
	// 连接
	if err = client.Connect(context.TODO()); err != nil {
		fmt.Println(err)
		return
	}

	// 2. 选择数据库my_db
	database = client.Database("cron")
	// 3. 选择表my_collection
	collection = database.Collection("log")

	// 删除数据
	// 定义删除条件结构体
	delCond = &DeleteCond{
		BeforeCond: StartTimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}
	// 执行删除
	delResult, err = collection.DeleteMany(context.TODO(),
		delCond)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数：", delResult.DeletedCount)
}
