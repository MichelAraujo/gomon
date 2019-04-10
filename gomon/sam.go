package gomon

import (
	"fmt"
	"os/exec"
)

func ExecutionAwsSam() {
	exec.Command("killall", "sam").Run()
	BuildBinary()

	fmt.Println("## Running AWS Sam CLI ##")
	go execSam()
}

func execSam() {
	sam, _ := exec.Command("sam", "local", "start-api", "-t", "sam.yaml").CombinedOutput()
	fmt.Println(string(sam))
}
