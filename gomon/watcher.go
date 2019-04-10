package gomon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watch(watcherPath string, modExecution string) {
	var watcher *fsnotify.Watcher

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(
		watcherPath,
		func(path string, fi os.FileInfo, err error) error {
			if fi.Mode().IsDir() {
				return watcher.Add(path)
			}

			return nil
		},
	)
	if err != nil {
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
					ExecutionAwsSam()
				case "binary":
					BuildBinary()

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
