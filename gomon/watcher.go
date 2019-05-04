package gomon

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func Watch(watcherPath string, modExecution string) {
	var watcher *fsnotify.Watcher
	var newWatcherError error

	watcher, newWatcherError = fsnotify.NewWatcher()
	if newWatcherError != nil {
		log.Println("ERROR: ", newWatcherError)
	}

	defer watcher.Close()

	watcher.Add(watcherPath)

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

						buildInfo, execCommandError := exec.Command("./main").CombinedOutput()
						if execCommandError != nil {
							log.Println("ERROR: ", execCommandError)
						}
						fmt.Println(string(buildInfo))
					}

					canShowOutput = false
				} else {
					canShowOutput = true
				}

			case watcherGoroutineError := <-watcher.Errors:
				fmt.Println("Error in execution the goroutine: ", watcherGoroutineError)
			}
		}
	}()

	<-done
}
