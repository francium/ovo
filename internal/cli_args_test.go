package internal

import (
	"os"
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

func TestParseArgsAllOptions(t *testing.T) {
	os.Args = []string{"ovo", "-h", "-v", "-c", "myfile.txt", "echo hi"}
	args := ParseArgs()
	assert.Equal(t, args.Arg0, "ovo")
	assert.Equal(t, args.Help, true)
	assert.Equal(t, args.Verbose, true)
	assert.Equal(t, args.ClearScreen, true)
	assert.Equal(t, args.WatchPath, "myfile.txt")
	assert.Equal(t, args.Cmd, "echo hi")
}
