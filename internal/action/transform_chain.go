package action

import (
	"strings"
	"unicode"

	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/buffer"
)

type Transform_Chain struct {
	transforms []TranBuf
	path       string
}

type TranBuf struct {
	tran   *Transform
	outBuf *buffer.Buffer
	trivia *Trivia
}

type Trivia struct {
	start  int
	end    int
	colour tcell.Color
}

func NewTransformChain(baseBuff *buffer.Buffer) *Transform_Chain {
	ts := make([]TranBuf, 1)
	ts[0] = TranBuf{tran: &Transform{}, outBuf: baseBuff}
	return &Transform_Chain{transforms: ts, path: baseBuff.Path}
}

func (tc *Transform_Chain) UpdateTransform(tcS string) (updated bool, err error) {
	// LH TODO This is a bit of a hack, could be better
	// Prepend pipe for the starting nil transform
	if len(tcS) > 0 {
		tcS = "|" + tcS

	}

	transStrings := strings.Split(tcS, "|")

	// We dont update the actual tc.Transforms until the whole chain succeeded
	newTrans := make([]TranBuf, len(transStrings))
	newTrans[0] = tc.transforms[0]

	var chain_err error

	x := 0
	// The first tran will always be the nil tran, no need to check it
	for i := 1; i < len(transStrings); i++ {

		if i >= len(tc.transforms) ||
			!tc.transforms[i].tran.Matches(transStrings[i]) {

			t := NewTransform(transStrings[i])

			// Here is the chain, the output of the previous transform is the input for the next one
			var outBuf *buffer.Buffer
			outBuf, chain_err = t.Exec(newTrans[i-1].outBuf, tc.path)
			if chain_err != nil {
				break
			}

			// calculate start and end of trimmed transform, used for display nice to haves
			start := len(transStrings[i]) -
				len(strings.TrimLeftFunc(transStrings[i], unicode.IsSpace))
			length := len(strings.TrimSpace(transStrings[i]))

			trivia := &Trivia{start: x + start, end: x + start + length, colour: newTrans[i-1].trivia.colour + 1}

			newTrans[i] = TranBuf{tran: t, outBuf: outBuf, trivia: trivia}
		} else {
			newTrans[i] = tc.transforms[i]
		}

		x += len(transStrings[i])
	}

	// Only now we know there was no errors do we update the actual transform chain
	if chain_err == nil {
		tc.transforms = newTrans
	}

	updated = chain_err == nil
	return updated, chain_err
}
