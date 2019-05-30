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
		delResp *clientv3.DeleteResponse
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
		fmt.Println("Error client connecting: ", err)
		return
	}

	kv = clientv3.NewKV(client)

	delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrevKV(), clientv3.WithPrefix())

	if err != nil {
		fmt.Println("Error deleting Kv: ", err)
		return
	}

	if len(delResp.PrevKvs) != 0 {
		for _, v := range delResp.PrevKvs {
			fmt.Println("删除了：", string(v.Key), string(v.Value))
		}
	}

	fmt.Println(delResp.Deleted)
}
