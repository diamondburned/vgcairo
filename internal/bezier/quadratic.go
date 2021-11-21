// Package bezier provides helper functions for Bezier curves.
package bezier

import (
	"github.com/diamondburned/gotk4/pkg/cairo"
	"gonum.org/v1/plot/vg"
)

// point exists solely for quadratic-to-cubic conversion.
type point struct{ x, y float64 }

func newPt2(pt float64) point { return point{pt, pt} }

func (pt point) Add(pt1 point) point {
	return point{pt.x + pt1.x, pt.y + pt1.y}
}

func (pt point) Sub(pt1 point) point {
	return point{pt.x - pt1.x, pt.y - pt1.y}
}

func (pt point) Mult(pt1 point) point {
	return point{pt.x * pt1.x, pt.y * pt1.y}
}

// Quadratic draws a quadratic bezier curve into the given Cairo context.
func Quadratic(t *cairo.Context, p1, pt vg.Point) {
	cx, cy := t.GetCurrentPoint()

	// https://stackoverflow.com/a/55034115

	// current point
	qp0 := point{cx, cy}
	// quadratic points
	qp1 := point{float64(p1.X), float64(p1.Y)}
	qp2 := point{float64(pt.X), float64(pt.Y)}

	// cubic points
	cp1 := qp0.Add(newPt2(2.0 / 3.0).Mult(qp1.Sub(qp0)))
	cp2 := qp2.Add(newPt2(2.0 / 3.0).Mult(qp1.Sub(qp2)))

	t.CurveTo(cp1.x, cp1.y, cp2.x, cp2.y, qp2.x, qp2.y)
}

// Cubic draws a cubic bezier curve into the given Cairo context.
func Cubic(t *cairo.Context, p1, p2, pt vg.Point) {
	t.CurveTo(
		float64(p1.X), float64(p1.Y),
		float64(p2.X), float64(p2.Y),
		float64(pt.X), float64(pt.Y),
	)
}
