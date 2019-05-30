package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
	)
	cmd = exec.Command("/bin/bash", "-c", "echo 1;")
	err := cmd.Run()

	fmt.Println(err)
}
