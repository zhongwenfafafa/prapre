package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
		DialTimeout: 5 * time.Second,
	}
	// 使用配置文件创建客户端
	client, err = clientv3.New(config)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// 创建kv操作对象用于读写 k-v
	kv = clientv3.NewKV(client)

	getResp, err = kv.Get(context.TODO(), "/cron/jobs/job2")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, v := range getResp.Kvs {
		fmt.Printf("(key, value, version) = (%s, %s, %d)\n", v.Key, v.Value, v.Version)
	}

}
