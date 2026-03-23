package utils

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func SelectOption(page *rod.Page, parentSelector string, optionXPath string) {
	page.MustElementX(parentSelector).MustClick()
	el := page.MustElementX(optionXPath)
	MoveMouseToElement(page, el)
	el.MustClick()
	page.MustWaitDOMStable()
}

func MoveMouseToElement(page *rod.Page, el *rod.Element) error {
	window, err := page.GetWindow()
	if err != nil {
		return err
	}

	startPos := proto.Point{
		X: RandRange(10, *window.Width),
		Y: RandRange(10, *window.Height),
	}

	// Get the target's center position
	box := el.MustShape().Box()
	targetPos := proto.Point{X: box.X + box.Width/2, Y: box.Y + box.Height/2}

	// Choose a random control point for the Bézier curve
	controlX := startPos.X + (targetPos.X-startPos.X)*(0.2+rand.Float64()*0.6)
	controlY := startPos.Y + (targetPos.Y-startPos.Y)*(0.2+rand.Float64()*0.6)
	controlPoint := proto.Point{X: controlX, Y: controlY}

	steps := 1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		p := BezierCurvePoint(startPos, controlPoint, targetPos, t)

		page.Mouse.MustMoveTo(p.X, p.Y)
		time.Sleep(RandomDelay(10, 30))
	}

	page.Mouse.MustMoveTo(targetPos.X, targetPos.Y)
	return nil
}

// randomDelay generates a random delay between min and max milliseconds.
func RandomDelay(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min+1)+min) * time.Millisecond
}

func RandRange(min, max int) float64 {
	return float64(rand.Intn(max-min+1) + min)
}

// bezierCurvePoint calculates a point on a 3-point Bézier curve.
func BezierCurvePoint(p0, p1, p2 proto.Point, t float64) proto.Point {
	x := (1-t)*(1-t)*p0.X + 2*(1-t)*t*p1.X + t*t*p2.X
	y := (1-t)*(1-t)*p0.Y + 2*(1-t)*t*p1.Y + t*t*p2.Y

	return proto.Point{X: x, Y: y}
}
