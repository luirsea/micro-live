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
	Sed
	Unknown
)

// A Transform contains information about how to execute a live bash command
type Transform struct {
	raw   string
	cType CommandType
	args  string
}

func NewTransform(raw string) *Transform {
	t := new(Transform)

	t.raw = raw

	split := strings.SplitN(raw, " ", 2)

	if len(split) < 2 {
		t.cType = Unknown
	} else {
		switch strings.ToLower(split[0]) {
		case "grep":
			t.cType = Grep
			t.args = split[1]
		case "sed":
			t.cType = Sed
			t.args = split[1]
		default:
			t.cType = Unknown
		}
	}

	return t
}

func (t *Transform) Exec(inBuf *buffer.Buffer, n int) (outBuf *buffer.Buffer, err error) {
	var cmd *exec.Cmd
	if cmd, err = t.getCmd(); err != nil {
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

	stderr, err := cmd.StderrPipe()
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

	slurp, _ := io.ReadAll(stderr)

	err = cmd.Wait()

	// An exit code of 1 is not an error for grep, it just means no matches
	if exitErr, ok := err.(*exec.ExitError); ok && t.cType == Grep && exitErr.ExitCode() == 1 {
		err = nil
	}

	if err != nil {
		err = errors.New(err.Error() + ":" + string(slurp) + ":" + outBuf.Line(0))
	}

	return
}

func (t *Transform) getCmd() (*exec.Cmd, error) {
	var cmd exec.Cmd

	// command types are checked so that command specific behavior can be defined (eg grep colouring matches)
	switch t.cType {
	case Grep:
		// LH TODO enable grep colouring matches
		cmd = *exec.Command("grep", t.args)
	case Sed:
		cmd = *exec.Command("sed", t.args)
	default:
		return nil, errors.New("Unknown command")
	}
	return &cmd, nil
}
