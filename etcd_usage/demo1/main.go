package main

import (
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
		DialTimeout: 5 * time.Second,
	}
	// 建立客户端
	client, err = clientv3.New(config)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	client = client

}
