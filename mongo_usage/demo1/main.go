package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
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
	database = client.Database("my_db")
	// 3. 选择表my_collection
	collection = database.Collection("my_collection")

	collection = collection
}
