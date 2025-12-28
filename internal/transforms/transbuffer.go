package transforms

import (
	"github.com/zyedidia/micro/v2/internal/buffer"
)

// The TransBuf displays messages and other info at the bottom of the screen.
// It is represented as a buffer and a message with a style.
type TransBuf struct {
	*buffer.Buffer

	Arg string

	HasFocus bool

	// LH TODO consider History UX
	// This map stores the history for all the different kinds of uses Prompt has
	// It's a map of history type -> history array
	History    map[string][]string
	HistoryNum int
	// HistorySearch indicates whether we are searching for history items
	// beginning with HistorySearchPrefix
	HistorySearch       bool
	HistorySearchPrefix string
}

// NewBuffer returns a new TransBuffer
func NewBuffer() *TransBuf {
	tb := new(TransBuf)
	tb.History = make(map[string][]string)

	tb.Buffer = buffer.NewBufferFromString("", "", buffer.BTInfo)
	// tb.LoadHistory()

	return tb
}

// Close performs any cleanup necessary when shutting down the TransBuffer
func (tb *TransBuf) Close() {
	// tb.SaveHistory()
}
