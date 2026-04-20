package utils

import (
	"fmt"
	"math/rand"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func SelectOption(page *rod.Page, parentSelector string, optionXPath string) error {
	parent, err := page.ElementX(parentSelector)
	if err != nil {
		return fmt.Errorf("dropdown parent not found (%s): %w", parentSelector, err)
	}
	if err := parent.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click dropdown parent: %w", err)
	}

	option, err := page.ElementX(optionXPath)
	if err != nil {
		return fmt.Errorf("option not found (%s): %w", optionXPath, err)
	}

	if err := MoveMouseToElement(page, option); err != nil {
		return fmt.Errorf("failed to move mouse to option: %w", err)
	}

	if err := option.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click option: %w", err)
	}

	if err := page.WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after selecting option: %w", err)
	}

	return nil
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

func FillInput(page *rod.Page, selector, value string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	return el.Input(value)
}

func FillInputTime(page *rod.Page, xpath string, t time.Time) error {
	el, err := page.ElementX(xpath)
	if err != nil {
		return fmt.Errorf("date input element not found (%s): %w", xpath, err)
	}
	if err := el.InputTime(t); err != nil {
		return fmt.Errorf("failed to input time into %s: %w", xpath, err)
	}
	return nil
}
