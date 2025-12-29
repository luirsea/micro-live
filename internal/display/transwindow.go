package display

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/screen"
	"github.com/zyedidia/micro/v2/internal/transforms"
)

type TransWindow struct {
	*transforms.TransBuf
	*View

	hscroll int
}

func NewTransWindow(b *transforms.TransBuf) *TransWindow {
	tw := new(TransWindow)
	tw.TransBuf = b
	tw.View = new(View)

	tw.Width, tw.Y = screen.Screen.Size()
	tw.Y--

	return tw
}

func (t *TransWindow) Resize(w, h int) {
	t.Width = w
	t.Y = h
}

func (t *TransWindow) SetBuffer(b *buffer.Buffer) {
	t.TransBuf.Buffer = b
}

func (t *TransWindow) Relocate() bool   { return false }
func (t *TransWindow) GetView() *View   { return t.View }
func (t *TransWindow) SetView(v *View)  {}
func (t *TransWindow) SetActive(b bool) {}
func (t *TransWindow) IsActive() bool   { return true }

func (t *TransWindow) LocFromVisual(vloc buffer.Loc) buffer.Loc {
	return buffer.Loc{-1, -1}
}

func (t *TransWindow) BufView() View {
	return View{
		X:         0,
		Y:         t.Y,
		Width:     t.Width,
		Height:    1,
		StartLine: SLoc{0, 0},
		StartCol:  0,
	}
}

func (t *TransWindow) Scroll(s SLoc, n int) SLoc        { return s }
func (t *TransWindow) Diff(s1, s2 SLoc) int             { return 0 }
func (t *TransWindow) SLocFromLoc(loc buffer.Loc) SLoc  { return SLoc{0, 0} }
func (t *TransWindow) VLocFromLoc(loc buffer.Loc) VLoc  { return VLoc{SLoc{0, 0}, loc.X} }
func (t *TransWindow) LocFromVLoc(vloc VLoc) buffer.Loc { return buffer.Loc{vloc.VisualX, 0} }

func (t *TransWindow) Clear() {
	for x := 0; x < t.Width; x++ {
		screen.SetContent(x, t.Y, ' ', nil, config.DefStyle)
	}
}

const transformLabel = "Transform: "

func (t *TransWindow) Display() {
	t.DisplayWithTrivia(nil)
}

func (t *TransWindow) DisplayWithTrivia(tc *transforms.Transform_Chain) {

	x := 0

	for _, c := range transformLabel {
		screen.SetContent(x, t.Y, c, nil, t.defStyle())
		x += runewidth.RuneWidth(c)
	}

	// Skip the nil transform
	curTransform := 1
	inTransform := false

	style := t.defStyle()

	// LH TODO this is getting a little deep, consider breaking up into sup function
	for i, c := range t.Buffer.Line(0) {
		if tc != nil && curTransform < len(tc.Transforms) {
			curTrivia := tc.Transforms[curTransform].Trivia
			if inTransform && i > curTrivia.End {
				style = t.defStyle()
				// We've hit the end of this transform start looking for the start of the next
				curTransform++
				inTransform = false
			} else if !inTransform && i >= curTrivia.Start {
				style = t.defStyle().Background(curTrivia.Colour)
				inTransform = true
			}
		} else {
			style = t.defStyle()
		}

		screen.SetContent(x, t.Y, c, nil, style)
		x += runewidth.RuneWidth(c)
	}

	if t.TransBuf.HasFocus {
		c := t.GetActiveCursor().Loc
		screen.ShowCursor(len(transformLabel)+c.X, t.Y)
	}
}

func (t *TransWindow) defStyle() tcell.Style {
	style := config.DefStyle

	return style
}
