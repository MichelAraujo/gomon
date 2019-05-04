package gomon

import (
	"fmt"
	"log"
	"os/exec"
)

func BuildBinary() {
	buildInfo, execCommandError := exec.Command("make", "build").CombinedOutput()
	if execCommandError != nil {
		log.Println("ERROR: ", execCommandError)
	}
	fmt.Println("## Build info: ", string(buildInfo))
}
