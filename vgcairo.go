// Package vgcairo provides a Cairo drawing backend for gonum.org/plot/vg.
package vgcairo

import (
	"image"
	"image/color"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"

	xfont "golang.org/x/image/font"
)

// Canvas implements the vg.Canvas interface, drawing to a cairo.Context.
type Canvas struct {
	t *cairo.Context
}

var _ vg.Canvas = (*Canvas)(nil)

// NewCanvas returns a new Cairo Canvas.
func NewCanvas(t *cairo.Context) Canvas {
	return Canvas{t: t}
}

// Context returns the Canvas' cairo.Context.
func (c Canvas) Context() *cairo.Context {
	return c.t
}

// SetLineWidth implements vg.Canvas.
func (c Canvas) SetLineWidth(l vg.Length) {
	c.t.SetLineWidth(float64(l))
}

// SetLineDash implements vg.Canvas.
func (c Canvas) SetLineDash(pattern []vg.Length, offset vg.Length) {
	// TODO: probably optimize this.
	dashes := make([]float64, len(pattern))
	for i, dash := range pattern {
		dashes[i] = float64(dash)
	}
	c.t.SetDash(dashes, float64(offset))
}

// SetColor implements vg.Canvas.
func (c Canvas) SetColor(clr color.Color) {
	if clr == nil {
		clr = color.Black
	}

	r, g, b, a := clr.RGBA()
	c.t.SetSourceRGBA(
		float64(r)/0xFFFF,
		float64(g)/0xFFFF,
		float64(b)/0xFFFF,
		float64(a)/0xFFFF,
	)
}

// Rotate implements vg.Canvas.
func (c Canvas) Rotate(rad float64) {
	c.t.Rotate(rad)
}

// Translate implements vg.Canvas.
func (c Canvas) Translate(pt vg.Point) {
	c.t.Translate(float64(pt.X), float64(pt.Y))
}

// Scale implements vg.Canvas.
func (c Canvas) Scale(x, y float64) {
	c.t.Scale(x, y)
}

// Push implements vg.Canvas.
func (c Canvas) Push() {
	c.t.Save()
}

// Pop implements vg.Canvas.
func (c Canvas) Pop() {
	c.t.Restore()
}

// Fill implements vg.Canvas.
func (c Canvas) Fill(path vg.Path) {
	c.t.Save()
	c.outline(path)
	c.t.Fill()
	c.t.Restore()
}

// Stroke implements vg.Canvas.
func (c Canvas) Stroke(path vg.Path) {
	c.t.Save()
	c.outline(path)
	c.t.Stroke()
	c.t.Restore()
}

func (c Canvas) outline(path vg.Path) {
	c.t.NewPath()

	for _, comp := range path {
		switch comp.Type {
		case vg.MoveComp:
			c.t.MoveTo(float64(comp.Pos.X), float64(comp.Pos.Y))
		case vg.LineComp:
			c.t.LineTo(float64(comp.Pos.X), float64(comp.Pos.Y))
		case vg.ArcComp:
			c.t.Arc(float64(comp.Pos.X), float64(comp.Pos.Y), float64(comp.Radius.Points()), comp.Start, comp.Angle)
		case vg.CurveComp:
			switch len(comp.Control) {
			case 1: // quadratic
				c.t.CurveTo(
					float64(comp.Control[0].X), float64(comp.Control[0].Y),
					// https://bugs.launchpad.net/inkscape/+bug/1009765/comments/1
					// http://caffeineowl.com/graphics/2d/vectorial/cubic2quad01.html
					float64(2*comp.Control[0].X+comp.Pos.X)/3, float64(2*comp.Control[0].Y+comp.Pos.Y)/3,
					float64(comp.Pos.X), float64(comp.Pos.Y),
				)
			case 2: // cubic
				c.t.CurveTo(
					float64(comp.Control[0].X), float64(comp.Control[0].Y),
					float64(comp.Control[1].X), float64(comp.Control[1].Y),
					float64(comp.Pos.X), float64(comp.Pos.Y),
				)
			default:
				panic("vgcairo: invalid number of control points")
			}
		case vg.CloseComp:
			c.t.ClosePath()
		}
	}
}

// FillString implements vg.Canvas.
func (c Canvas) FillString(f font.Face, pt vg.Point, text string) {
	weight := cairo.FONT_WEIGHT_NORMAL
	switch f.Font.Weight {
	case xfont.WeightSemiBold, xfont.WeightBold, xfont.WeightExtraBold, xfont.WeightBlack:
		weight = cairo.FONT_WEIGHT_BOLD
	}

	c.t.Save()
	c.t.SelectFontFace(string(f.Font.Typeface), cairo.FONT_SLANT_NORMAL, weight)
	c.t.MoveTo(float64(pt.X), float64(pt.Y))
	c.t.ShowText(text)
	c.t.Restore()
}

// DrawImage implements vg.Canvas.
func (c Canvas) DrawImage(rect vg.Rectangle, img image.Image) {
	rectSz := rect.Size()
	imagSz := img.Bounds().Size()

	surface := cairo.CreateSurfaceFromImage(img)

	c.t.Save()
	defer c.t.Restore()

	if !vgPtEq(rectSz, imagSz) {
		c.t.Scale(
			float64(rectSz.X)/float64(imagSz.X),
			float64(rectSz.Y)/float64(imagSz.Y),
		)
	}

	c.t.SetSourceSurface(surface, rect.Min.X.Points(), rect.Min.Y.Points())
	c.t.Paint()
}

func vgPtEq(pt1 vg.Point, pt2 image.Point) bool {
	return int(pt1.X) == pt2.X && int(pt1.Y) == pt2.Y
}
