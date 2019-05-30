package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"` // 开始时间
	EndTime   int64 `bson:"endTime"`   // 结束时间
}

type LogRecord struct {
	JobName   string `bson:"jobName"` // 任务名
	Command   string `bson:"command"` // shell命令
	Err       string `bson:"err"`     // 脚本错误
	Content   string `bson:"content"` // 脚本输出
	TimePoint `bson:"timePoint"`
}

func main() {
	var (
		client        *mongo.Client
		err           error
		database      *mongo.Database
		collection    *mongo.Collection
		record        *LogRecord
		insertOneResp *mongo.InsertOneResult
		docId         primitive.ObjectID
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
	// 4. 插入一条记录
	record = &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	if insertOneResp, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}

	// _id: 默认生成的一个全局唯一ID，ObjectID：12字节的二进制
	docId = insertOneResp.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID：", docId.Hex()) // 转16进制
}
