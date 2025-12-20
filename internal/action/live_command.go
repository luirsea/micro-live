package action

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/zyedidia/micro/v2/internal/buffer"
)

type CommandType int

const (
	NilCommand = iota
	Grep
	Unknown
)

// A LiveCommand contains information about how to execute a live bash command
type LiveCommand struct {
	raw    string
	cType  CommandType
	args   string
	outBuf *buffer.Buffer
}

func NewLiveCommand(buf *buffer.Buffer, raw string) *LiveCommand {
	l := new(LiveCommand)

	l.raw = raw

	split := strings.SplitN(raw, " ", 2)

	if len(split) < 2 {
		l.cType = Unknown
	} else {
		switch split[0] {
		case "grep":
			l.cType = Grep
			l.args = split[1]
		default:
			l.cType = Unknown
		}
	}
	l.outBuf = buf
	return l
}

func (lc *LiveCommand) Exec(inBuf *buffer.Buffer, n int) (outBuf *buffer.Buffer, err error) {

	var cmd *exec.Cmd
	if cmd, err = lc.getCmd(); err != nil {
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	go func() {
		defer stdin.Close()
		// LH TODO this has the be the most inefficient way to stream the buffer to stdin
		io.Writer.Write(stdin, inBuf.LineArray.Bytes())
	}()

	outBuf = buffer.NewBuffer(stdout, 0, "", buffer.BTDefault, *new(buffer.Command))
	outBuf.SetName(fmt.Sprintf("[%d]%s", n, inBuf.GetName()))

	err = cmd.Wait()

	return
}

func (lc *LiveCommand) getCmd() (*exec.Cmd, error) {
	var cmd exec.Cmd

	// command types are checked so that command specific behavior can be defined (eg grep colouring matches)
	switch lc.cType {
	case Grep:
		// LH TODO enable grep colouring matches
		cmd = *exec.Command("grep", lc.args)
	default:
		return nil, errors.New("Unknown command")
	}
	return &cmd, nil
}
