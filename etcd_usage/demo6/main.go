package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"time"
)

func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		lease        clientv3.Lease
		leaseResp    *clientv3.LeaseGrantResponse
		leaseId      clientv3.LeaseID
		kv           clientv3.KV
		putResp      *clientv3.PutResponse
		getResp      *clientv3.GetResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
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
	// 申请一个lease（租约）
	lease = clientv3.NewLease(client)
	// 申请一个10秒的租约
	leaseResp, err = lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Put一个k-v，让它与租约关联起来，从而实现10秒后过期
	leaseId = leaseResp.ID

	// 5秒后自动取消上下文
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)

	// 自动续租, 5秒后会自动取消续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else { // 每秒会续租一次
					fmt.Println("收到自动续租应答", keepResp.ID)
				}
			}
		}
	END:
	}()

	kv = clientv3.NewKV(client)

	putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功：", putResp.Header.Revision)

	// 写一个for循环查询key是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("key-value 过期了")
			break
		}

		fmt.Println("还没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
