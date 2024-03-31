package internal

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgsHelp(t *testing.T) {
	os.Args = []string{"ovo", "-h"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, true)
	assert.Equal(t, args.Verbose, false)
	assert.Equal(t, args.ClearScreen, false)
	assert.Equal(t, args.WatchPath, "")
	assert.Equal(t, args.Cmd, "")
}

func TestParseArgsVerbose(t *testing.T) {
	os.Args = []string{"ovo", "-v"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, true)
	assert.Equal(t, args.Verbose, true)
	assert.Equal(t, args.ClearScreen, false)
	assert.Equal(t, args.WatchPath, "")
	assert.Equal(t, args.Cmd, "")
}

func TestParseArgsWatchPath(t *testing.T) {
	os.Args = []string{"ovo", "myfile.txt"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, true)
	assert.Equal(t, args.Verbose, false)
	assert.Equal(t, args.ClearScreen, false)
	assert.Equal(t, args.WatchPath, "myfile.txt")
	assert.Equal(t, args.Cmd, "")
}

func TestParseArgsWatchPathAndCmd(t *testing.T) {
	os.Args = []string{"ovo", "myfile.txt", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, false)
	assert.Equal(t, args.Verbose, false)
	assert.Equal(t, args.ClearScreen, false)
	assert.Equal(t, args.WatchPath, "myfile.txt")
	assert.Equal(t, args.Cmd, "echo hi")
}

func TestParseArgsDashesAndCmd(t *testing.T) {
	os.Args = []string{"ovo", "--", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, false)
	assert.Equal(t, args.Verbose, false)
	assert.Equal(t, args.ClearScreen, false)
	assert.Equal(t, args.WatchPath, "")
	assert.Equal(t, args.Cmd, "echo hi")
}

func TestSigDefault(t *testing.T) {
	os.Args = []string{"ovo", "--", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Signal, os.Interrupt)
}

func TestSigKill(t *testing.T) {
	os.Args = []string{"ovo", "--sigkill", "--", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Signal, os.Kill)
}

func TestSigInt(t *testing.T) {
	os.Args = []string{"ovo", "--sigint", "--", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Signal, os.Interrupt)
}

func TestSigTerm(t *testing.T) {
	os.Args = []string{"ovo", "--sigterm", "--", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Signal, syscall.SIGTERM)
}

func TestParseArgsAllOptions(t *testing.T) {
	os.Args = []string{"ovo", "-h", "--sigkill", "-v", "-c", "myfile.txt", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, true)
	assert.Equal(t, args.Verbose, true)
	assert.Equal(t, args.ClearScreen, true)
	assert.Equal(t, args.Signal, os.Kill)
	assert.Equal(t, args.WatchPath, "myfile.txt")
	assert.Equal(t, args.Cmd, "echo hi")
}
