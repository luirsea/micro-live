package transforms

import (
	"strings"
	"unicode"

	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/buffer"
)

type Transform_Chain struct {
	Transforms []TranBuf
	Path       string
}

type TranBuf struct {
	tran   *Transform
	outBuf *buffer.Buffer
	Trivia *Trivia
}

type Trivia struct {
	Start  int
	End    int
	Colour tcell.Color
}

func (tb *TranBuf) Buffer() *buffer.Buffer {
	return tb.outBuf
}

func NewTransformChain(baseBuff *buffer.Buffer) *Transform_Chain {
	ts := make([]TranBuf, 1)
	ts[0] = TranBuf{tran: &Transform{}, outBuf: baseBuff, Trivia: new(Trivia)}
	return &Transform_Chain{Transforms: ts, Path: baseBuff.Path}
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
	newTrans[0] = tc.Transforms[0]

	var chain_err error

	x := 0
	// The first tran will always be the nil tran, no need to check it
	for i := 1; i < len(transStrings); i++ {

		if i >= len(tc.Transforms) ||
			!tc.Transforms[i].tran.Matches(transStrings[i]) {

			t := NewTransform(transStrings[i])

			// Here is the chain, the output of the previous transform is the input for the next one
			var outBuf *buffer.Buffer
			outBuf, chain_err = t.Exec(newTrans[i-1].outBuf, tc.Path, t.CommandLabel())
			if chain_err != nil {
				break
			}

			// calculate start and end of trimmed transform, used for display nice to haves
			start := len(transStrings[i]) -
				len(strings.TrimLeftFunc(transStrings[i], unicode.IsSpace))
			length := len(strings.TrimSpace(transStrings[i]))

			trivia := &Trivia{Start: x + start, End: x + start + length - 1, Colour: colourPalette[i%len(colourPalette)]}

			newTrans[i] = TranBuf{tran: t, outBuf: outBuf, Trivia: trivia}
		} else {
			newTrans[i] = tc.Transforms[i]
		}

		x += len(transStrings[i]) + 1 // add one for joining pipe char
	}

	// Only now we know there was no errors do we update the actual transform chain
	if chain_err == nil {
		tc.Transforms = newTrans
	}

	updated = chain_err == nil
	return updated, chain_err
}

var colourPalette = [5]tcell.Color{
	tcell.ColorDarkCyan,
	tcell.ColorDarkGreen,
	tcell.ColorDarkOrange,
	tcell.ColorDarkKhaki,
	tcell.ColorDarkMagenta}
