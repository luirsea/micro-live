package action

import (
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/display"
	"github.com/zyedidia/micro/v2/internal/screen"
	"github.com/zyedidia/micro/v2/internal/transforms"
)

type TransPane struct {
	*BufPane
	*transforms.TransBuf
	TransChain  *transforms.Transform_Chain
	transWindow *display.TransWindow
}

func NewTranPane(tb *transforms.TransBuf, w *display.TransWindow, tab *Tab) *TransPane {
	tp := new(TransPane)
	tp.TransBuf = tb
	tp.transWindow = w
	tp.BufPane = NewBufPane(tb.Buffer, w, tab)
	tp.TransChain = transforms.NewTransformChain(buffer.NewBufferFromString("<no buffer open to transform>", "", buffer.BTDefault))

	return tp
}

func (tp *TransPane) StartTransform(b *buffer.Buffer) {
	tp.TransChain = transforms.NewTransformChain(b)
	tp.HasFocus = true
}

func NewTransBar() *TransPane {
	tb := transforms.NewBuffer()
	w := display.NewTransWindow(tb)
	return NewTranPane(tb, w, nil) // LH TODO the nil tab breaks defualt bufpane binds, (infobar does the same thing so I think it must use different binds)
}

func (tp *TransPane) Display() {
	tp.transWindow.DisplayWithTrivia(tp.TransChain)
}

func (t *TransPane) Close() {
	t.TransBuf.Close()
	t.BufPane.Close()
}

func (t *TransPane) HandleEvent(event tcell.Event) {
	switch e := event.(type) {
	case *tcell.EventResize:
		// TODO
	case *tcell.EventKey:

		t.BufPane.HandleEvent(e)

		resp := string(t.LineBytes(0))

		if success, err := t.TransChain.UpdateTransform(resp); success != true {
			InfoBar.Error(err)
		} else {
			InfoBar.Message("")
			t.UpdateTabs()
		}
	default:
		t.BufPane.HandleEvent(event)
	}
}

func (t *TransPane) UpdateTabs() {

	active := Tabs.Active()
	n := len(Tabs.List)
	// LH TODO this must be inefficient
	Tabs.RemoveAll()

	width, height := screen.Screen.Size()
	iOffset := config.GetGlobalBarsOffset()

	transforms := t.TransChain.Transforms

	for _, tb := range transforms {
		tp := NewTabFromBuffer(0, 0, width, height-iOffset, tb.Buffer())
		tp.SetColor(tb.Trivia.Colour)
		Tabs.AddTab(tp)
	}

	if n >= len(transforms) &&
		active < len(transforms) {
		// LH TODO this is somewhat nieve, could be a totally different transform be we are showing it anyway
		Tabs.SetActive(active)
	} else {
		Tabs.SetActive(len(Tabs.List) - 1)
	}
}
