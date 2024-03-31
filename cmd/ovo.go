package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/francium/ovo/internal"
	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

const defaultWatchFile = ".watchfile"
const clearScreenEscapeSequence = "\x1bc"

func main() {
	args := internal.ParseArgs()
	if args.WatchPath == "" {
		args.WatchPath = defaultWatchFile
	}

	if args.Help {
		internal.DisplayHelp()
		os.Exit(0)
	}

	internal.InitLogging(args.Verbose)

	if args.WatchPath == defaultWatchFile {
		internal.CreateIfNotExists(defaultWatchFile)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Info("Watching ", args.WatchPath)
	err = watcher.Add(args.WatchPath)
	if err != nil {
		log.Fatal(err)
	}

	cancel := make(chan struct{})
	runner_lock := sync.Mutex{}

	i := 0

	go handleFSEvent(watcher, func(path string) {
		i++
		log_prefix := fmt.Sprintf("(%d) ", i)

		log.Info(log_prefix, "Execution routine start")

		select {
		case cancel <- struct{}{}:
		default:
		}

		log.Info(log_prefix, "Execution routine waiting for runner lock ")
		runner_lock.Lock()
		log.Info(log_prefix, "Execution routine acquired runner lock ")

		cmd := exec.Command(
			"bash", "-c", args.Cmd,
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		done := make(chan struct{})

		killing := false

		go func() {
			log.Info(log_prefix, "Cancel routine start ")

			<-cancel

			log.Info(log_prefix, "Cancel routine received cancel sign")

			if cmd != nil {
				log.Info(log_prefix, "Cancel routine killing cmd")
				killing = true
				cmd.Process.Signal(args.Signal)
			} else {
				log.Warn(log_prefix, "Cmd is nil")
			}

			log.Info(log_prefix, "Cancel routine end ")
		}()

		go func() {
			log.Info(log_prefix, "Runner routine start ")

			if args.ClearScreen {
				fmt.Print(clearScreenEscapeSequence)
			}

			bold := color.New(color.Bold).SprintFunc()
			fmt.Println(bold("> ", args.Cmd))

			log.Info(log_prefix, "Runner routine invoking command ")
			err := cmd.Run()
			if err != nil && !killing {
				log.Error(log_prefix, "Failed to run command: ", err)
			}

			select {
			case done <- struct{}{}:
			default:
			}

			log.Info(log_prefix, "Runner routine end ")
		}()

		// It's possible this logic below isn't 100% bulletproof. Especially
		// the way the cancel is triggered in the go routine below and the
		// done/cancel block is below that.

		// Here we want to first block until one of done or cancel send a
		// message
		select {
		case <-done:
		case <-cancel:
		}

		// Setup an async task to send a message to trigger the cancellation
		go func() {
			select {
			case cancel <- struct{}{}:
			default:
			}
		}()

		// Wait until until the other channel sends a message.
		select {
		case <-done:
		case <-cancel:
		}

		if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			log.Fatal(log_prefix, "Process did not exit")
		}

		log.Info(log_prefix, "Execution routine releasing lock ")
		runner_lock.Unlock()

		log.Info(log_prefix, "Execution routine end ")
	})

	// Block
	<-make(chan struct{})
}

func handleFSEvent(watcher *fsnotify.Watcher, cb func(path string)) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Error("Failed to read event from channel")
				return
			}
			log.Infof("path=%s, op=%s", event.Name, event.Op)
			go cb(event.Name)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Error:", err)
		}
	}
}
