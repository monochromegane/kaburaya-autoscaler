package autoscaler

import "math"

type Component interface {
	Work(float64) float64
}

type RoundAndMinimum struct {
	Minimum float64
}

func (c *RoundAndMinimum) Work(x float64) float64 {
	return math.Max(math.Round(x), c.Minimum)
}
