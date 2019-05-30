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
		getReps *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2380",
		},
		DialTimeout: 5 * time.Second,
	}

	client, err = clientv3.New(config)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	kv = clientv3.NewKV(client)

	getReps, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())

	if err != nil {
		fmt.Println("Error kv : ", err)
		return
	}

	fmt.Printf("%v", getReps.Count)

}
