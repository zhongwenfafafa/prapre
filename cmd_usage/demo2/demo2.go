package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)
	// 生产cmd
	cmd = exec.Command("/bin/bash", "-c", "sleep 5;ls -l")

	// 执行命令，捕获子进程的输出（pipe）
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// 打印子进程输出
	fmt.Printf("%s", output)
}
