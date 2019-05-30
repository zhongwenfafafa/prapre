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
		op      clientv3.Op
		opResp  clientv3.OpResponse
		putResp *clientv3.PutResponse
	)

	config = clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2379",
			"127.0.0.1:2380",
		},
		DialTimeout: 5 * time.Second,
	}

	client, err = clientv3.New(config)

	if err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	op = clientv3.OpPut("/cron/jobs/job7", "", clientv3.WithPrevKV())

	opResp, err = kv.Do(context.TODO(), op)
	if err != nil {
		fmt.Println(err)
		return
	}

	putResp = opResp.Put()
	fmt.Println("写入Revision: ", putResp.Header.Revision)
	fmt.Println("创建Revision：", putResp.PrevKv.CreateRevision)
	fmt.Println("更新Revision：", putResp.PrevKv.ModRevision)
}
