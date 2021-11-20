package vgcairo

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"testing"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestPNG(t *testing.T) {
	const w = 1000
	const h = 1000

	t.Run("simple", func(t *testing.T) {
		surface := cairo.CreateImageSurface(cairo.FORMAT_ARGB32, w, h)
		c := NewCanvas(cairo.Create(surface))

		p := newSimplePlot()
		p.Draw(draw.NewCanvas(c, w, h))

		if testing.Verbose() {
			surface.WriteToPNG(tmpPNG(t, "simple"))
		}
	})

	t.Run("labels", func(t *testing.T) {
		surface := cairo.CreateImageSurface(cairo.FORMAT_ARGB32, w, h)
		c := NewCanvas(cairo.Create(surface))

		p := newLabelsPlot(t)
		p.Draw(draw.NewCanvas(c, w, h))

		if testing.Verbose() {
			surface.WriteToPNG(tmpPNG(t, "labels"))
		}
	})
}

func tmpPNG(t *testing.T, suffix string) string {
	f, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("vgcairo-%s-*.png", suffix))
	if err != nil {
		t.Fatal("CreateTemp error:", err)
	}

	t.Log("new PNG at", f.Name())

	t.Cleanup(func() {
		f.Close()

		if !testing.Verbose() {
			os.Remove(f.Name())
		}
	})

	return f.Name()
}

// https://github.com/gonum/plot/blob/v0.10.0/vg/vggio/vggio_example_test.go
func newSimplePlot() *plot.Plot {
	p := plot.New()
	p.Title.Text = "My title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	quad := plotter.NewFunction(func(x float64) float64 {
		return x * x
	})
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := plotter.NewFunction(func(x float64) float64 {
		return math.Pow(2, x)
	})
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := plotter.NewFunction(func(x float64) float64 {
		return 10*math.Sin(x) + 50
	})
	sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	sin.Width = vg.Points(4)
	sin.Color = color.RGBA{R: 255, A: 255}

	p.Add(quad, exp, sin)
	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("10*sin(x)+50", sin)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 100

	p.Add(plotter.NewGrid())
	return p
}

func newLabelsPlot(t *testing.T) *plot.Plot {
	p := plot.New()
	p.Title.Text = "Labels"
	p.X.Min = -1
	p.X.Max = +1
	p.Y.Min = -1
	p.Y.Max = +1

	const (
		left   = 0.00
		middle = 0.02
		right  = 0.04
	)

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: -0.8 + left, Y: -0.5},   // Aq + y-align bottom
			{X: -0.6 + middle, Y: -0.5}, // Aq + y-align center
			{X: -0.4 + right, Y: -0.5},  // Aq + y-align top

			{X: -0.8 + left, Y: +0.5}, // ditto for Aq\nAq
			{X: -0.6 + middle, Y: +0.5},
			{X: -0.4 + right, Y: +0.5},

			{X: +0.0 + left, Y: +0}, // ditto for Bg\nBg\nBg
			{X: +0.2 + middle, Y: +0},
			{X: +0.4 + right, Y: +0},
		},
		Labels: []string{
			"Aq", "Aq", "Aq",
			"Aq\nAq", "Aq\nAq", "Aq\nAq",

			"Bg\nBg\nBg",
			"Bg\nBg\nBg",
			"Bg\nBg\nBg",
		},
	})
	if err != nil {
		t.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range labels.TextStyle {
		sty := &labels.TextStyle[i]
		sty.Font.Size = vg.Length(34)
	}
	labels.TextStyle[0].YAlign = draw.YBottom
	labels.TextStyle[1].YAlign = draw.YCenter
	labels.TextStyle[2].YAlign = draw.YTop

	labels.TextStyle[3].YAlign = draw.YBottom
	labels.TextStyle[4].YAlign = draw.YCenter
	labels.TextStyle[5].YAlign = draw.YTop

	labels.TextStyle[6].YAlign = draw.YBottom
	labels.TextStyle[7].YAlign = draw.YCenter
	labels.TextStyle[8].YAlign = draw.YTop

	lred, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: -0.8 + left, Y: +0.5},
			{X: +0.0 + left, Y: +0},
		},
		Labels: []string{
			"Aq", "Bg",
		},
	})
	if err != nil {
		t.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range lred.TextStyle {
		sty := &lred.TextStyle[i]
		sty.Font.Size = vg.Length(34)
		sty.Color = color.RGBA{R: 255, A: 255}
		sty.YAlign = draw.YBottom
	}

	m5 := plotter.NewFunction(func(float64) float64 { return -0.5 })
	m5.LineStyle.Color = color.RGBA{R: 255, A: 255}

	l0 := plotter.NewFunction(func(float64) float64 { return 0 })
	l0.LineStyle.Color = color.RGBA{G: 255, A: 255}

	p5 := plotter.NewFunction(func(float64) float64 { return +0.5 })
	p5.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(labels, lred, m5, l0, p5)
	p.Add(plotter.NewGrid())
	p.Add(plotter.NewGlyphBoxes())
	return p
}
