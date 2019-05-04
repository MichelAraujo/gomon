package gomon

import (
	"fmt"
	"log"
	"os/exec"
)

func ExecutionAwsSam() {
	exec.Command("killall", "sam").Run()
	BuildBinary()

	fmt.Println("## Running AWS Sam CLI ##")
	go execSam()
}

func execSam() {
	sam, execCommandError := exec.Command("sam", "local", "start-api", "-t", "sam.yaml").CombinedOutput()
	if execCommandError != nil {
		log.Println("ERROR: ", execCommandError)
	}
	fmt.Println(string(sam))
}
