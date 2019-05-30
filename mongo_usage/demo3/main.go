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
		interManyResp *mongo.InsertManyResult
		logArr        []interface{}
		insertId      interface{}
		recId         primitive.ObjectID
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

	// 批量插入多条document
	logArr = []interface{}{record, record, record}
	// 发起插入
	if interManyResp, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}

	// 推特开源的：twee的ID生成算法 ID关联时间
	// snowflake 算法规则：毫秒/微妙级别的当前时间 + 机器ID + 当前好秒/微妙内的自增ID(每当好秒产生变化可会重置自增ID为0，继续自增)
	for _, insertId = range interManyResp.InsertedIDs {
		// 拿着interface反射成ObjectID
		recId = insertId.(primitive.ObjectID)
		fmt.Println("自增ID：", recId)
	}
}
