package action

import (
	"strings"

	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/screen"
)

type Transform_Chain struct {
	transforms []TranBuf
}

type TranBuf struct {
	tran   *Transform
	outBuf *buffer.Buffer
}

func NewTransformChain(baseBuff *buffer.Buffer) *Transform_Chain {
	ts := make([]TranBuf, 1)
	ts[0] = TranBuf{tran: &Transform{}, outBuf: baseBuff}
	return &Transform_Chain{transforms: ts}
}

func (tc *Transform_Chain) Exec(tcS string) (updated bool, err error) {
	// LH TODO This is a bit of a hack, could be better
	// Prepend pipe for the starting nil transform
	tcS = "|" + tcS
	transStrings := strings.Split(tcS, "|")

	// We dont update the actual tc.Transforms until the whole chain succeeded
	newTrans := make([]TranBuf, len(transStrings))
	newTrans[0] = tc.transforms[0]

	var chain_err error
	// The first tran will always be the nil tran, no need to check it
	for i := 1; i < len(transStrings); i++ {

		if i >= len(tc.transforms) ||
			transStrings[i] != tc.transforms[i].tran.raw {
			t := NewTransform(transStrings[i])

			// Here is the chain, the output of the previous transform is the input for the next one
			var outBuf *buffer.Buffer
			outBuf, chain_err = t.Exec(newTrans[i-1].outBuf, i)
			if chain_err != nil {
				break
			}

			newTrans[i] = TranBuf{tran: t, outBuf: outBuf}
		} else {
			newTrans[i] = tc.transforms[i]
		}
	}

	// Only now we know there was no errors do we update the actual transform chain
	if chain_err == nil {
		tc.transforms = newTrans
	}

	updated = chain_err == nil
	return updated, chain_err
}

func (tc *Transform_Chain) UpdateDisplay() {

	active := Tabs.Active()
	n := len(Tabs.List)
	// LH TODO this must be inefficient
	Tabs.RemoveAll()

	width, height := screen.Screen.Size()
	iOffset := config.GetInfoBarOffset()

	for _, tb := range tc.transforms {
		tp := NewTabFromBuffer(0, 0, width, height-iOffset, tb.outBuf)
		Tabs.AddTab(tp)
	}

	if n >= len(tc.transforms) &&
		active < len(tc.transforms) {
		// LH TODO this is somewhat nieve, could be a totally different transform be we are showing it anyway
		Tabs.SetActive(active)
	} else {
		Tabs.SetActive(len(Tabs.List) - 1)
	}
}
