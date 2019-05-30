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
		putResp *clientv3.PutResponse
	)

	config = clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
		DialTimeout: 5 * time.Second,
	}
	// 创建一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	// 用于读写etcd的健值对
	kv = clientv3.NewKV(client)

	putResp, err = kv.Put(context.TODO(), "/cron/jobs/job2", "hello13", clientv3.WithPrevKV())

	if err != nil {
		fmt.Println(err)
		return
	}

	if putResp.PrevKv != nil {
		fmt.Println("PrevValue: ", string(putResp.PrevKv.Value))
	}
	fmt.Println("Revision: ", putResp.Header.Revision)
}
