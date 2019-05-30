package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"time"
)

func main() {
	var (
		config            clientv3.Config
		client            *clientv3.Client
		err               error
		lease             clientv3.Lease
		leaseGrantResp    *clientv3.LeaseGrantResponse
		leaseId           clientv3.LeaseID
		leaseKeepResp     *clientv3.LeaseKeepAliveResponse
		leaseKeepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		ctx               context.Context
		cancelFun         context.CancelFunc
		kv                clientv3.KV
		txn               clientv3.Txn
		txnResp           *clientv3.TxnResponse
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
	// lease实现锁自动过期：
	// op操作
	// tnx事务：if else then
	lease = clientv3.NewLease(client)

	// 申请一个租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// 获取租约ID
	leaseId = leaseGrantResp.ID

	ctx, cancelFun = context.WithCancel(context.TODO())

	// 释放锁，确保函数退出后租约会自动停止
	defer cancelFun()
	defer lease.Revoke(context.TODO(), leaseId)

	// 租约续租
	if leaseKeepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 续约应答协程处理
	go func() {
		for {
			select {
			case leaseKeepResp = <-leaseKeepRespChan:
				if leaseKeepResp != nil {
					fmt.Println("收到自动续租应答", leaseKeepResp.ID)
				} else {
					fmt.Println("租约已过期")
					goto END
				}
			}
		}
	END:
	}()

	// if不存在key那么设置它，then设置它，else抢锁失败
	// 获取kv处理对象
	kv = clientv3.NewKV(client)
	// 创建事务
	txn = kv.Txn(context.TODO())
	// 定义事务
	txn.If(
		clientv3.Compare(
			clientv3.CreateRevision("/cron/lock/job9"),
			"=",
			0,
		)).
		Then(
			clientv3.OpPut(
				"/cron/lock/job9",
				"XXX",
				clientv3.WithLease(leaseId),
			)).
		Else(clientv3.OpGet("/cron/lock/job9"))
	// 提交事务
	txnResp, err = txn.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 业务处理
	fmt.Println("业务处理")
	time.Sleep(5 * time.Second)
}
