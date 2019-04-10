package gomon

import (
	"fmt"
	"log"
	"os/exec"
)

func BuildBinary() {
	buildInfo, err := exec.Command("make", "build").CombinedOutput()
	if err != nil {
		log.Println("Could not execute command: ", err)
	}
	fmt.Println("## Build info: ", string(buildInfo))
}
