package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 任务的执行时间
type TimePoint struct {
	StartTime int64 `bson:"startTime"` // 开始时间
	EndTime   int64 `bson:"endTime"`   // 结束时间
}

// 一条日志
type LogRecord struct {
	JobName   string `bson:"jobName"` // 任务名
	Command   string `bson:"command"` // shell命令
	Err       string `bson:"err"`     // 脚本错误
	Content   string `bson:"content"` // 脚本输出
	TimePoint `bson:"timePoint"`
}

// jobName查询条件
type FindByJobName struct {
	JobName string `bson:"jobName"` //JobName赋值为job10
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		cond       *FindByJobName
		cursor     *mongo.Cursor
		record     *LogRecord
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

	// 读取数据
	// mongo 读取回来的是bson，需要反序列化为LogRecord对象
	// 按照jobName字段过滤，找出jobName为job10 的记录
	cond = &FindByJobName{
		JobName: "job10",
	}
	// 发起查询 + 翻页
	cursor, err = collection.Find(context.TODO(), cond,
		options.Find().SetSkip(0),
		options.Find().SetLimit(2))

	// 释放cursor游标
	defer cursor.Close(context.TODO())

	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		// 反序列化bson对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		// 打印日志
		fmt.Println(*record)
	}
}
