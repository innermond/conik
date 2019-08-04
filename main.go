package main

import (
	"flag"
	"fmt"
	"log"
	"math"
)

var (
	t, b, h, w float64
	p, l, a, g float64
	c, d       bool
)

func param() {
	flag.Float64Var(&t, "T", 0, "Top diameter of the cone")
	flag.Float64Var(&b, "B", 0, "Bottom diameter of the cone")
	flag.Float64Var(&h, "H", 0, "Height of the label")
	flag.Float64Var(&w, "W", 0, "Width of the straightened label")
	flag.Float64Var(&g, "G", 0, "Gap between ends of label")
	flag.BoolVar(&c, "C", false, "Print only numbers")
	flag.BoolVar(&d, "D", false, "Print straighten unfolded cone")
	flag.Parse()
}

func main() {
	param()

	if t == b {
		log.Fatal("not a con, a cilinder")
	}

	flip := false
	if t > b {
		flip = true
		t, b = b, t
	}

	a2 := math.Asin(0.5 * (b - t) / h)
	a = 2 * a2
	if a >= 3 || (b-t)*0.5 >= h {
		log.Fatalf("imposible label having height %.2f for cone with diameters %.2f %.2f/n", h, b, t)
	}
	// top circle unwrapped
	l = math.Pi * t
	// (p+h2)/(b/2) = p/(t/2)
	// segment on top of top-circle to complete the entire cone
	p := 0.5 * t / math.Sin(a2)
	// entire cone radius
	q := p + h

	if w > 0.0 {
		a = w / q
		a2 = 0.5 * a
	}

	if g > 0.0 {
		ag := g / q
		a -= ag
		a2 = 0.5 * a
	}

	// height of encompassing rect for curved label
	hx := q - math.Cos(a2)*p

	// end point
	px := 0.0
	py := p
	// bezier control point
	p1x := math.Tan(a2) * p
	p1y := p
	// start point
	p2x := math.Sin(a) * p
	p2y := math.Cos(a) * p

	// start
	qx := 0.0
	qy := q
	// control
	q1x := math.Tan(a2) * q
	q1y := q
	// end
	q2x := math.Sin(a) * q
	q2y := math.Cos(a) * q

	transform := "translate(%f,%f) "
	rectw := q * a
	transform = fmt.Sprintf(transform, 0.5*rectw, -1.0*q+hx)
	if flip {
		transform = "translate(%f,%f) scale(1,-1) "
		transform = fmt.Sprintf(transform, q2x, q)
	}
	transform += "rotate(%f)"
	degrees := a2 * 180 / math.Pi
	transform = fmt.Sprintf(transform, degrees)
	curve := `
<path
	d="M%f %f Q%f %f %f %f L%f %f Q%f %f %f %f Z"
	transform="%s"
	fill="none" stroke="black" stroke-width="0.1"
/>`
	curve = fmt.Sprintf(curve, p2x, p2y, p1x, p1y, px, py, qx, qy, q1x, q1y, q2x, q2y, transform)

	svg := `
<svg
	width="%f%s" height="%f%s"
	viewBox="%f %f %f %f"
	version="1.1"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink">
	%s%s
</svg>`
	ox, oy := 0.0, 0.0
	o1x, o1y := rectw, hx
	unit := "mm"
	rect := ""
	if d {
		rect = `
<rect
	x="%f" y="%f"
	width="%f" height="%f"
	fill="none" stroke="black" stroke-width="0.1"
/>`
		rect = fmt.Sprintf(rect, 0.0, 0.0, rectw, h)
	}
	svg = fmt.Sprintf(
		svg,
		o1x, unit,
		hx, unit,
		ox, oy, o1x, o1y,
		curve,
		rect,
	)
	if !c {
		fmt.Println(svg)
	} else {
		fmt.Printf("t:%f b:%f h:%f p:%f p1x:%f p1y:%f a:%f q:%f q1x:%f q1y:%f q2x:%f q2y:%f o1x:%f o1y:%f\n", t, b, h, p, p1x, p1y, a*180/math.Pi, q, q1x, q1y, q2x, q2y, o1x, o1y)
	}
}
