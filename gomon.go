package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

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

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(watcherPath, watchDir); err != nil {
		fmt.Println("Error in add watcher path: ", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-watcher.Events:
				fmt.Println("### Changes detected, rebuild ... ###")

				switch modExecution {
				case "sam":
					executionAwsSam()
				case "binary":
					buildBinary()

					buildInfo, err := exec.Command("./main").CombinedOutput()
					if err != nil {
						log.Println("Could not execute command: ", err)
					}
					fmt.Println(string(buildInfo))
				}

			case err := <-watcher.Errors:
				fmt.Println("Error in execution the goroutine: ", err)
			}
		}
	}()

	<-done
}

func setParameters(inputArgs []string) (string, string) {
	currentDir, _ := os.Getwd()

	modExecution := modExecutionDefault
	watcherPath := currentDir

	for i, args := range inputArgs {
		switch args {
		case "--mod":
			modArg := inputArgs[i+1]

			switch modArg {
			case "sam":
				modExecution = "sam"
				executionAwsSam()
			case modExecutionDefault:
				buildBinary()
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

func executionAwsSam() {
	exec.Command("killall", "sam").Run()
	buildBinary()

	fmt.Println("## Running AWS Sam CLI ##")
	go execSam()
}

func buildBinary() {
	buildInfo, err := exec.Command("make", "build").CombinedOutput()
	if err != nil {
		log.Println("Could not execute command: ", err)
	}
	fmt.Println("## Build info: ", string(buildInfo))
}

func execSam() {
	sam, _ := exec.Command("sam", "local", "start-api", "-t", "sam.yaml").CombinedOutput()
	fmt.Println(string(sam))
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
