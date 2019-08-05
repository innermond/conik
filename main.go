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
	flag.Float64Var(&t, "T", 0.0, "Top diameter of the cone")
	flag.Float64Var(&b, "B", 0.0, "Bottom diameter of the cone")
	flag.Float64Var(&h, "H", 0.0, "Height of the label")
	flag.Float64Var(&w, "W", 0.0, "Width of the straightened label")
	flag.Float64Var(&g, "G", 0.0, "Gap between ends of label")
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

	// half of cone angle (in radians)
	a2 := math.Asin(0.5 * (b - t) / h)
	a = 2 * a2
	// segment on top of top-circle to complete the entire cone
	p := 0.5 * t / math.Sin(a2)
	// entire cone radius
	q := p + h

	// adjust angle regarding requested label width
	if w > 0.0 {
		a = w / q
		a2 = 0.5 * a
	}
	// adjust angle to take in account a gap betweem ends
	if g != 0.0 {
		ag := g / q
		a -= ag
		a2 = 0.5 * a
	}

	// unfolded cone no more than sem-circle - 3 radians
	if a >= 3 || (b-t)*0.5 >= h {
		log.Fatalf("imposible label having height %.2f for cone with diameters %.2f %.2f\n", h, b, t)
	}

	// height of encompassing rect for curved label
	hx := q - math.Cos(a2)*p
	// width of encompassing rect for curved label
	wx := q * a

	// end point
	px := 0.0
	py := p
	// start point
	p2x := math.Sin(a) * p
	p2y := math.Cos(a) * p

	// start
	qx := 0.0
	qy := q
	// end
	q2x := math.Sin(a) * q
	q2y := math.Cos(a) * q

	// reposition for a nice, natural view of label
	transform := "translate(%f,%f) "
	transform = fmt.Sprintf(transform, 0.5*wx, -1.0*(q-hx))
	if flip {
		transform = "translate(%f,%f) scale(1,-1) "
		transform = fmt.Sprintf(transform, 0.5*wx, q)
	}
	transform += "rotate(%f)"
	degrees := a2 * 180 / math.Pi
	transform = fmt.Sprintf(transform, degrees)

	// curved label
	curve := `
<path
	d="M%f %f A%f %f 0 0 1 %f %f L%f %f A%f %f 0 0 0 %f %f Z"
	transform="%s"
	fill="none" stroke="black" stroke-width="0.1"
/>`
	curve = fmt.Sprintf(curve, p2x, p2y, p, p, px, py, qx, qy, q, q, q2x, q2y, transform)

	// straightened label
	rect := ""
	if d {
		rect = `
<rect
	x="%f" y="%f"
	width="%f" height="%f"
	fill="none" stroke="black" stroke-width="0.1"
/>`
		rect = fmt.Sprintf(rect, 0.0, 0.0, wx, h)
	}

	// final composition
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
	o1x, o1y := wx, hx
	unit := "mm"
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
		fmt.Printf("t:%f b:%f h:%f p:%f a:%f q:%f q2x:%f q2y:%f o1x:%f o1y:%f\n", t, b, h, p, a*180/math.Pi, q, q2x, q2y, o1x, o1y)
	}
}
