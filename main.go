package main

import (
	"flag"
	"fmt"
	"math"
)

var (
	t, b, h       float64
	r, p, l, a, g float64
	k, c, d       bool
)

func init() {
	flag.Float64Var(&t, "T", 0, "Top diameter of the cone")
	flag.Float64Var(&b, "B", 0, "Bottom diameter of the cone")
	flag.Float64Var(&h, "H", 0, "Height of the cone")
	flag.Float64Var(&g, "G", 0, "Gap between ends, expressed as length")
	flag.BoolVar(&k, "K", false, "Draw entire unfolded cone")
	flag.BoolVar(&c, "C", false, "Print only numbers")
	flag.BoolVar(&d, "D", false, "Print straighten unfolded cone")
	flag.Parse()
}

func main() {
	t, b = math.Min(t, b), math.Max(t, b)
	//hack
	if t == b {
		t = b * 0.000001
	}
	r = math.Sqrt(math.Pow((b-t)/2, 2) + math.Pow(h, 2))
	l = math.Pi * t
	p = (t * r) / (b - t)
	q := p + r
	// radians
	a = l / p
	// take gap into acount
	g = math.Atan(g / q)
	a = a + g
	// degrees
	//a = a * 180 / math.Pi
	// half a
	a2 := a / 2

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

	curve := "M%f %f Q%f %f %f %f L%f %f Q%f %f %f %f Z"
	curve = fmt.Sprintf(curve, p2x, p2y, p1x, p1y, px, py, qx, qy, q1x, q1y, q2x, q2y)

	degrees := a2 * 180 / math.Pi

	svg := `<svg width="%f%s" height="%f%s" viewBox="%f %f %f %f" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><path fill="none" stroke="black" stroke-width="0.1" d="%s" transform="rotate(%f)"/>%s%s</svg>`
	ox, oy := 0.0, 0.0
	o1x, o1y := math.Tan(a)*q, q
	unit := "mm"
	kon := ""
	if k {
		kon = `<path d="M%f %f V%f H%f Z" transform="rotate(%f)" fill="none" stroke="black" stroke-width="0.1" />`
		kon = fmt.Sprintf(kon, ox, oy, o1y, o1x, degrees)
	}
	rect := ""
	if d {
		rect = `<rect x="%f" y="%f" width="%f" height="%f" fill="none" stroke="black" stroke-width="0.1" />`
		rect = fmt.Sprintf(rect, 0.0, 0.0, math.Pi*b, r)
	}
	svg = fmt.Sprintf(svg, o1x, unit, q, unit, ox, oy, o1x, o1y, curve, degrees, kon, rect)
	if !c {
		fmt.Println(svg)
	} else {
		fmt.Printf("t:%f b:%f h:%f r:%f p:%f p1x:%f p1y:%f a:%f g:%f q:%f q1x:%f q1y:%f q2x:%f q2y:%f o1x:%f o1y:%f\n", t, b, h, r, p, p1x, p1y, a*180/math.Pi, g, q, q1x, q1y, q2x, q2y, o1x, o1y)
	}
}
