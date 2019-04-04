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

func main() {
	currentDir, _ := os.Getwd()

	modExecution := "binary"
	watcherPath := currentDir

	for i, args := range os.Args {
		switch args {
		case "--mod":
			modArg := os.Args[i+1]

			switch modArg {
			case "sam":
				modExecution = "sam"
				executionAwsSam()
			}
		case "--path":
			watcherPath = os.Args[i+1]
		case "--version":
			showVersion()
		default:
			buildBinary()
		}
	}

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

func showVersion() {
	fmt.Println("Version: V1.0.0")
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
