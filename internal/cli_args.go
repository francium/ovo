package internal

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
)

const skipWatchPathMarker = "--"

type CliArgs struct {
	Arg0        string
	Help        bool
	Verbose     bool
	Signal      os.Signal
	StartNow    bool
	ShowTime    bool
	ClearScreen bool
	Cmd         string
	WatchPath   string
}

func ParseArgs() CliArgs {
	args := CliArgs{
		Arg0: path.Base(os.Args[0]),
	}

loop:
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			args.Help = true
		case "-v", "--verbose":
			args.Verbose = true
		case "--now":
			args.StartNow = true
		case "--time":
			args.ShowTime = true
		case "--sigkill":
			args.Signal = os.Kill
		case "--sigint":
			args.Signal = os.Interrupt
		case "--sigterm":
			args.Signal = syscall.SIGTERM
		case "-c", "--clear":
			args.ClearScreen = true
		default:
			if args.WatchPath == "" {
				args.WatchPath = arg
			} else {
				args.Cmd = strings.Join(os.Args[i:], " ")
				break loop
			}
		}
	}

	if args.Cmd == "" {
		args.Help = true
	}

	if args.Signal == nil {
		args.Signal = os.Interrupt
	}

	if args.WatchPath == skipWatchPathMarker {
		args.WatchPath = ""
	}

	return args
}

func DisplayHelp() {
	argv0 := path.Base(os.Args[0])
	fmt.Fprintf(
		os.Stderr,
		`usage: %s [flags] [path] [--] command...

Runs the specified command every time a file system event occurs on or within
the specified path. If the path is not specified, then the default .watchfile
is used. By default a SIGINT is sent to the command to stop it's execution when
a new file system event occurs.

Options must be specified before the command.

Flags:
  -h, --help     Display help
  -v, --verbose  Verbose logging
  -c, --clear    Clear screen every time
  --time         Display timestamp when command is invoked
  --now          Invoke command immediately on start
  --sigint       send SIGINT signal to command (default)
  --sigkill      Send SIGKILL signal to command
  --sigterm      Send SIGTERM signal to command
`,
		argv0,
	)
}
