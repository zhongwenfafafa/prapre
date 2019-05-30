package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config             clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watcherRespChan    <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
		ctx                context.Context
		cancelFun          context.CancelFunc
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

	kv = clientv3.NewKV(client)

	go func() {
		for {
			_, _ = kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")

			_, _ = kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	// 先Get到当前值，在监听后续变化
	getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7")

	if err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println(string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID，单调递增
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本往后监听: ", watchStartRevision)

	// 5秒后cancel
	ctx, cancelFun = context.WithCancel(context.TODO())

	time.AfterFunc(5*time.Second, func() {
		cancelFun()
	})

	watcherRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))
	// 处理k-v的变化事件
	for watchResp = range watcherRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value),
					"Revision：",
					event.Kv.CreateRevision,
					event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision：", event.Kv.ModRevision)
			}
		}
	}
}
