package gomon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Watch(watcherPath string, modExecution string) {
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

	/**
	 * This variable resolve the problem of the show duplicate output to more than one event
	 *
	 * Note: Depending on how the file is modified (using IDE for example), more than one event is triggered for the same modification
	 */
	canShowOutput := true

	go func() {
		for {
			select {
			case <-watcher.Events:

				if canShowOutput {
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

					canShowOutput = false
				} else {
					canShowOutput = true
				}

			case err := <-watcher.Errors:
				fmt.Println("Error in execution the goroutine: ", err)
			}
		}
	}()

	<-done
}
