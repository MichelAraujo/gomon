package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gomon/gomon"
)

const modExecutionDefault = "binary"

/**
 * GOMON - V1.0.0
 *
 * GNU GENERAL PUBLIC LICENSE
 * https://github.com/MichelAraujo/gomon/blob/master/LICENSE
 *
 * Compiles and executes golang codes as files are changed
 */
func main() {
	modExecution, watcherPath := setParameters(os.Args)
	gomon.Watch(watcherPath, modExecution)
}

func setParameters(inputArgs []string) (string, string) {
	currentDir, getWdError := os.Getwd()
	if getWdError != nil {
		log.Println("ERROR: ", getWdError.Error())
	}

	modExecution := modExecutionDefault
	watcherPath := currentDir

	for i, args := range inputArgs {
		switch args {
		case "--mod":
			modArg := inputArgs[i+1]

			switch modArg {
			case "sam":
				modExecution = "sam"
				gomon.ExecutionAwsSam()
			case modExecutionDefault:
				gomon.BuildBinary()
			}
		case "--path":
			watcherPath = inputArgs[i+1]
		case "--version":
			showVersion()
		}
	}

	return modExecution, watcherPath
}

func showVersion() {
	fmt.Println("Version: V1.0.0")
	os.Exit(3)
}
