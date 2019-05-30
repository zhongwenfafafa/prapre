package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type forkProcessResult struct {
	err    error
	output []byte
}

func main() {
	// 执行一个cmd，让它在一个协程中执行，让它执行2秒：sleep 2; echo hello;
	// 1秒的时候，我们杀死cmd
	var (
		cxt           context.Context
		cancelFunc    context.CancelFunc
		cmd           *exec.Cmd
		resultChan    chan *forkProcessResult
		processResult *forkProcessResult
	)
	// 创建一个结果队列
	resultChan = make(chan *forkProcessResult, 1000)

	cxt, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err    error
		)

		// 第一个参数是一个Context 上下文
		cmd = exec.CommandContext(cxt, "/bin/bash", "-c", "sleep 2; echo hello")
		// 执行任务捕获输出
		output, err = cmd.CombinedOutput()
		// 把任务输出结果，传给main协程
		resultChan <- &forkProcessResult{
			err:    err,
			output: output,
		}
	}()

	// main 继续执行
	time.Sleep(1 * time.Second)

	// 取消上下文
	cancelFunc()

	// 在main协程里，等待子协程的退出，并打印任务执行结果
	processResult = <-resultChan

	if processResult.err != nil {
		fmt.Println(processResult.err)
		return
	}

	fmt.Println(string(processResult.output))
}
