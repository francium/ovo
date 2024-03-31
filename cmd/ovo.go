package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/francium/ovo/internal"
	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/fatih/color"
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

	i := 0

	go handleFSEvent(watcher, func(path string) {
		i++
		j := i

		log.Info("Execution routine start ", j)

		select {
		case cancel <- struct{}{}:
		default:
		}

		cmd := exec.Command(
			"bash", "-c", args.Cmd,
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		done := make(chan struct{})

		killing := false

		go func() {
			log.Info("Cancel routine start ", j)

			<-cancel

			log.Info("Cancel routine received cancel sign")

			if cmd != nil {
				log.Info("Cancel routine killing cmd")
				killing = true
				cmd.Process.Signal(args.Signal)
			} else {
				log.Info("cmd is nil")
			}

			log.Info("Cancel routine end ", j)
		}()

		go func() {
			log.Info("Runner routine start ", j)

			if args.ClearScreen {
				fmt.Print(clearScreenEscapeSequence)
			}

			bold := color.New(color.Bold).SprintFunc()
			fmt.Println(bold("> ", args.Cmd))

			err := cmd.Run()
			if err != nil && !killing {
				log.Error("Failed to run command: ", err)
			}

			select {
			case done <- struct{}{}:
			default:
			}

			log.Info("Runner routine end ", j)
		}()

		select {
		case <-done:
		case <-cancel:
		}

		select {
		case cancel <- struct{}{}:
		default:
		}

		if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			log.Fatal("Process did not exit")
		}

		log.Info("Execution routine end ", j)
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
